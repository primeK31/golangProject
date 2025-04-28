package domain

import (
	"context"
	"time"
)

type Bet struct {
	UserID    int
	EventID  string 
	Status string
	CreatedAt  time.Time 
	Money int
	Coefficient float32
	ExpectedResult float32
}

type BetService interface {
	Authenticate(ctx context.Context, email, password string) (*User, error)
	GetProfile(ctx context.Context) (*User, error)
	GetCurrentUser(ctx context.Context) (*User, error)
}
