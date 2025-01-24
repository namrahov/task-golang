-- +migrate Up

CREATE TABLE roles_permissions
(
    id            BIGSERIAL NOT NULL PRIMARY KEY,
    role_id       BIGINT    NOT NULL,
    permission_id BIGINT    NOT NULL,
    CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
    CONSTRAINT fk_permission FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE
);
