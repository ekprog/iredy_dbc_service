-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks
(
    id         SERIAL PRIMARY KEY NOT NULL,
    user_id    bigint             not null,
    project_id bigint,

    name       varchar(255)       not null,
    "desc"       varchar(1000),
    priority   integer            not null default 0,
    done       boolean            not null default false,

    created_at timestamp(0)       NOT NULL DEFAULT now(),
    updated_at timestamp(0)       NOT NULL DEFAULT now(),
    deleted_at timestamp(0)                DEFAULT null,

    constraint fk_user_id foreign key (user_id) REFERENCES users (id) ON DELETE CASCADE,
    constraint fk_project_id foreign key (project_id) REFERENCES projects (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd
