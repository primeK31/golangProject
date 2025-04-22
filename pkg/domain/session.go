package domain

import "time"

type Session struct {
	Token 	   string
    UserID     int
    ExpiresAt  time.Time
    CreatedAt  time.Time
    UserAgent  string
    IPAddress  string
}
