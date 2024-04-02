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

CREATE TABLE IF NOT EXISTS quest (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    -- This could be a JSON string detailing quest requirements like required items or previous quest completion.
    requirements TEXT,
    rewards TEXT -- JSON
);

CREATE TABLE IF NOT EXISTS quest_item (
    quest_id TEXT NOT NULL,
    item_id TEXT NOT NULL,
    quantity_required INTEGER NOT NULL,
    FOREIGN KEY (quest_id) REFERENCES quest (id),
    FOREIGN KEY (item_id) REFERENCES item (id),
    PRIMARY KEY (quest_id, item_id)
);

CREATE TABLE IF NOT EXISTS user_quest (
    user_id TEXT NOT NULL,
    quest_id TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (user_id, quest_id)
);

-- +goose Down
DROP TABLE IF EXISTS item;

DROP TABLE IF EXISTS character;

DROP TABLE IF EXISTS character_class;

DROP TABLE IF EXISTS quest_item;

DROP TABLE IF EXISTS quest;