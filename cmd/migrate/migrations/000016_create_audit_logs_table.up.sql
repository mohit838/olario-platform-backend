CREATE TABLE audit_logs (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT REFERENCES tenants(id) ON DELETE SET NULL,
  actor_user_id BIGINT,
  action VARCHAR(120) NOT NULL,
  entity_type VARCHAR(120) NOT NULL,
  entity_id BIGINT,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  ip_address INET,
  user_agent TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_actor_user_id ON audit_logs(actor_user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_metadata ON audit_logs USING GIN(metadata);

COMMENT ON TABLE audit_logs IS 'Initial audit log storage in Postgres for important business and security events.';
COMMENT ON COLUMN audit_logs.action IS 'Event action such as dev.full_circle, order.created, employee.deactivated, or report.requested.';
COMMENT ON COLUMN audit_logs.entity_type IS 'Business object type affected by the action.';
COMMENT ON COLUMN audit_logs.metadata IS 'Additional structured event details that do not need dedicated columns yet.';
