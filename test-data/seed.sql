-- Seed initial data so test suites have stable SKUs to work against.

INSERT INTO products (sku, name, price_cents, stock, category) VALUES
    ('SKU-001', 'Mechanical Keyboard',    19999, 100, 'electronics'),
    ('SKU-002', 'Wireless Mouse',          4500, 100, 'electronics'),
    ('SKU-003', 'USB-C Cable 2m',           999, 500, 'accessories'),
    ('SKU-004', 'Webcam HD',               7800,  50, 'electronics'),
    ('SKU-LIMITED', 'Limited Edition Mug', 2500,   2, 'merch'),
    ('SKU-RACE', 'Inventory Race Item',    1500,   1, 'merch')
ON CONFLICT (sku) DO NOTHING;
