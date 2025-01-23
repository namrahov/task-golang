-- +migrate Up

CREATE TABLE IF NOT EXISTS users_roles (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Indexes to optimize queries
CREATE INDEX IF NOT EXISTS idx_users_roles_user_id ON users_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_users_roles_role_id ON users_roles(role_id);