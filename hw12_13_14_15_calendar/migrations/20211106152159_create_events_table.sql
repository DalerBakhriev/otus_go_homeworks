-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
   id serial PRIMARY KEY,
   title TEXT NOT NULL,
   creation_time TIMESTAMP NOT NULL,
   start_time TIMESTAMP NOT NULL,
   end_time TIMESTAMP NOT NULL,
   description TEXT NULL,
   time_before_notification TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
