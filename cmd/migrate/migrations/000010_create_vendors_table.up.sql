CREATE TABLE vendors (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  name VARCHAR(160) NOT NULL,
  contact_person VARCHAR(120),
  phone VARCHAR(30),
  email VARCHAR(255),
  address TEXT,
  notes TEXT,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (tenant_id, id),
  UNIQUE (tenant_id, name)
);

CREATE INDEX idx_vendors_tenant_id ON vendors(tenant_id);
CREATE INDEX idx_vendors_phone ON vendors(phone);
CREATE INDEX idx_vendors_is_active ON vendors(is_active);

COMMENT ON TABLE vendors IS 'Tenant-scoped suppliers or salesmen who provide grocery products.';
COMMENT ON COLUMN vendors.tenant_id IS 'Scopes the vendor to one tenant.';
COMMENT ON COLUMN vendors.contact_person IS 'Optional human contact name for the vendor or salesman.';
COMMENT ON COLUMN vendors.is_active IS 'Disables a vendor without deleting product purchase history.';

ALTER TABLE products
  ADD CONSTRAINT fk_products_primary_vendor
  FOREIGN KEY (tenant_id, primary_vendor_id) REFERENCES vendors(tenant_id, id) ON DELETE SET NULL (primary_vendor_id);

COMMENT ON CONSTRAINT fk_products_primary_vendor ON products IS 'Keeps preferred product vendor tenant-safe.';
