package repositories

import (
    "context"
    "database/sql"
    "golangproject/pkg/domain"
    "fmt"
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
	/* fmt.Println(session.Token)
	fmt.Println(session.UserID)
	fmt.Println(session.ExpiresAt)
	fmt.Println(session.CreatedAt)
	fmt.Println(session.UserAgent)
	fmt.Println(session.IPAddress) */
    _, err := r.db.ExecContext(ctx,
        "INSERT INTO sessions (token, user_id, expires_at, created_at, user_agent, ip_address) VALUES (?, ?, ?, ?, ?, ?)",
        session.Token,
        session.UserID,
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
        "SELECT token, user_id, expires_at, created_at, user_agent, ip_address FROM sessions WHERE token = ?",
        token,
    ).Scan(
        &session.Token,
        &session.UserID,
        &session.ExpiresAt,
        &session.CreatedAt,
        &session.UserAgent,
        &session.IPAddress,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve session: %w", err)
    }
    return &session, nil
}


func (r *sessionRepo) Delete(ctx context.Context, token string) error {
    _, err := r.db.ExecContext(ctx,
        "DELETE FROM sessions WHERE token = ?",
        token,
    )
    if err != nil {
        return fmt.Errorf("failed to delete session: %w", err)
    }
    return nil
}
