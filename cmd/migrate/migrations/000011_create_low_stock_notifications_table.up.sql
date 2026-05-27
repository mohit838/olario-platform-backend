CREATE TABLE low_stock_notifications (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  product_id BIGINT NOT NULL,
  threshold_quantity INT NOT NULL CHECK (threshold_quantity >= 0),
  current_quantity INT NOT NULL CHECK (current_quantity >= 0),
  status VARCHAR(30) NOT NULL DEFAULT 'open',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  resolved_at TIMESTAMPTZ,
  FOREIGN KEY (tenant_id, product_id) REFERENCES products(tenant_id, id) ON DELETE CASCADE,
  CHECK (status IN ('open', 'resolved', 'ignored'))
);

CREATE INDEX idx_low_stock_notifications_tenant_id ON low_stock_notifications(tenant_id);
CREATE INDEX idx_low_stock_notifications_product_id ON low_stock_notifications(product_id);
CREATE INDEX idx_low_stock_notifications_status ON low_stock_notifications(status);

COMMENT ON TABLE low_stock_notifications IS 'Tracks low-stock events when product quantity reaches its alert threshold.';
COMMENT ON COLUMN low_stock_notifications.threshold_quantity IS 'Product threshold used when the notification was created.';
COMMENT ON COLUMN low_stock_notifications.current_quantity IS 'Stock quantity at the time the notification was created.';
COMMENT ON COLUMN low_stock_notifications.status IS 'Open until stock is handled, resolved, or intentionally ignored.';
