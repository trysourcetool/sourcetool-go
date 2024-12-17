package sourcetool

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

type contextKey string

const sessionIDKey = contextKey("session_id")

type Session struct {
	ID    uuid.UUID
	State *State
}

func NewSession(id uuid.UUID) *Session {
	return &Session{
		ID:    id,
		State: NewState(),
	}
}

// WithSessionID returns a new context with the session ID
func WithSessionID(ctx context.Context, sessionID uuid.UUID) context.Context {
	return context.WithValue(ctx, sessionIDKey, sessionID)
}

// SessionIDFromContext retrieves the session ID from the context
func SessionIDFromContext(ctx context.Context) (uuid.UUID, error) {
	v := ctx.Value(sessionIDKey)
	if v == nil {
		return uuid.Nil, fmt.Errorf("session ID not found in context")
	}

	sessionID, ok := v.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid session ID type in context")
	}

	return sessionID, nil
}
