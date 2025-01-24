package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
	"google.golang.org/protobuf/proto"
)

type Client interface {
	RegisterHandler(MessageHandlerFunc)
	Enqueue(string, proto.Message)
	EnqueueWithResponse(string, proto.Message) (*websocketv1.Message, error)
	Close() error
	Wait() error
}

type Config struct {
	URL            string
	APIKey         string
	InstanceID     uuid.UUID
	PingInterval   time.Duration
	ReconnectDelay time.Duration
	OnReconnecting func()
	OnReconnected  func()
}

type client struct {
	config Config
	dialer *websocket.Dialer

	conn   *websocket.Conn
	connMu sync.RWMutex

	messageQueue chan *websocketv1.Message
	done         chan error

	responses map[string]chan *websocketv1.Message
	respMu    sync.RWMutex

	handler   MessageHandlerFunc
	handlerMu sync.RWMutex
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
		messageQueue: make(chan *websocketv1.Message, 100),
		done:         make(chan error, 1),
		dialer:       websocket.DefaultDialer,
		responses:    make(map[string]chan *websocketv1.Message),
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	go c.pingPongLoop()
	go c.readMessages()
	go c.sendEnqueuedMessagesLoop()

	return c, nil
}

func (c *client) RegisterHandler(handler MessageHandlerFunc) {
	c.handlerMu.Lock()
	defer c.handlerMu.Unlock()
	c.handler = handler
}

func (c *client) handleMessage(msg *websocketv1.Message) error {
	// Handle responses
	c.respMu.RLock()
	respChan, exists := c.responses[msg.Id]
	c.respMu.RUnlock()

	if exists {
		respChan <- msg
		return nil
	}

	// Handle calls
	c.handlerMu.RLock()
	handler := c.handler
	c.handlerMu.RUnlock()

	if handler == nil {
		return fmt.Errorf("no message handler registered")
	}

	return handler(msg)
}

func (c *client) connect() error {
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	header.Set("X-Instance-Id", c.config.InstanceID.String())

	conn, _, err := c.dialer.Dial(c.config.URL, header)
	if err != nil {
		return err
	}

	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(c.config.PingInterval * 2))
	})

	c.connMu.Lock()
	c.conn = conn
	c.connMu.Unlock()

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
		c.connMu.RLock()
		conn := c.conn
		c.connMu.RUnlock()

		if conn == nil {
			return
		}

		if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
			log.Printf("ping failed: %v", err)
			c.connMu.Lock()
			conn.Close()
			c.conn = nil
			c.connMu.Unlock()

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
		var data []byte
		c.connMu.RLock()
		conn := c.conn
		c.connMu.RUnlock()

		if conn == nil {
			time.Sleep(time.Second)
			continue
		}

		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				c.done <- nil
				return
			}

			log.Printf("read error: %v", err)
			c.connMu.Lock()
			conn.Close()
			c.conn = nil
			c.connMu.Unlock()

			go func() {
				if err := c.reconnect(); err != nil {
					log.Printf("reconnection failed: %v", err)
				}
			}()
			continue
		}

		msg, err := unmarshalMessage(data)
		if err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		if err := c.handleMessage(msg); err != nil {
			log.Printf("error handling message: %v", err)
			c.sendException(msg.Id, err)
		}
	}
}

func (c *client) sendEnqueuedMessagesLoop() {
	defer close(c.messageQueue)

	const batchInterval = 10 * time.Millisecond
	var messageBuffer []*websocketv1.Message

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
				c.connMu.RLock()
				conn := c.conn
				c.connMu.RUnlock()

				if conn == nil {
					time.Sleep(time.Second)
					continue
				}

				var remainingMessages []*websocketv1.Message
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

func (c *client) send(msg *websocketv1.Message) error {
	data, err := marshalMessage(msg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *client) Enqueue(id string, payload proto.Message) {
	msg, err := NewMessage(id, payload)
	if err != nil {
		log.Printf("error creating message: %v", err)
		return
	}
	c.messageQueue <- msg
}

func (c *client) EnqueueWithResponse(id string, payload proto.Message) (*websocketv1.Message, error) {
	respChan := make(chan *websocketv1.Message, 1)
	c.respMu.Lock()
	c.responses[id] = respChan
	c.respMu.Unlock()

	defer func() {
		c.respMu.Lock()
		delete(c.responses, id)
		c.respMu.Unlock()
	}()

	c.Enqueue(id, payload)

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
	c.connMu.Lock()
	defer c.connMu.Unlock()

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Printf("error closing connection: %v", err)
			return err
		}
		c.conn = nil
	}

	c.done <- nil
	return nil
}
