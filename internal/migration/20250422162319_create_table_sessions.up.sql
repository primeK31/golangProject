CREATE TABLE IF NOT EXISTS sessions (
    token TEXT NOT NULL,
    user_uuid BINARY(16) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_agent VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    INDEX idx_user_uuid (user_uuid),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
