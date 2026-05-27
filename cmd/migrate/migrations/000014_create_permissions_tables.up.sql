CREATE TABLE permissions (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(120) NOT NULL UNIQUE,
  menu_key VARCHAR(80) NOT NULL,
  action_key VARCHAR(80) NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (menu_key, action_key)
);

CREATE TABLE role_permissions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  role_id BIGINT NOT NULL,
  permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (tenant_id, role_id, permission_id),
  FOREIGN KEY (tenant_id, role_id) REFERENCES roles(tenant_id, id) ON DELETE CASCADE
);

CREATE INDEX idx_permissions_menu_key ON permissions(menu_key);
CREATE INDEX idx_role_permissions_tenant_id ON role_permissions(tenant_id);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);

COMMENT ON TABLE permissions IS 'Menu/action permissions used to hide screens and block actions by role.';
COMMENT ON COLUMN permissions.menu_key IS 'UI or API area, such as revenue, products, orders, or employees.';
COMMENT ON COLUMN permissions.action_key IS 'Allowed action, such as view, create, update, delete, or export.';
COMMENT ON TABLE role_permissions IS 'Tenant-scoped assignment of permissions to roles.';
