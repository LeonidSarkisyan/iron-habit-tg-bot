-- +goose Up
-- +goose StatementBegin
ALTER TABLE habits
    ADD COLUMN time_completed INTEGER DEFAULT 60 CHECK (time_completed >= 15 AND time_completed <= 300),
    ADD COLUMN time_warning INTEGER DEFAULT 15 CHECK (time_warning >= 5 AND time_warning <= 60);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE habits
    DROP COLUMN IF EXISTS time_warning,
    DROP COLUMN IF EXISTS time_completed;
-- +goose StatementEnd
