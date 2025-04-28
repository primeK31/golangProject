package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Token 	   string
    UserUUID   uuid.UUID
    ExpiresAt  time.Time
    CreatedAt  time.Time
    UserAgent  string
    IPAddress  string
}
