-- +goose Up
CREATE TABLE IF NOT EXISTS code_verification (
    id              UUID PRIMARY KEY,
    problem_id      UUID NOT NULL,
    tests_code      VARCHAR(10000) NOT NULL,
    total_number    INTEGER NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS code_verification;
