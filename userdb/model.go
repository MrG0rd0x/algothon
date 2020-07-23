package userdb

import (
	"context"
	"time"
)

// UserDB is
type UserDB interface {
	Get(context.Context, string) (User, error)
	Save(context.Context, *User, time.Duration) error
	Verify(context.Context) bool
	Login(context.Context, string, string) bool
	Register(context.Context, string, string) error
}

// User contains information about a user session
type User struct {
	username string
	token    []byte
	ttl      time.Duration
}
