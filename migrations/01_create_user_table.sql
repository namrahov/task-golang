-- +migrate Up

CREATE TABLE users
(
    id                  BIGSERIAL    NOT NULL PRIMARY KEY,
    username            VARCHAR(255),
    email               VARCHAR(255) NOT NULL UNIQUE,
    password            VARCHAR(256) NOT NULL,
    phone_number        VARCHAR(16),
    accept_notification BOOLEAN   DEFAULT FALSE,
    is_active           BOOLEAN   DEFAULT FALSE,
    inactivated_date    DATE,
    full_name           VARCHAR(255),
    description         TEXT,
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP DEFAULT NOW()
);
