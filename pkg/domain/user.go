package domain

import "context"

type User struct {
	ID       int
	Email    string
	Password string
	Username     string
	Balance int
}

type UserService interface {
	Authenticate(ctx context.Context, email, password string) (*User, error)
	GetProfile(ctx context.Context) (*User, error)
	GetCurrentUser(ctx context.Context) (*User, error)
}