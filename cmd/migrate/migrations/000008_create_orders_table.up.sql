CREATE TABLE orders (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  order_number VARCHAR(80) NOT NULL,
  customer_id BIGINT,
  created_by BIGINT,
  status VARCHAR(30) NOT NULL DEFAULT 'draft',
  subtotal_amount NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (subtotal_amount >= 0),
  discount_amount NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
  tax_amount NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (tax_amount >= 0),
  total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
  paid_amount NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (paid_amount >= 0),
  loyalty_points_earned INT NOT NULL DEFAULT 0 CHECK (loyalty_points_earned >= 0),
  payment_method VARCHAR(50),
  notes TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (tenant_id, id),
  UNIQUE (tenant_id, order_number),
  FOREIGN KEY (tenant_id, customer_id) REFERENCES customers(tenant_id, id) ON DELETE SET NULL (customer_id),
  FOREIGN KEY (tenant_id, created_by) REFERENCES users(tenant_id, id) ON DELETE SET NULL (created_by),
  CHECK (status IN ('draft', 'confirmed', 'paid', 'cancelled', 'refunded'))
);

CREATE INDEX idx_orders_tenant_id ON orders(tenant_id);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_created_by ON orders(created_by);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);

COMMENT ON TABLE orders IS 'Tenant-scoped grocery checkout records.';
COMMENT ON COLUMN orders.tenant_id IS 'Scopes the order to one tenant.';
COMMENT ON COLUMN orders.order_number IS 'Tenant-unique human-readable order or invoice number.';
COMMENT ON COLUMN orders.customer_id IS 'Optional customer linked to the order.';
COMMENT ON COLUMN orders.created_by IS 'Optional user who created or processed the order.';
COMMENT ON COLUMN orders.status IS 'Simple order lifecycle for v1 checkout workflows.';
COMMENT ON COLUMN orders.subtotal_amount IS 'Total before discount and tax.';
COMMENT ON COLUMN orders.discount_amount IS 'Order-level discount amount.';
COMMENT ON COLUMN orders.tax_amount IS 'Order-level tax amount.';
COMMENT ON COLUMN orders.total_amount IS 'Final payable amount after discount and tax.';
COMMENT ON COLUMN orders.paid_amount IS 'Amount collected from the customer.';
COMMENT ON COLUMN orders.loyalty_points_earned IS 'Loyalty points earned by the customer from this order.';
COMMENT ON COLUMN orders.payment_method IS 'Payment method label such as cash, card, mobile, or other local method.';
