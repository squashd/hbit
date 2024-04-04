-- +goose Up
ALTER TABLE
    character RENAME TO character_backup;

CREATE TABLE IF NOT EXISTS character_state (
    user_id TEXT NOT NULL,
    class_id TEXT NOT NULL,
    character_level INTEGER NOT NULL DEFAULT 1,
    experience INTEGER NOT NULL DEFAULT 0,
    health INTEGER NOT NULL DEFAULT 50,
    mana INTEGER NOT NULL DEFAULT 50,
    strength INTEGER NOT NULL DEFAULT 5,
    dexterity INTEGER NOT NULL DEFAULT 5,
    intelligence INTEGER NOT NULL DEFAULT 5,
    event_id TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (class_id) REFERENCES character_class (id)
);

INSERT INTO
    character_state (
        user_id,
        class_id,
        character_level,
        experience,
        health,
        mana,
        strength,
        dexterity,
        intelligence,
        event_id
    )
SELECT
    user_id,
    class_id,
    character_level,
    experience,
    health,
    mana,
    strength,
    dexterity,
    intelligence,
    "no_event" AS event_id
FROM
    character_backup;

CREATE INDEX IF NOT EXISTS character_state_user_id_index ON character_state (user_id);

ALTER TABLE
    quest RENAME TO quest_backup;

CREATE TABLE IF NOT EXISTS quest (
    quest_id TEXT PRIMARY KEY,
    quest_type TEXT NOT NULL NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    title TEXT UNIQUE NOT NULL,
    details TEXT NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
    quest (
        quest_id,
        title,
        quest_type,
        description,
        details
    )
SELECT
    id,
    title,
    "" AS quest_type,
    description,
    requirements AS details
FROM
    quest_backup;

ALTER TABLE
    user_quest RENAME TO user_quest_backup;

CREATE TABLE IF NOT EXISTS user_quest (
    user_id TEXT NOT NULL,
    quest_id TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    event_id TEXT NOT NULL,
    details TEXT NOT NULL DEFAULT '{}'
);

INSERT INTO
    user_quest (
        user_id,
        quest_id,
        completed,
        event_id,
        details
    )
SELECT
    user_id,
    quest_id,
    completed,
    "no_event" AS event_id,
    '{}' AS details
FROM
    user_quest_backup;

ALTER TABLE
    item RENAME TO item_backup;

CREATE TABLE IF NOT EXISTS item (
    item_id TEXT PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    href TEXT NOT NULL DEFAULT '',
    str_boost INTEGER NOT NULL DEFAULT 0,
    dex_boost INTEGER NOT NULL DEFAULT 0,
    int_boost INTEGER NOT NULL DEFAULT 0,
    slot TEXT NOT NULL DEFAULT 'none',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
    item (
        item_id,
        name,
        description,
        str_boost,
        dex_boost,
        int_boost,
        slot
    )
SELECT
    id,
    name,
    description,
    0 AS str_boost,
    0 AS dex_boost,
    0 AS int_boost,
    item_type AS slot
FROM
    item_backup;

-- +goose Down
ALTER TABLE
    character_backup RENAME TO character;

INSERT INTO
    character
SELECT
    DISTINCT user_id,
    class_id,
    character_level,
    experience,
    health,
    mana,
    strength,
    dexterity,
    intelligence
FROM
    (
        SELECT
            *,
            ROW_NUMBER() OVER (
                PARTITION BY user_id
                ORDER BY
                    timestamp DESC
            ) AS rn
        FROM
            character_state
    )
WHERE
    rn = 1;

DROP TABLE IF EXISTS character_state;

ALTER TABLE
    quest_backup RENAME TO quest;

INSERT INTO
    quest
SELECT
    id,
    title,
    description,
    requirements
FROM
    quest_backup;

DROP TABLE quest;

ALTER TABLE
    user_quest RENAME TO user_quest_backup_temp;

ALTER TABLE
    user_quest_backup RENAME TO user_quest;

INSERT INTO
    user_quest
SELECT
    user_id,
    quest_id,
    completed,
    event_id,
    details
FROM
    user_quest_backup_temp;

DROP TABLE IF EXISTS user_quest_backup_temp;

ALTER TABLE
    item RENAME TO item_temp;

CREATE TABLE IF NOT EXISTS item (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    attributes TEXT NOT NULL DEFAULT '{}',
    requirements TEXT NOT NULL DEFAULT '{}'
);

INSERT INTO
    item (id, slot)
SELECT
    id,
    slot,
    JSON_OBJECT(
        'str_boost',
        str_boost,
        'dex_boost',
        dex_boost,
        'int_boost',
        int_boost
    ) AS attributes,
    '{}' AS requirements
FROM
    item_temp;

ALTER TABLE
    item_backup RENAME TO item;

DROP TABLE item_temp;