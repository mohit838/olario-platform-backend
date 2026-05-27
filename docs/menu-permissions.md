# Menu And Permission Seeding Guide

This project stores frontend navigation and permissions in the database so the
UI can be flexible later.

## Tables

- `ui_menus`: menu and submenu tree.
- `permissions`: action permissions such as `view`, `manage`, `review`, or
  `use`.
- `superadmin_permissions`: permissions assigned to platform superadmin
  accounts.
- `role_permissions`: permissions assigned to tenant roles such as owner or
  employee.

## Menu Tree Rule

Top-level menus have `parent_id = NULL`.

Submenus point to a parent menu:

```text
sa.tenants
  sa.tenants.list
  sa.tenants.invitations
  sa.tenants.deactivation
```

Frontend code should read menus ordered by `parent_id` and `sort_order`, then
hide any menu that has no matching permission for the current user.

## Permission Rule

Each permission has:

```text
code       stable backend/frontend permission key
scope      superadmin or tenant
menu_key   menu/submenu this permission belongs to
action_key action inside that menu
```

Examples:

```text
superadmin.tenants.view
superadmin.invitations.manage
products.view
orders.manage
revenue.view
```

## Seeded Access

The seed command creates:

- one superadmin account
- one tenant
- one tenant owner/admin
- two tenant employees

The seeded superadmin receives every `superadmin` permission.

The seeded tenant owner receives every `tenant` permission.

The seeded employees receive a small cashier-friendly set:

- dashboard view
- product view
- inventory view
- POS use
- order view/manage
- customer view/manage
- low-stock view

Employees do not receive revenue or employee-management permissions by default.
That lets the frontend hide sensitive revenue and admin menus for normal staff.

## Example Read Queries

Superadmin menu permissions:

```sql
SELECT
  m.parent_id,
  m.menu_key,
  m.title,
  m.route_path,
  m.icon_name,
  m.sort_order,
  p.code AS permission_code,
  p.action_key
FROM superadmins s
JOIN superadmin_permissions sp ON sp.superadmin_id = s.id
JOIN permissions p ON p.id = sp.permission_id
JOIN ui_menus m ON m.id = p.ui_menu_id
WHERE s.email = 'superadmin@example.test'
ORDER BY COALESCE(m.parent_id, m.id), m.sort_order, p.action_key;
```

Tenant role menu permissions:

```sql
SELECT
  r.code AS role_code,
  m.parent_id,
  m.menu_key,
  m.title,
  m.route_path,
  p.code AS permission_code,
  p.action_key
FROM tenants t
JOIN roles r ON r.tenant_id = t.id
JOIN role_permissions rp ON rp.tenant_id = r.tenant_id AND rp.role_id = r.id
JOIN permissions p ON p.id = rp.permission_id
JOIN ui_menus m ON m.id = p.ui_menu_id
WHERE t.slug = 'seed-grocery'
ORDER BY r.code, COALESCE(m.parent_id, m.id), m.sort_order, p.action_key;
```

## How To Add A New Menu Later

1. Add a new row to `ui_menus` in the seeder.
2. If it is a submenu, set `parent_id` with a lookup to the parent menu key.
3. Add one or more `permissions` rows for that menu.
4. Assign those permissions to `superadmin_permissions` or `role_permissions`.
5. Run migrations if schema changed, then run `make seed`.

## Swagger Note

The current Swagger endpoint is a small starter OpenAPI document. When we move
to annotation-generated Swagger, run:

```bash
make swagger-generate
```

after route or DTO changes.
