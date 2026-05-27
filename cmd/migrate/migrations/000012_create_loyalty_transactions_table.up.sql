CREATE TABLE loyalty_transactions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  customer_id BIGINT NOT NULL,
  order_id BIGINT,
  points INT NOT NULL,
  reason VARCHAR(80) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  FOREIGN KEY (tenant_id, customer_id) REFERENCES customers(tenant_id, id) ON DELETE CASCADE,
  FOREIGN KEY (tenant_id, order_id) REFERENCES orders(tenant_id, id) ON DELETE SET NULL (order_id)
);

CREATE INDEX idx_loyalty_transactions_tenant_id ON loyalty_transactions(tenant_id);
CREATE INDEX idx_loyalty_transactions_customer_id ON loyalty_transactions(customer_id);
CREATE INDEX idx_loyalty_transactions_order_id ON loyalty_transactions(order_id);

COMMENT ON TABLE loyalty_transactions IS 'Ledger of loyalty point changes for customers.';
COMMENT ON COLUMN loyalty_transactions.points IS 'Positive or negative point change.';
COMMENT ON COLUMN loyalty_transactions.reason IS 'Reason for points, such as order_purchase, bonus, or adjustment.';
