-- +goose Up
create table auth_user (
                      id serial primary key,
                      name varchar(50) not null,
                      email varchar(50) not null,
                      role smallint not null,
                      created_at timestamp not null default now(),
                      updated_at timestamp
);

-- +goose Down
drop table auth_user;

