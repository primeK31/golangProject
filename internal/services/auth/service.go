package auth

import (
	"errors"
	"fmt"
	"golangproject/pkg/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
    jwtSecret string
}

type Claims struct {
	UserUUID uuid.UUID `json:"user_uuid"`
	jwt.RegisteredClaims
}

func New(jwtSecret string) *Service {
    return &Service{jwtSecret: jwtSecret}
}

func (s *Service) GenerateToken(user *domain.User) (string, error) {
    claims := jwt.MapClaims{
        "user_uuid": user.UUID,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}

// ParseToken validates and parses JWT token
func (s *Service) ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(s.jwtSecret), nil
    })

    if err != nil {
        return nil, handleJWTError(err)
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        if claims.ExpiresAt.Time.Before(time.Now()) {
            return nil, fmt.Errorf("token expired")
        }
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token claims")
}

func handleJWTError(err error) error {
    switch {
    case errors.Is(err, jwt.ErrTokenExpired):
        return fmt.Errorf("token expired")
    case errors.Is(err, jwt.ErrTokenMalformed):
        return fmt.Errorf("malformed token")
    case errors.Is(err, jwt.ErrTokenSignatureInvalid):
        return fmt.Errorf("invalid signature")
    default:
        return fmt.Errorf("token parsing error: %w", err)
    }
}
