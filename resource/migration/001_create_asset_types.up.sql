CREATE TABLE asset_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    is_cash BOOLEAN NOT NULL,
    is_liability BOOLEAN NOT NULL,
    sequence INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL
);

CREATE UNIQUE INDEX idx_asset_type ON asset_types(name);