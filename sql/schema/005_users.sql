-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red boolean not null
default false;

-- +goose Down
ALTER TABLE users
DROP COLUMN is_chirpy_red;