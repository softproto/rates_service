-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS usdtrub(
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    ask REAL NOT NULL,
    bid REAL NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS usdtrub;
-- +goose StatementEnd
