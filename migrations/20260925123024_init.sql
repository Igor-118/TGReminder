-- +goose Up
-- SELECT 'up SQL query';

CREATE TABLE users (
    id integer primary key,
    'name' text,
    age int,
    time_zone text,
)

CREATE TABLE dates (
    dateid integer primary key
    userid integer foreign key,
    'description' text,
    time_announce timestamp,
)

-- +goose Down
SELECT 'down SQL query';
