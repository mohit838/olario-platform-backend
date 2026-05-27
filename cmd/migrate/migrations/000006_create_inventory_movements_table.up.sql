CREATE TABLE inventory_movements (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  product_id BIGINT NOT NULL,
  movement_type VARCHAR(30) NOT NULL,
  quantity INT NOT NULL CHECK (quantity > 0),
  note TEXT,
  created_by BIGINT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  FOREIGN KEY (tenant_id, product_id) REFERENCES products(tenant_id, id) ON DELETE CASCADE,
  FOREIGN KEY (tenant_id, created_by) REFERENCES users(tenant_id, id) ON DELETE SET NULL (created_by),
  CHECK (movement_type IN ('stock_in', 'stock_out', 'adjustment', 'sale', 'return'))
);

CREATE INDEX idx_inventory_movements_tenant_id ON inventory_movements(tenant_id);
CREATE INDEX idx_inventory_movements_product_id ON inventory_movements(product_id);
CREATE INDEX idx_inventory_movements_created_at ON inventory_movements(created_at);

COMMENT ON TABLE inventory_movements IS 'Audit trail for stock changes without coupling stock history to product updates.';
COMMENT ON COLUMN inventory_movements.tenant_id IS 'Scopes the inventory movement to one tenant.';
COMMENT ON COLUMN inventory_movements.product_id IS 'Product whose stock changed.';
COMMENT ON COLUMN inventory_movements.movement_type IS 'Reason for stock movement, such as stock_in, sale, return, or adjustment.';
COMMENT ON COLUMN inventory_movements.quantity IS 'Positive quantity moved; movement_type explains the direction or reason.';
COMMENT ON COLUMN inventory_movements.note IS 'Optional human note explaining the stock change.';
COMMENT ON COLUMN inventory_movements.created_by IS 'Optional user who recorded the stock movement.';
