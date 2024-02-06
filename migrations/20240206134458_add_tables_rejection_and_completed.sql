-- +goose Up
-- +goose StatementBegin
CREATE TABLE rejections (
   id SERIAL PRIMARY KEY,
   text TEXT,
   datetime_rejection TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   habit_id INT,
   FOREIGN KEY (habit_id) REFERENCES habits(id)
);

CREATE TABLE executions (
   id SERIAL PRIMARY KEY,
   datetime_completed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   habit_id INT,
   FOREIGN KEY (habit_id) REFERENCES habits(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rejections;
DROP TABLE IF EXISTS executions;
-- +goose StatementEnd
