-- +goose Up
CREATE TABLE IF NOT EXISTS achievement (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    requirement TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS user_achievement (
    user_id TEXT NOT NULL,
    achievement_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, achievement_id)
);

CREATE INDEX user_achievement_user_id_idx ON user_achievement (user_id);

CREATE TABLE IF NOT EXISTS task_log (
    event_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    task_id TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    payload jsonb NOT NULL
);

CREATE INDEX task_log_user_id_idx ON task_log (user_id);

CREATE TABLE IF NOT EXISTS quest_log(
    event_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    quest_id TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    payload jsonb NOT NULL
);

CREATE INDEX quest_log_user_id_idx ON quest_log (user_id);

-- +goose Down
DROP TABLE IF EXISTS achievement;

DROP TABLE IF EXISTS user_achievement;

DROP TABLE IF EXISTS task_log;

DROP TABLE IF EXISTS quest_log;