-- +migrate Up

CREATE TABLE permissions
(
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    url         TEXT NOT NULL,
    http_method VARCHAR(16) NOT NULL,
    description TEXT
);
