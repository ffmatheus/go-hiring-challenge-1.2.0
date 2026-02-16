INSERT INTO categories (code, name) VALUES
('CTG001', 'Clothing'),
('CTG002', 'Shoes'),
('CTG003', 'Accessories');

UPDATE products
SET category_id = (SELECT id FROM categories WHERE code = 'CTG001')
WHERE code IN ('PROD001', 'PROD004', 'PROD007');

UPDATE products
SET category_id = (SELECT id FROM categories WHERE code = 'CTG002')
WHERE code IN ('PROD002', 'PROD006');

UPDATE products
SET category_id = (SELECT id FROM categories WHERE code = 'CTG003')
WHERE code IN ('PROD003', 'PROD005', 'PROD008');
