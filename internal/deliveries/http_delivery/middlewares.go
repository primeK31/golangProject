package http_delivery

import (
	"context"
	"errors"
	"net/http"
	"strings"
	//"time"

	"golangproject/internal/app/config"
	"golangproject/internal/services/auth"
	"golangproject/internal/services/middleware"
	//"golangproject/internal/services/session.go"
	"golangproject/internal/services/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


func AuthMiddleware(authService *auth.Service, userService *user.Service) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract the token from the Authorization header
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
                http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
                return
            }

            token := strings.TrimPrefix(authHeader, "Bearer ")

            userUUID, err := ValidateToken(r.Context(), token)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            user, err := userService.GetProfile(r.Context(), userUUID)
            if user == nil || err != nil {
                http.Error(w, "User not found", http.StatusUnauthorized)
                return
            }

            // Inject user into context
            ctx := context.WithValue(r.Context(), middleware.CurrentUserKey, userUUID)

            // Call the next handler with the updated context
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

/* func SessionMiddleware(sessionService *session.Service) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            cookie, err := r.Cookie("session_id")
            if err != nil {
                respondWithError(w, http.StatusUnauthorized, "Session required")
                return
            }

            session, err := sessionService.GetSession(r.Context(), cookie.Value)
            if err != nil {
                respondWithError(w, http.StatusUnauthorized, "Invalid session")
                return
            }

            if time.Now().After(session.ExpiresAt) {
                respondWithError(w, http.StatusUnauthorized, "Session expired")
                return
            }

            ctx := context.WithValue(r.Context(), middleware.CurrentUserKey, session.UserUUID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
} */


func ValidateToken(ctx context.Context, tokenString string) (uuid.UUID, error) {
    token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Ensure the signing method is HMAC
        cfg := config.LoadConfig()
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(cfg.JWTSecret), nil
    })

    if err != nil {
        return uuid.Nil, err
    }

    claims, ok := token.Claims.(*auth.Claims)
    if !ok || !token.Valid {
        return uuid.Nil, errors.New("invalid token")
    }

    return claims.UserUUID, nil
}
