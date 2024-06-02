-- +goose Up
-- +goose StatementBegin
create table access_rule (
                           id serial primary key,
                           endpoint varchar(255) not null,
                           role smallint not null,
                           created_at timestamp not null default now(),
                           updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table access_rule;
-- +goose StatementEnd
