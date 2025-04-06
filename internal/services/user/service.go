package user

import (
	"context"
	"errors"
	"fmt"

	"golangproject/internal/repositories"
	"golangproject/internal/services/middleware"
	"golangproject/pkg/domain"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrInvalidPassword   = errors.New("invalid password")
)

type Service struct {
	repo       repositories.UserRepository
}

func New(repo repositories.UserRepository) *Service {
    return &Service{repo: repo}
}

func (s *Service) HashPassword(password string) (string, error) {
	// Добавляем "пепец" к паролю перед хэшированием
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *Service) comparePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password)) == nil 
}

func (s *Service) Register(ctx context.Context, user domain.User) (*domain.User, error) {
	hashedPassword, err := s.HashPassword(user.Password)
    if err != nil {
        return nil, err
    }
    user.Password = hashedPassword

    return s.repo.AddUser(ctx, user)
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		//fmt.Println("not found by email")
		return nil, ErrUserNotFound
	}

	if !s.comparePassword(user.Password, password) {
		//fmt.Println("password")
		return nil, ErrInvalidPassword
	}

	return user, nil
}

func (s *Service) GetCurrentUser(ctx context.Context) (*domain.User, error) {
    // Получаем пользователя из контекста
    user, ok := ctx.Value(middleware.CurrentUserKey).(*domain.User)
	//fmt.Println("test")
	//fmt.Println(user.Username)
    if !ok || user == nil {
        return nil, ErrUnauthorized
    }
    
    // При необходимости обновляем данные из БД
    freshUser, err := s.repo.GetByID(ctx, user.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to refresh user data: %w", err)
    }
    
    return freshUser, nil
}

func (s *Service) GetProfile(ctx context.Context, id int) (*domain.User, error) {
	userID := id
	//fmt.Println("PROFILE")
	//fmt.Println(userID)
	
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	
	return user, nil
}