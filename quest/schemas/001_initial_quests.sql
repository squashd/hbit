-- +goose Up
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

-- +goose Down
DROP TABLE IF EXISTS quest_item;

DROP TABLE IF EXISTS quest;