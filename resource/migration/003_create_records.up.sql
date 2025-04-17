CREATE TABLE records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    asset_id INTEGER NOT NULL,
    date DATE NOT NULL,
    bought_value DECIMAL(14,2) DEFAULT 0,
    current_value DECIMAL(14,2) DEFAULT 0,
    realized_value DECIMAL(14,2) DEFAULT 0,
    note TEXT DEFAULT NULL,
    FOREIGN KEY (asset_id) REFERENCES assets(id)
);

CREATE UNIQUE INDEX idx_record ON records(asset_id, date);
CREATE INDEX idx_record_date ON records(date);