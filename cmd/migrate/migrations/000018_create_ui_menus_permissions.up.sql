CREATE TABLE ui_menus (
  id BIGSERIAL PRIMARY KEY,
  parent_id BIGINT REFERENCES ui_menus(id) ON DELETE CASCADE,
  scope VARCHAR(30) NOT NULL CHECK (scope IN ('superadmin', 'tenant')),
  menu_key VARCHAR(120) NOT NULL UNIQUE,
  title VARCHAR(120) NOT NULL,
  route_path VARCHAR(200),
  icon_name VARCHAR(80),
  sort_order INTEGER NOT NULL DEFAULT 0,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE permissions
  ADD COLUMN scope VARCHAR(30) NOT NULL DEFAULT 'tenant' CHECK (scope IN ('superadmin', 'tenant')),
  ADD COLUMN ui_menu_id BIGINT REFERENCES ui_menus(id) ON DELETE SET NULL;

CREATE TABLE superadmin_permissions (
  id BIGSERIAL PRIMARY KEY,
  superadmin_id BIGINT NOT NULL REFERENCES superadmins(id) ON DELETE CASCADE,
  permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (superadmin_id, permission_id)
);

CREATE INDEX idx_ui_menus_parent_id ON ui_menus(parent_id);
CREATE INDEX idx_ui_menus_scope_sort ON ui_menus(scope, sort_order);
CREATE INDEX idx_permissions_scope ON permissions(scope);
CREATE INDEX idx_permissions_ui_menu_id ON permissions(ui_menu_id);
CREATE INDEX idx_superadmin_permissions_superadmin_id ON superadmin_permissions(superadmin_id);

COMMENT ON TABLE ui_menus IS 'Frontend navigation tree for tenant and superadmin dashboards. parent_id creates submenu nesting.';
COMMENT ON COLUMN ui_menus.scope IS 'Navigation owner: superadmin dashboard or tenant dashboard.';
COMMENT ON COLUMN ui_menus.menu_key IS 'Stable frontend key used by API responses and permission mapping.';
COMMENT ON COLUMN ui_menus.route_path IS 'Frontend route path; nullable for parent menu groups.';
COMMENT ON COLUMN ui_menus.icon_name IS 'Frontend icon name such as layout-dashboard, users, shopping-cart, or shield.';
COMMENT ON COLUMN ui_menus.sort_order IS 'Controls menu ordering inside the same parent.';
COMMENT ON COLUMN permissions.scope IS 'Permission owner: superadmin platform action or tenant business action.';
COMMENT ON COLUMN permissions.ui_menu_id IS 'Optional menu/submenu this permission belongs to for flexible frontend rendering.';
COMMENT ON TABLE superadmin_permissions IS 'Permission assignments for platform superadmin accounts until full platform roles are added.';
