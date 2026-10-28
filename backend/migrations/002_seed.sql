INSERT INTO products (sku, name, price, stock, is_active) VALUES
('CRM-001', 'Keramik Lantai 40x40 cm', 85000, 120, TRUE),
('CRM-002', 'Keramik Lantai 60x60 cm', 145000, 75, TRUE),
('CRM-003', 'Keramik Dinding 25x40 cm', 95000, 60, TRUE),
('CRM-004', 'Granit Tile 60x60 cm', 185000, 45, TRUE),
('CRM-005', 'Granit Tile 80x80 cm', 275000, 30, TRUE),
('CRM-006', 'Perekat Keramik 25 kg', 75000, 50, TRUE),
('CRM-007', 'Nat Keramik 1 kg', 18000, 80, TRUE),
('SAN-001', 'Kloset Duduk', 1250000, 12, TRUE),
('SAN-002', 'Wastafel Keramik', 675000, 15, TRUE),
('SAN-003', 'Keran Air Stainless', 185000, 25, TRUE),
('SAN-004', 'Shower Set Kamar Mandi', 450000, 0, FALSE),
('HOM-001', 'Kitchen Sink Stainless', 550000, 10, TRUE)
ON CONFLICT (sku) DO NOTHING;

INSERT INTO stock_movements (product_id, type, qty, note)
SELECT p.id, 'in', p.stock, 'Initial stock'
FROM products p
WHERE NOT EXISTS (SELECT 1 FROM stock_movements WHERE product_id = p.id)
  AND p.stock > 0;
