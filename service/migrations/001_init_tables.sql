CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(63) NOT NULL
);

CREATE TABLE IF NOT EXISTS games (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    mines INT DEFAULT 10,
    cols INT DEFAULT 8,
    rows INT DEFAULT 8,
    owner_id INT REFERENCES users (id),
    status VARCHAR(31) DEFAULT 'open',
    max_players INT DEFAULT 2,
    is_public BOOLEAN DEFAULT true,
    invite_token VARCHAR(16) UNIQUE DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_invite_token ON games (invite_token);

CREATE TABLE IF NOT EXISTS players (
    user_id INT REFERENCES users (id),
    game_id VARCHAR(36) REFERENCES games (id),
    CONSTRAINT unique_user_game UNIQUE (user_id, game_id)
);
CREATE INDEX idx_user_id_game_id ON players (user_id, game_id);
