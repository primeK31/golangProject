package http_delivery

import (
	"encoding/json"
	"errors"
	"net/http"

	_ "golangproject/cmd/app/docs"
	"golangproject/internal/repositories"
	"golangproject/internal/services/auth"
	"golangproject/internal/services/user"
	"golangproject/pkg/domain"
    "golangproject/pkg/reqresp"
)


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


// @Summary User login
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param credentials body AuthRequest true "Login credentials"
// @Success 200 {object} reqresp.AuthResponse
// @Failure 400 {object} reqresp.ErrorResponse
// @Failure 401 {object} reqresp.ErrorResponse
// @Failure 500 {object} reqresp.ErrorResponse
// @Router /login [post]
func LoginHandler(authService *auth.Service, userService *user.Service) http.HandlerFunc {
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

        respondWithJSON(w, http.StatusOK, reqresp.AuthResponse{Token: token})
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

func ProfileHandler(userService *user.Service) http.HandlerFunc {
    // Получаем пользователя из контекста через сервис
    return func(w http.ResponseWriter, r *http.Request) {
        //fmt.Println("lol")
        user, err := userService.GetCurrentUser(r.Context())
        //fmt.Println("lol")
        if err != nil {
            respondWithError(w, http.StatusUnauthorized, "Authentication required")
            return
        }

        // Формируем безопасный ответ (без чувствительных данных)
        safeUser := struct {
            ID       int    `json:"id"`
            Username string `json:"username"`
            Email    string `json:"email"`
        }{
            ID:       user.ID,
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
