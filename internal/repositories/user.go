package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"golangproject/pkg/domain"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user exists")
)

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
    AddUser(ctx context.Context, user domain.User) (*domain.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
    const query = `SELECT id, email, username FROM users WHERE id = ?`
    
    var user domain.User
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Email,
        &user.Username,
    )
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT 
			id, 
			email, 
			password
		FROM users 
		WHERE email = ?`

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)

	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("%w: email %s", ErrUserNotFound, email)
	case err != nil:
		log.Printf("Error fetching user by email: %v", err)
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

func (r *userRepo) AddUser(ctx context.Context, user domain.User) (*domain.User, error) {
	const query = `
        INSERT INTO users (
            username, 
            email, 
            password
        ) VALUES (?, ?, ?)`
	// Выполняем запрос с контекстом
	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	)

	if err != nil {
		// Обработка ошибки дублирования уникального поля
		if isDuplicateError(err) {
			return nil, ErrUserExists
		}
        //fmt.Println("lol")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// Вспомогательная функция для проверки дублирования
func isDuplicateError(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return mysqlErr.Number == 1062 // Код ошибки дублирования для MySQL
	}
	// Для PostgreSQL используйте проверку на pgErr.Code == "23505"
	return false
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
    const query = `
        SELECT 
            id, 
            email, 
            password
        FROM users 
        WHERE email = ?`

    var user domain.User

    err := r.db.QueryRowContext(ctx, query, username).Scan(
        &user.ID,
        &user.Email,
        &user.Password,
    )

    switch {
    case err == sql.ErrNoRows:
        return nil, fmt.Errorf("%w: username %s", ErrUserNotFound, username)
    case err != nil:
        log.Printf("Error fetching user by username: %v", err)
        return nil, fmt.Errorf("database error: %w", err)
    }
	return &domain.User{Username: username}, nil
}
