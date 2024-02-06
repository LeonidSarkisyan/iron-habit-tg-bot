-- +goose Up
-- +goose StatementBegin
ALTER TABLE habits
    ADD COLUMN time_completed INTEGER DEFAULT 60 CHECK (time_completed >= 0 AND time_completed <= 60),
    ADD COLUMN time_warning INTEGER DEFAULT 5 CHECK (time_warning >= 0 AND time_warning <= 60);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE habits
    DROP COLUMN IF EXISTS time_warning,
    DROP COLUMN IF EXISTS time_completed;
-- +goose StatementEnd
