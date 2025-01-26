-- +migrate Up

CREATE TABLE boards
(
    id         BIGSERIAL PRIMARY KEY, -- Primary key, auto-incrementing
    name       VARCHAR(255) NOT NULL, -- Column for board name, with a max length of 255
    created_by VARCHAR(64)  NOT NULL, -- Column for the creator's name
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);