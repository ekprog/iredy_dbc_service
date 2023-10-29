-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dbc_challenges
(
    id              SERIAL PRIMARY KEY NOT NULL,
    owner_id        bigint                      default null,
    category_id     bigint                      default null,

    -- Активность трекается пользователем?
    is_auto_track   bool               not null default false,

    visibility_type varchar(255)       not null default 'private',

    name            varchar(255)       not null,
    image           varchar(255)                default null,
    "desc"          varchar(1000),

    created_at      timestamp(0)       NOT NULL DEFAULT now(),
    updated_at      timestamp(0)       NOT NULL DEFAULT now(),
    deleted_at      timestamp(0)                DEFAULT null,

    unique (owner_id, name),

    constraint fk_user_id foreign key (owner_id) REFERENCES users (id) ON DELETE CASCADE,
    constraint fk_category_id foreign key (category_id) REFERENCES dbc_challenge_categories (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dbc_challenges;
-- +goose StatementEnd
