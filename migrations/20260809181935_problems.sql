-- +goose Up
CREATE TABLE IF NOT EXISTS problems (
    id UUID PRIMARY KEY,
    name VARCHAR(70) UNIQUE NOT NULL,
    description VARCHAR(10000) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS problems;
