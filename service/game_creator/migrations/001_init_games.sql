CREATE TABLE IF NOT EXISTS games (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    mines INT DEFAULT 40,
    cols INT DEFAULT 16,
    rows INT DEFAULT 16,
    owner_id INT REFERENCES users (id),
    status VARCHAR(128) DEFAULT 'open',
    max_players INT DEFAULT 2,
    is_public BOOLEAN DEFAULT true,
    invite_token VARCHAR(16) UNIQUE DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_invite_token ON games (invite_token);
