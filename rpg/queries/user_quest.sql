CREATE TABLE IF NOT EXISTS user_quest (
    user_id TEXT NOT NULL,
    quest_id TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    event_id TEXT NOT NULL,
    details TEXT NOT NULL DEFAULT '{}'
);

-- :name UpdateUserQuest :one
INSERT INTO
    user_quest (user_id, quest_id, completed, event_id, details)
VALUES
    (?, ?, ?, ?, ?) returning *;