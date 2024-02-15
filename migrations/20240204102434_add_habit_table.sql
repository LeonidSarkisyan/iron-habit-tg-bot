-- +goose Up
-- +goose StatementBegin
CREATE TABLE habits (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL
);

CREATE TABLE timestamps (
    id SERIAL PRIMARY KEY,
    day VARCHAR(20) NOT NULL,
    time_ VARCHAR(5) NOT NULL,
    habit_id INTEGER REFERENCES habits(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS timestamps;
DROP TABLE IF EXISTS habits;
-- +goose StatementEnd
