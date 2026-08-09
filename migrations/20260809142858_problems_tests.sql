-- +goose Up
CREATE TABLE IF NOT EXISTS problems_tests (
    id UUID PRIMARY KEY,
    problem_id UUID UNIQUE NOT NULL,
    tests_code VARCHAR(10000) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    submissions_count INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS problems_tests;
