CREATE TABLE categories (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  parent_id BIGINT,
  name VARCHAR(100) NOT NULL,
  code_prefix VARCHAR(20) NOT NULL,
  description TEXT,
  sort_order INT NOT NULL DEFAULT 0,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (tenant_id, id),
  UNIQUE (tenant_id, name),
  UNIQUE (tenant_id, code_prefix),
  FOREIGN KEY (tenant_id, parent_id) REFERENCES categories(tenant_id, id) ON DELETE SET NULL (parent_id)
);

CREATE INDEX idx_categories_tenant_id ON categories(tenant_id);
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_is_active ON categories(is_active);

COMMENT ON TABLE categories IS 'Tenant-scoped grocery product categories, with optional parent-child hierarchy.';
COMMENT ON COLUMN categories.tenant_id IS 'Scopes the category to one tenant.';
COMMENT ON COLUMN categories.parent_id IS 'Optional parent category for nested category navigation.';
COMMENT ON COLUMN categories.name IS 'Tenant-unique category name.';
COMMENT ON COLUMN categories.code_prefix IS 'Short tenant-unique prefix used when generating manual product codes.';
COMMENT ON COLUMN categories.sort_order IS 'Manual ordering value for category lists.';
COMMENT ON COLUMN categories.is_active IS 'Hides or disables a category without deleting products.';
