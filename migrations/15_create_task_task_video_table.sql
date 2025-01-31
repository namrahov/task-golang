-- +migrate Up

CREATE TABLE IF NOT EXISTS task_task_video (
    task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    task_video_id BIGINT NOT NULL REFERENCES task_video(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, task_video_id)
    );

-- Indexes to optimize queries
CREATE INDEX IF NOT EXISTS idx_task_task_video_task_id ON task_task_video(task_id);
CREATE INDEX IF NOT EXISTS idx_task_task_video_task_video_id ON task_task_video(task_video_id);