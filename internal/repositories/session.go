package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"golangproject/pkg/domain"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByToken(ctx context.Context, token string) (*domain.Session, error)
	Delete(ctx context.Context, token string) error
}

type sessionRepo struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepo{db: db}
}

func (r *sessionRepo) Create(ctx context.Context, session *domain.Session) error {
	// No need to parse the UUID as it's already a uuid.UUID type in the struct
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO sessions (token, user_uuid, expires_at, created_at, user_agent, ip_address)
         VALUES (?, UUID_TO_BIN(?), ?, ?, ?, ?)`, 
		session.Token,
		uuid.New().String(), // Store UUID as is, database driver will handle conversion
		session.ExpiresAt,
		session.CreatedAt,
		session.UserAgent,
		session.IPAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *sessionRepo) GetByToken(ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session

	err := r.db.QueryRowContext(ctx,
		"SELECT token, user_uuid, expires_at, created_at, user_agent, ip_address FROM sessions WHERE token = ?",
		token,
	).Scan(
		&session.Token,
		&session.UserUUID,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UserAgent,
		&session.IPAddress,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	return &session, nil
}

func (r *sessionRepo) Delete(ctx context.Context, token string) error {
    fmt.Println("delete session")
	result, err := r.db.ExecContext(ctx,
		"DELETE FROM sessions WHERE token = ?",
		token,
	)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Check if any row was affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}
	
	return nil
}
