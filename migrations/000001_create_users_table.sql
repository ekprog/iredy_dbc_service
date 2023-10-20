-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id          bigint unique,
    score       integer      not null default 0,
    score_daily integer      not null default 0,

    created_at  timestamp(0) NOT NULL DEFAULT now(),
    updated_at  timestamp(0) NOT NULL DEFAULT now(),
    deleted_at  timestamp(0)          DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
