-- +goose Up
CREATE TABLE IF NOT EXISTS user_settings (
    user_id TEXT PRIMARY KEY,
    theme TEXT CHECK(theme IN ('light', 'dark')) NOT NULL DEFAULT 'light',
    display_name TEXT NOT NULL,
    email TEXT NOT NULL DEFAULT '',
    email_notifications BOOLEAN NOT NULL DEFAULT FALSE,
    reset_time TEXT NOT NULL DEFAULT '00:00',
    user_timezone TEXT NOT NULL DEFAULT 'UTC',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_settings_email ON user_settings (email);

CREATE INDEX IF NOT EXISTS idx_user_settings_reset_time ON user_settings (reset_time);

-- +goose Down
DROP TABLE IF EXISTS user_settings;