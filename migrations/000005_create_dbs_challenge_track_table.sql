-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dbc_challenges_tracks
(
    id           SERIAL PRIMARY KEY          NOT NULL,
    user_id      bigint                      not null,
    challenge_id bigint                               default null,

    "date"       timestamp without time zone not null DEFAULT null,
    "done"       bool                        not null DEFAULT false,

    created_at   timestamp(0)                NOT NULL DEFAULT now(),
    updated_at   timestamp(0)                NOT NULL DEFAULT now(),

    unique (user_id, "date"),

    constraint fk_user_id foreign key (user_id) REFERENCES users (id) ON DELETE CASCADE,
    constraint fk_category_id foreign key (challenge_id) REFERENCES dbc_challenges (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dbc_challenges_tracks;
-- +goose StatementEnd
