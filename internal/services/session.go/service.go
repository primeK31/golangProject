package session

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golangproject/internal/repositories"
	"golangproject/pkg/domain"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Service struct {
    repo            repositories.SessionRepository
    sessionDuration time.Duration
    secretKey       []byte
}

func New(repo repositories.SessionRepository, sessionDuration time.Duration, secretKey []byte) *Service {
    return &Service{
        repo:            repo,
        sessionDuration: sessionDuration,
        secretKey:       secretKey,
    }
}

func (s *Service) CreateSession(ctx context.Context, userUUID uuid.UUID, r *http.Request) (*domain.Session, error) {
    token, err := s.generateJWT(userUUID, r)
    if err != nil {
        return nil, fmt.Errorf("failed to generate JWT: %w", err)
    }

    session := &domain.Session{
        Token:     token,
        UserUUID:    userUUID,
        ExpiresAt: time.Now().Add(s.sessionDuration),
        CreatedAt: time.Now(),
        UserAgent: r.UserAgent(),
        IPAddress: getIPAddress(r),
    }

    if err := s.repo.Create(ctx, session); err != nil {
        return nil, fmt.Errorf("failed to save session: %w", err)
    }

    return session, nil
}

func (s *Service) GetSession(ctx context.Context, tokenStr string) (*domain.Session, error) {
    _, err := s.parseJWT(tokenStr)
    if err != nil {
        return nil, fmt.Errorf("failed to validate JWT: %w", err)
    }

    // Fetch session by token (not by user ID)
    session, err := s.repo.GetByToken(ctx, tokenStr)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve session: %w", err)
    }

    return session, nil
}


func (s *Service) DeleteSession(ctx context.Context, tokenStr string) error {
    _, err := s.parseJWT(tokenStr)
    if err != nil {
        return fmt.Errorf("failed to validate JWT: %w", err)
    }
    fmt.Println(tokenStr)

    return s.repo.Delete(ctx, tokenStr)
}

func (s *Service) generateJWT(userUUID uuid.UUID, r *http.Request) (string, error) {
    claims := jwt.MapClaims{
        "user_uuid": userUUID,
        "exp":     time.Now().Add(s.sessionDuration).Unix(), // Expiry time
        "iat":     time.Now().Unix(), // Issued at time
        "user_agent": r.UserAgent(),
        "ip_address": getIPAddress(r),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(s.secretKey)
    if err != nil {
        return "", fmt.Errorf("failed to sign JWT: %w", err)
    }

    return tokenString, nil
}

func (s *Service) parseJWT(tokenStr string) (int, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        // Ensure that the token's method is HMAC
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return s.secretKey, nil
    })
    if err != nil {
        return 0, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID, ok := claims["user_uuid"].(float64) // user_id is float64 in JWT claims
        if !ok {
            return 0, fmt.Errorf("user_uuid not found in token")
        }
        return int(userID), nil
    }

    return 0, fmt.Errorf("invalid token")
}

func getIPAddress(r *http.Request) string {
    ip := r.Header.Get("X-Forwarded-For")
    if ip == "" {
        ip = r.RemoteAddr
    }
    return ip
}
