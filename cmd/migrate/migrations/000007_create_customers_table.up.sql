CREATE TABLE customers (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  name VARCHAR(120) NOT NULL,
  email VARCHAR(255),
  phone VARCHAR(30),
  address TEXT,
  notes TEXT,
  email_opt_in BOOLEAN NOT NULL DEFAULT FALSE,
  loyalty_points_balance INT NOT NULL DEFAULT 0 CHECK (loyalty_points_balance >= 0),
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (tenant_id, id)
);

CREATE INDEX idx_customers_tenant_id ON customers(tenant_id);
CREATE INDEX idx_customers_phone ON customers(phone);
CREATE INDEX idx_customers_is_active ON customers(is_active);

COMMENT ON TABLE customers IS 'Tenant-scoped grocery customers for orders, receipts, and future loyalty workflows.';
COMMENT ON COLUMN customers.tenant_id IS 'Scopes the customer to one tenant.';
COMMENT ON COLUMN customers.name IS 'Customer display name used on orders and receipts.';
COMMENT ON COLUMN customers.email IS 'Optional customer email for receipts or communication.';
COMMENT ON COLUMN customers.phone IS 'Optional customer phone for lookup and communication.';
COMMENT ON COLUMN customers.email_opt_in IS 'Customer consent flag for email receipts and messages.';
COMMENT ON COLUMN customers.loyalty_points_balance IS 'Current loyalty balance for future rewards and bonuses.';
COMMENT ON COLUMN customers.is_active IS 'Disables customer selection without deleting order history.';
