package auth

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("AuthDatabase: not found")

type AuthDatabase interface {
	// If email could not be found error is present
	AccountExists(ctx context.Context, email string) (string, error)
	AddSession(ctx context.Context, sessionId string, email string) (error)
	GetSession(ctx context.Context, sessionId string, since time.Time) (AccountDetails, error)
}

