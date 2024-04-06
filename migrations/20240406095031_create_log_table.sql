-- +goose Up
create table log (
                           id serial primary key,
                           action varchar(100) not null,
                           model varchar(50),
                           model_id varchar(50),
                           query varchar(100),
                           query_row text,
                           created_at timestamp not null default now()
);

-- +goose Down
drop table log;

