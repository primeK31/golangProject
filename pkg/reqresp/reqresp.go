package reqresp


type AuthRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Email    string `json:"email"`
    Username string `json:"username"`
    Password string `json:"password"`
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
