-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS projects
(
    id         SERIAL PRIMARY KEY not null,
    user_id    bigint             not null,

    name       varchar(255)       not null,
    "desc"     varchar(1000)               default null,
    color      varchar(100)                default null,

    "order"    int                         default 0,
    parent_id  bigint                      default null,

    created_at timestamp(0)       NOT NULL DEFAULT now(),
    updated_at timestamp(0)       NOT NULL DEFAULT now(),
    deleted_at timestamp(0)                DEFAULT NULL,

    constraint fk_user_id foreign key (user_id) REFERENCES users (id) ON DELETE CASCADE,
    constraint fk_parent_id foreign key (parent_id) REFERENCES projects (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS projects;
-- +goose StatementEnd
