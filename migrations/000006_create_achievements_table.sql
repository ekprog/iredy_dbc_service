-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS achievements
(
    id           SERIAL PRIMARY KEY NOT NULL,
    user_id      bigint             not null,

    title        varchar(255)       not null,
    "desc"       varchar(1000)      not null,
    image        varchar(255)                default null,

    trigger_name varchar(255)       not null,


    created_at   timestamp(0)       NOT NULL DEFAULT now(),
    updated_at   timestamp(0)       NOT NULL DEFAULT now(),

    unique (title),

    constraint fk_user_id foreign key (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS achievements;
-- +goose StatementEnd
