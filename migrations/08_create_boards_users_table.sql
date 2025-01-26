-- +migrate Up

CREATE TABLE IF NOT EXISTS users_boards (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    board_id BIGINT NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, board_id)
    );

-- Indexes to optimize queries
CREATE INDEX IF NOT EXISTS idx_users_boards_user_id ON users_boards(user_id);
CREATE INDEX IF NOT EXISTS idx_users_boards_board_id ON users_boards(board_id);