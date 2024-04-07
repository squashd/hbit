-- +goose Up
ALTER TABLE
    task
ADD
    COLUMN difficulty TEXT NOT NULL DEFAULT 'easy';

-- +goose Down
ALTER TABLE
    task DROP COLUMN difficulty;