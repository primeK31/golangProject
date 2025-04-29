package http_delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	_ "golangproject/docs/swagger"
	"golangproject/internal/repositories"
	"golangproject/internal/services/auth"
	"golangproject/internal/services/session.go"
	"golangproject/internal/services/user"
	"golangproject/pkg/domain"
	"golangproject/pkg/reqresp"

	"github.com/google/uuid"
)


// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

// RegisterHandler godoc
// @Summary User registration
// @Description Register new user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body reqresp.RegisterRequest true "Registration data"
// @Success 201 {object} reqresp.RegisterResponse
// @Failure 400 {object} reqresp.ErrorResponse
// @Failure 409 {object} reqresp.ErrorResponse
// @Failure 500 {object} reqresp.ErrorResponse
// @Router /register [post]
func RegisterHandler(userService *user.Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req reqresp.RegisterRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            respondWithError(w, http.StatusBadRequest, "invalid request")
            return
        }

        _, err := userService.Register(r.Context(), domain.User{
            Username: req.Username,
            Email:    req.Email,
            Password: req.Password,
        })

        switch {
        case errors.Is(err, repositories.ErrUserExists):
            respondWithError(w, http.StatusConflict, "user already exists")
        case err != nil:
            respondWithError(w, http.StatusInternalServerError, "registration failed")
        default:
            respondWithJSON(w, http.StatusCreated, reqresp.RegisterResponse{Message: "User registered successfully"})
        }
    }
}


// LoginHandler godoc
// @Summary User login
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body reqresp.AuthRequest true "Login credentials"
// @Success 200 {object} reqresp.AuthResponse
// @Failure 400 {object} reqresp.ErrorResponse
// @Failure 401 {object} reqresp.ErrorResponse
// @Failure 500 {object} reqresp.ErrorResponse
// @Router /login [post]
func LoginHandler(authService *auth.Service, userService *user.Service, sessionService *session.Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req reqresp.AuthRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            respondWithError(w, http.StatusBadRequest, "Invalid request")
            return
        }

        user, err := userService.Authenticate(r.Context(), req.Email, req.Password)
        if err != nil {
            respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
            return
        }

        token, err := authService.GenerateToken(user)
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
            return
        }

        session, err := sessionService.CreateSession(r.Context(), user.UUID, r)
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Failed to create session")
            return
        }

        // Set the JWT token in a secure HTTP-only cookie
        http.SetCookie(w, &http.Cookie{
            Name:     "session_id",
            Value:    token,         // The generated JWT token
            Expires:  session.ExpiresAt,
            HttpOnly: true,          // Ensure the cookie is not accessible via JavaScript
            Secure:   true,          // Set to true in production with HTTPS
            SameSite: http.SameSiteStrictMode, // Strict mode for cookie security
            Path:     "/",           // Cookie is available for the whole site
        })

        // Respond with the JWT token (optional, as it's already in the cookie)
        respondWithJSON(w, http.StatusOK, reqresp.AuthResponse{Token: token})
    }
}


// LogoutHandler godoc
// @Summary User logout
// @Description Logout user and clear session cookie
// @Tags auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} reqresp.RegisterResponse
// @Failure 400 {object} reqresp.ErrorResponse
// @Failure 500 {object} reqresp.ErrorResponse
// @Router /logout [post]
func LogoutHandler(sessionService *session.Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("session_id")
        if err != nil {
            respondWithError(w, http.StatusBadRequest, "No active session")
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:     "session_id",
            Value:    "",
            Path:     "/",                         // must match the original Path
            Domain:   cookie.Domain,              // if you set Domain originally
            Expires:  time.Unix(0, 0),             // in the past
            MaxAge:   -1,                          // delete immediately
            HttpOnly: true,
            Secure:   true,                        // HTTPS only in prod
            SameSite: http.SameSiteStrictMode,
        })

        respondWithJSON(w, http.StatusOK, map[string]string{
            "message": "Successfully logged out",
        })
    }
}


type UserHandler struct {
    userService domain.UserService
}

func NewUserHandler(us domain.UserService) *UserHandler {
    return &UserHandler{
        userService: us,
    }
}


// ProfileHandler godoc
// @Summary Get user profile
// @Description Get user details from JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} reqresp.RegisterResponse
// @Failure 401 {object} reqresp.ErrorResponse
// @Failure 500 {object} reqresp.ErrorResponse
// @Router /profile [get] // Fixed from [post] to [get]
func ProfileHandler(userService *user.Service) http.HandlerFunc {
    // Получаем пользователя из контекста через сервис
    return func(w http.ResponseWriter, r *http.Request) {
        user, err := userService.GetCurrentUser(r.Context())
        if err != nil {
            respondWithError(w, http.StatusUnauthorized, "Authentication required")
            return
        }

        safeUser := struct {
            UUID     uuid.UUID    `json:"uuid"`
            Username string `json:"username"`
            Email    string `json:"email"`
        }{
            UUID:       user.UUID,
            Username: user.Username,
            Email:    user.Email,
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(safeUser)
    }
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}
