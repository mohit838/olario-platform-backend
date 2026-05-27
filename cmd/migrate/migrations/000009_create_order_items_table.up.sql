CREATE TABLE order_items (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  order_id BIGINT NOT NULL,
  product_id BIGINT,
  product_name VARCHAR(160) NOT NULL,
  quantity INT NOT NULL CHECK (quantity > 0),
  unit_price NUMERIC(12, 2) NOT NULL CHECK (unit_price >= 0),
  discount_amount NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
  line_total NUMERIC(12, 2) NOT NULL CHECK (line_total >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  FOREIGN KEY (tenant_id, order_id) REFERENCES orders(tenant_id, id) ON DELETE CASCADE,
  FOREIGN KEY (tenant_id, product_id) REFERENCES products(tenant_id, id) ON DELETE SET NULL (product_id)
);

CREATE INDEX idx_order_items_tenant_id ON order_items(tenant_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

COMMENT ON TABLE order_items IS 'Line items for grocery orders.';
COMMENT ON COLUMN order_items.tenant_id IS 'Scopes the order item to one tenant.';
COMMENT ON COLUMN order_items.order_id IS 'Parent order for this line item.';
COMMENT ON COLUMN order_items.product_id IS 'Optional product reference; kept nullable so old orders survive product deletion.';
COMMENT ON COLUMN order_items.product_name IS 'Snapshot of product name at checkout time so old receipts stay readable.';
COMMENT ON COLUMN order_items.quantity IS 'Purchased quantity for this line item.';
COMMENT ON COLUMN order_items.unit_price IS 'Product unit price captured at checkout time.';
COMMENT ON COLUMN order_items.discount_amount IS 'Line-level discount amount.';
COMMENT ON COLUMN order_items.line_total IS 'Final line amount after quantity and discount.';
