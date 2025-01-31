-- +migrate Up

CREATE TABLE task_video
(
    id         BIGSERIAL PRIMARY KEY,
    file_type  VARCHAR(255),
    file_path  VARCHAR(255) UNIQUE NOT NULL,
    file_name  VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);
