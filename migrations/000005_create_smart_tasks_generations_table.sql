-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS smart_tasks_gen
(
    id            SERIAL PRIMARY KEY NOT NULL,
    user_id       bigint             not null,
    smart_task_id bigint             not null,

    period        varchar(255)       not null,
    datetime      timestamp(0)       not null,

    created_at    timestamp(0)       NOT NULL DEFAULT now(),
    updated_at    timestamp(0)       NOT NULL DEFAULT now(),
    deleted_at    timestamp(0)                DEFAULT null,

    constraint fk_smart_task_id foreign key (smart_task_id) REFERENCES smart_tasks (id) ON DELETE CASCADE,
    constraint fk_user_id foreign key (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS smart_tasks_gen;
-- +goose StatementEnd
