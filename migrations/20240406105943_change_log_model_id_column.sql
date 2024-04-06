-- +goose Up
alter table log
    alter column model_id type integer using model_id::integer;

-- +goose Down
alter table log
    alter column model_id type varchar(50) using model_id::varchar(50);