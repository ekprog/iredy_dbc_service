-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dbc_challenge_categories
(
    id         SERIAL PRIMARY KEY not null,
    user_id    bigint             not null,

    name       varchar(255)       not null,

    created_at timestamp(0)       NOT NULL DEFAULT now(),
    updated_at timestamp(0)       NOT NULL DEFAULT now(),
    deleted_at timestamp(0)                DEFAULT NULL,

    unique (user_id, name),

    constraint fk_user_id foreign key (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dbc_challenge_categories;
-- +goose StatementEnd
