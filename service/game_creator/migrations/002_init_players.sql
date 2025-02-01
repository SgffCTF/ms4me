CREATE TABLE IF NOT EXISTS players (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users (id),
    game_id VARCHAR(36) REFERENCES games (id),
    CONSTRAINT unique_user_game UNIQUE (user_id, game_id)
);
CREATE INDEX idx_user_id_game_id ON players (user_id, game_id);
