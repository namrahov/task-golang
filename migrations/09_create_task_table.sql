-- +migrate Up

CREATE TABLE tasks
(
    id             BIGSERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    priority       VARCHAR(64)  NOT NULL,
    status         VARCHAR(64)  NOT NULL,

    -- Allow these columns to default to NULL explicitly:
    created_by_id  BIGINT DEFAULT NULL,
    changed_by_id  BIGINT DEFAULT NULL,
    assigned_by_id BIGINT DEFAULT NULL,
    assigned_to_id BIGINT DEFAULT NULL,

    board_id   BIGINT,
    deadline   TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_board FOREIGN KEY (board_id) REFERENCES boards (id)
);
