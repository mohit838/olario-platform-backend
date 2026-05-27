CREATE TABLE roles (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  name VARCHAR(50) NOT NULL,
  code VARCHAR(50),
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (tenant_id, id),
  UNIQUE (tenant_id, name),
  UNIQUE (tenant_id, code)
);

CREATE INDEX idx_roles_tenant_id ON roles(tenant_id);
CREATE INDEX idx_roles_is_active ON roles(is_active);

COMMENT ON TABLE roles IS 'Tenant-scoped user roles such as owner, manager, cashier, or staff.';
COMMENT ON COLUMN roles.tenant_id IS 'Scopes this role to one tenant and prevents cross-tenant reuse.';
COMMENT ON COLUMN roles.name IS 'Human-readable role name shown in admin screens.';
COMMENT ON COLUMN roles.code IS 'Optional stable role code for application logic; display names can change.';
COMMENT ON COLUMN roles.is_active IS 'Disables a role without deleting users or history.';
