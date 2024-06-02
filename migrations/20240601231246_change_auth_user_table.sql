-- +goose Up
-- +goose StatementBegin
alter table auth_user
    rename column name to username;

alter table auth_user
    add password_hash text;

create unique index auth_user__username
    on auth_user (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table auth_user
    rename column username to name;

alter table auth_user
    drop column password_hash;

drop index auth_user__username;
-- +goose StatementEnd
