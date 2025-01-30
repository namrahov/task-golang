-- +migrate Up

CREATE TABLE IF NOT EXISTS task_task_image (
    task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    task_image_id BIGINT NOT NULL REFERENCES task_image(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, task_image_id)
);

-- Indexes to optimize queries
CREATE INDEX IF NOT EXISTS idx_task_task_image_task_id ON task_task_image(task_id);
CREATE INDEX IF NOT EXISTS idx_task_task_image_task_image_id ON task_task_image(task_image_id);