-- +migrate Up

CREATE TABLE roles
(
    id   BIGSERIAL   NOT NULL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE
);

-- Insert admin and user roles
INSERT INTO roles (name)
VALUES ('user');
INSERT INTO roles (name)
VALUES ('admin');
