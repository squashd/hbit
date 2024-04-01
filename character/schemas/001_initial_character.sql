-- +goose Up
CREATE TABLE IF NOT EXISTS character_class (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS character (
    user_id TEXT NOT NULL PRIMARY KEY,
    class_id TEXT NOT NULL,
    character_level INTEGER NOT NULL DEFAULT 1,
    experience INTEGER NOT NULL DEFAULT 0,
    health INTEGER NOT NULL DEFAULT 50,
    mana INTEGER NOT NULL DEFAULT 50,
    strength INTEGER NOT NULL DEFAULT 5,
    dexterity INTEGER NOT NULL DEFAULT 5,
    intelligence INTEGER NOT NULL DEFAULT 5,
    FOREIGN KEY (class_id) REFERENCES character_class (id)
);

CREATE INDEX IF NOT EXISTS character_user_id_index ON character (user_id);

CREATE TABLE IF NOT EXISTS item (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    item_type TEXT NOT NULL DEFAULT '',
    -- e.g., "consumable", "equipment"
    attributes TEXT -- This could be a JSON string detailing item-specific attributes like attack power, defense, etc.
);

-- +goose Down
DROP TABLE IF EXISTS item;

DROP TABLE IF EXISTS character;

DROP TABLE IF EXISTS character_class;