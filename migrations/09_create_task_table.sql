-- +migrate Up

CREATE TABLE tasks
(
    id          BIGINT PRIMARY KEY,      -- Primary key for the Task table
    name        VARCHAR(255) NOT NULL,   -- Name of the task
    priority    VARCHAR(64)  NOT NULL,   -- Priority (can map to an enum in the application)
    status      VARCHAR(64)  NOT NULL,   -- Status (can map to an enum in the application)
    created_by  BIGINT       NOT NULL,   -- Foreign key to the User who created the task
    changed_by  BIGINT,                  -- Foreign key to the User who last changed the task
    assigned_to BIGINT,                  -- Foreign key to the User assigned to the task
    board_id    BIGINT,                  -- Foreign key to the associated Board
    deadline    TIMESTAMP    NOT NULL,   -- Deadline for the task
    created_at  TIMESTAMP DEFAULT NOW(), -- Timestamp for task creation
    updated_at  TIMESTAMP DEFAULT NOW(), -- Timestamp for the last update
    CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT fk_changed_by FOREIGN KEY (changed_by) REFERENCES users (id),
    CONSTRAINT fk_assigned_to FOREIGN KEY (assigned_to) REFERENCES users (id),
    CONSTRAINT fk_board FOREIGN KEY (board_id) REFERENCES boards (id)
);
