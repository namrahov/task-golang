-- +migrate Up

CREATE TABLE users
(
    id                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username            VARCHAR(255),
    email               VARCHAR(255) NOT NULL UNIQUE,
    password            VARCHAR(256) NOT NULL,
    phone_number        VARCHAR(16),
    accept_notification BOOLEAN   DEFAULT TRUE,
    is_active           BOOLEAN   DEFAULT TRUE,
    inactivated_date    DATE,
    full_name           VARCHAR(255),
    description         TEXT,
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP DEFAULT NOW()
);
