CREATE TABLE bought_value_offsets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    asset_id INTEGER NOT NULL,
    effective_date DATE NOT NULL,
    offset_price DECIMAL(14,2) NOT NULL,
    note TEXT DEFAULT NULL,
    FOREIGN KEY (asset_id) REFERENCES assets(id)
);

CREATE INDEX idx_bought_value_offsets ON bought_value_offsets(asset_id);