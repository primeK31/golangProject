CREATE TABLE sessions (
    token TEXT NOT NULL,
    user_id INT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    user_agent VARCHAR(512) NOT NULL,
    ip_address VARCHAR(45) NOT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
