-- +goose Up

CREATE TABLE users (
    id bigserial primary key,
    tg_id bigint not null unique,
    username text,
    age int,
    time_zone text not null
);

CREATE TABLE reminders (
    dateid bigserial primary key,
    user_id bigint not null references users(id),
    description text,
    time_announce timestamp
);

-- +goose Down
DROP TABLE dates;
DROP TABLE users;
