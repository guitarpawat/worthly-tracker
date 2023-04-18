CREATE TABLE assets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    broker TEXT NOT NULL,
    type_id INTEGER NOT NULL,
    default_increment DECIMAL(14,2) DEFAULT 0,
    sequence INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL,
    FOREIGN KEY (type_id) REFERENCES asset_types(id)
);

CREATE UNIQUE INDEX idx_asset ON assets(name, broker);