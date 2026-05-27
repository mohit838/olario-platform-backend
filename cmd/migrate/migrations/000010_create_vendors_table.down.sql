ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_primary_vendor;

DROP TABLE IF EXISTS vendors;
