package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client interface {
	Enqueue(string, MessageMethod, any)
	EnqueueWithResponse(string, MessageMethod, any) (*Message, error)
	Close() error
	Wait() error
}

type Config struct {
	URL            string
	APIKey         string
	PingInterval   time.Duration
	ReconnectDelay time.Duration
	CallHandler    func(*Message) error
	OnReconnecting func()
	OnReconnected  func()
}

type client struct {
	config       Config
	conn         *websocket.Conn
	messageQueue chan *Message
	done         chan error
	mu           sync.Mutex
	dialer       *websocket.Dialer
	responses    map[string]chan *Message
	respMu       sync.RWMutex
}

func NewClient(config Config) (Client, error) {
	if config.PingInterval == 0 {
		config.PingInterval = 1 * time.Second
	}
	if config.ReconnectDelay == 0 {
		config.ReconnectDelay = 1 * time.Second
	}

	c := &client{
		config:       config,
		messageQueue: make(chan *Message, 100),
		done:         make(chan error, 1),
		dialer:       websocket.DefaultDialer,
		responses:    make(map[string]chan *Message),
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	go c.pingPongLoop()
	go c.readMessages()
	go c.sendEnqueuedMessagesLoop()

	return c, nil
}

func (c *client) connect() error {
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	conn, _, err := c.dialer.Dial(c.config.URL, header)
	if err != nil {
		return err
	}

	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(c.config.PingInterval * 2))
	})

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	return nil
}

func (c *client) reconnect() error {
	if c.config.OnReconnecting != nil {
		c.config.OnReconnecting()
	}

	for {
		err := c.connect()
		if err == nil {
			if c.config.OnReconnected != nil {
				c.config.OnReconnected()
			}
			return nil
		}

		log.Printf("reconnection failed: %v, retrying in %v", err, c.config.ReconnectDelay)
		time.Sleep(c.config.ReconnectDelay)
	}
}

func (c *client) pingPongLoop() {
	ticker := time.NewTicker(c.config.PingInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.mu.Lock()
		conn := c.conn
		c.mu.Unlock()

		if conn == nil {
			return
		}

		if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
			log.Printf("ping failed: %v", err)
			c.mu.Lock()
			conn.Close()
			c.conn = nil
			c.mu.Unlock()

			go func() {
				if err := c.reconnect(); err != nil {
					log.Printf("reconnection failed: %v", err)
				}
			}()
			return
		}
	}
}

func (c *client) readMessages() {
	for {
		var msg *Message
		c.mu.Lock()
		conn := c.conn
		c.mu.Unlock()

		if conn == nil {
			time.Sleep(time.Second)
			continue
		}

		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				c.done <- nil
				return
			}

			log.Printf("read error: %v", err)
			c.mu.Lock()
			conn.Close()
			c.conn = nil
			c.mu.Unlock()

			go func() {
				if err := c.reconnect(); err != nil {
					log.Printf("reconnection failed: %v", err)
				}
			}()
			continue
		}

		if msg == nil {
			log.Printf("received nil message")
			continue
		}

		if err := c.handleMessage(msg); err != nil {
			log.Printf("error handling message: %v", err)
		}
	}
}

func (c *client) sendEnqueuedMessagesLoop() {
	defer close(c.messageQueue)

	const batchInterval = 10 * time.Millisecond
	var messageBuffer []*Message

	for {
		select {
		case msg, ok := <-c.messageQueue:
			if !ok {
				for _, m := range messageBuffer {
					c.send(m)
				}
				return
			}
			messageBuffer = append(messageBuffer, msg)
		default:
			if len(messageBuffer) > 0 {
				c.mu.Lock()
				conn := c.conn
				c.mu.Unlock()

				if conn == nil {
					time.Sleep(time.Second)
					continue
				}

				var remainingMessages []*Message
				for _, msg := range messageBuffer {
					if err := c.send(msg); err != nil {
						remainingMessages = append(remainingMessages, msg)
						log.Printf("error sending message: %v", err)
						break
					}
					time.Sleep(time.Millisecond)
				}

				messageBuffer = remainingMessages

				if len(remainingMessages) == 0 {
					time.Sleep(batchInterval)
				}
			}
		}
	}
}

func (c *client) handleMessage(msg *Message) error {
	switch msg.Kind {
	case MessageKindCall:
		if err := c.config.CallHandler(msg); err != nil {
			return err
		}
	case MessageKindResponse:
		c.respMu.RLock()
		respChan, exists := c.responses[msg.ID]
		c.respMu.RUnlock()

		if exists {
			respChan <- msg
			return nil
		}
	}
	return nil
}

func (c *client) send(msg *Message) error {
	return c.conn.WriteJSON(msg)
}

func (c *client) Enqueue(id string, method MessageMethod, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("error marshalling data: %v", err)
		return
	}

	c.messageQueue <- &Message{
		ID:      id,
		Kind:    MessageKindCall,
		Method:  method,
		Payload: b,
	}
}

func (c *client) EnqueueWithResponse(id string, method MessageMethod, data any) (*Message, error) {
	respChan := make(chan *Message, 1)
	c.respMu.Lock()
	c.responses[id] = respChan
	c.respMu.Unlock()

	defer func() {
		c.respMu.Lock()
		delete(c.responses, id)
		c.respMu.Unlock()
	}()

	c.Enqueue(id, method, data)

	select {
	case resp := <-respChan:
		return resp, nil
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

func (c *client) Wait() error {
	return <-c.done
}

func (c *client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		// if err := c.conn.WriteControl(
		// 	websocket.CloseMessage,
		// 	websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		// 	time.Now().Add(time.Second),
		// ); err != nil {
		// 	log.Printf("error sending close message: %v", err)
		// 	return err
		// }
		if err := c.conn.Close(); err != nil {
			log.Printf("error closing connection: %v", err)
			return err
		}
		c.conn = nil
	}

	c.done <- nil

	return nil
}
