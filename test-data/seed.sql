CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    sku TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    price_cents BIGINT NOT NULL CHECK (price_cents >= 0),
    stock INT NOT NULL DEFAULT 0,
    category TEXT NOT NULL DEFAULT 'general',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO products (sku, name, price_cents, stock, category) VALUES
    ('SKU-001', 'Mechanical Keyboard',    19999, 100, 'electronics'),
    ('SKU-002', 'Wireless Mouse',          4500, 100, 'electronics'),
    ('SKU-003', 'USB-C Cable 2m',           999, 500, 'accessories'),
    ('SKU-004', 'Webcam HD',               7800,  50, 'electronics'),
    ('SKU-LIMITED', 'Limited Edition Mug', 2500,   2, 'merch'),
    ('SKU-RACE', 'Inventory Race Item',    1500,   1, 'merch')
ON CONFLICT (sku) DO NOTHING;
