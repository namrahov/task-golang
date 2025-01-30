-- +migrate Up

CREATE TABLE task_attachment_file
(
    id                 BIGSERIAL PRIMARY KEY,
    user_id            BIGINT DEFAULT NULL,
    task_id            BIGINT DEFAULT NULL,
    attachment_file_id BIGINT DEFAULT NULL,

    CONSTRAINT unique_attachment_file_id UNIQUE (attachment_file_id), -- Correct way to define UNIQUE constraint

    CONSTRAINT fk_task_attachment_file_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT fk_task_attachment_file_tasks FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE SET NULL,
    CONSTRAINT fk_task_attachment_file_attachment FOREIGN KEY (attachment_file_id) REFERENCES attachment_file (id) ON DELETE SET NULL
);
