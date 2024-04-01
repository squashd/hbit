-- +goose Up
CREATE TABLE IF NOT EXISTS user_quest (
    user_id INTEGER NOT NULL,
    quest_id INTEGER NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (user_id, quest_id)
);

-- +goose Down