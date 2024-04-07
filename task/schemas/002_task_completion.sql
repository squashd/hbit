-- +goose Up
ALTER TABLE
    task RENAME TO task_backup;

CREATE TABLE IF NOT EXISTS task (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    text TEXT NOT NULL DEFAULT '',
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    task_type TEXT NOT NULL DEFAULT 'daily',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
    task (
        id,
        user_id,
        title,
        text,
        is_completed,
        task_type,
        created_at,
        updated_at
    )
SELECT
    id,
    user_id,
    title,
    text,
    false AS is_completed,
    task_type,
    created_at,
    updated_at
FROM
    task_backup;

DROP TABLE task_backup;

-- +goose Down
ALTER TABLE
    task RENAME TO task_backup;

CREATE TABLE IF NOT EXISTS task (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    text TEXT NOT NULL DEFAULT '',
    data JSON NOT NULL DEFAULT '{}',
    task_type TEXT NOT NULL DEFAULT 'daily',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
    task (
        id,
        user_id,
        title,
        text,
        data,
        task_type,
        created_at,
        updated_at
    )
SELECT
    id,
    user_id,
    title,
    text,
    '{}' AS data,
    task_type,
    created_at,
    updated_at
FROM
    task_backup;