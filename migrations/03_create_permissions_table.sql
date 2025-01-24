-- +migrate Up

CREATE TABLE permissions
(
    id          BIGSERIAL   NOT NULL PRIMARY KEY,
    url         TEXT        NOT NULL,
    http_method VARCHAR(16) NOT NULL,
    description TEXT
);
