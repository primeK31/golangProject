package reqresp

import (
    _ "golangproject/docs/swagger"
    "github.com/google/uuid"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type BetRequest struct {
	Amount  float64 `json:"amount" validate:"required,gt=0"`
	PredictedOutcome    string `json:"predictedOutcome" validate:"required,gt=1"`
	EventID string  `json:"eventId" validate:"required,uuid"`
}

type BetResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	BetID   uuid.UUID `json:"bet_id"`
}

type ServiceUnavailableResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	RetryAt int64  `json:"retry_after"`
}

type HandlerResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
}
