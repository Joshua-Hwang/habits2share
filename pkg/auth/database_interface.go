package auth

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("AuthDatabase: not found")

type AuthDatabase interface {
	// If email could not be found error is present
	// returns the userId
	GetUserIdFromEmail(ctx context.Context, email string) (string, error)
	AddSession(ctx context.Context, sessionId string, userId string) (error)
	GetUserIdFromSession(ctx context.Context, sessionId string, since time.Time) (string, error)
}

