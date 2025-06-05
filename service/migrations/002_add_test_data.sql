-- Пользователи (все с паролем 'password')
INSERT INTO users (username, password) VALUES
('alice',  '$2b$10$u1O7ghD48U2gxYtM/v8FvuU6/x1wX6VQuT9K4hFjOT.9.9TVA6uE2'),
('bob',    '$2b$10$u1O7ghD48U2gxYtM/v8FvuU6/x1wX6VQuT9K4hFjOT.9.9TVA6uE2'),
('charlie','$2b$10$u1O7ghD48U2gxYtM/v8FvuU6/x1wX6VQuT9K4hFjOT.9.9TVA6uE2'),
('diana',  '$2b$10$u1O7ghD48U2gxYtM/v8FvuU6/x1wX6VQuT9K4hFjOT.9.9TVA6uE2');

-- Игры (mines, cols, rows, max_players — по умолчанию)
INSERT INTO games (id, title, owner_id) VALUES
('game-uuid-001', 'Alice vs Bob', 1),
('game-uuid-002', 'Charlie waiting', 3),
('game-uuid-003', 'Public Game 1', 2),
('game-uuid-004', 'Diana match', 4);

INSERT INTO players (user_id, game_id) VALUES
(1, 'game-uuid-001'), -- Alice
(2, 'game-uuid-001'), -- Bob
(3, 'game-uuid-002'), -- Charlie
(2, 'game-uuid-003'), -- Bob
(4, 'game-uuid-004'); -- Diana
