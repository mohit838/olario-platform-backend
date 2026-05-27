CREATE TABLE products (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  category_id BIGINT,
  primary_vendor_id BIGINT,
  sku VARCHAR(80),
  product_code VARCHAR(80) NOT NULL,
  barcode VARCHAR(100),
  name VARCHAR(160) NOT NULL,
  description TEXT,
  unit VARCHAR(30) NOT NULL DEFAULT 'pcs',
  price NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
  cost_price NUMERIC(12, 2) CHECK (cost_price IS NULL OR cost_price >= 0),
  stock_quantity INT NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
  alert_at_stock INT NOT NULL DEFAULT 5 CHECK (alert_at_stock >= 0),
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  is_featured BOOLEAN NOT NULL DEFAULT FALSE,
  is_discounted BOOLEAN NOT NULL DEFAULT FALSE,
  discount_rate NUMERIC(5, 2) NOT NULL DEFAULT 0 CHECK (discount_rate >= 0 AND discount_rate <= 100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (tenant_id, id),
  UNIQUE (tenant_id, sku),
  UNIQUE (tenant_id, product_code),
  UNIQUE (tenant_id, barcode),
  FOREIGN KEY (tenant_id, category_id) REFERENCES categories(tenant_id, id) ON DELETE SET NULL (category_id)
);

CREATE INDEX idx_products_tenant_id ON products(tenant_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_products_is_deleted ON products(is_deleted);

COMMENT ON TABLE products IS 'Tenant-scoped grocery items available for sale or inventory tracking.';
COMMENT ON COLUMN products.tenant_id IS 'Scopes the product to one tenant.';
COMMENT ON COLUMN products.category_id IS 'Optional tenant-safe category for browsing and reporting.';
COMMENT ON COLUMN products.primary_vendor_id IS 'Preferred vendor or salesman for restocking; foreign key is added after vendors exist.';
COMMENT ON COLUMN products.sku IS 'Optional tenant-unique internal stock keeping unit.';
COMMENT ON COLUMN products.product_code IS 'Tenant-unique manual product code generated from the category prefix, usable when barcode scanning is unavailable.';
COMMENT ON COLUMN products.barcode IS 'Optional tenant-unique barcode used for checkout or inventory scanning.';
COMMENT ON COLUMN products.unit IS 'Selling unit such as pcs, kg, litre, pack, or box.';
COMMENT ON COLUMN products.price IS 'Current selling price for the product.';
COMMENT ON COLUMN products.cost_price IS 'Optional purchase/cost price for margin reporting.';
COMMENT ON COLUMN products.stock_quantity IS 'Current stock count kept simple for v1; detailed changes are recorded in inventory_movements.';
COMMENT ON COLUMN products.alert_at_stock IS 'Low-stock threshold for alerts; defaults to 5 but can be changed per product.';
COMMENT ON COLUMN products.is_deleted IS 'Soft-delete flag so historical orders can remain valid.';
COMMENT ON COLUMN products.discount_rate IS 'Percentage discount from 0 to 100 when product discounting is enabled.';
