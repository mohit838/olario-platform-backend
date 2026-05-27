DROP TABLE IF EXISTS superadmin_permissions;

ALTER TABLE permissions
  DROP COLUMN IF EXISTS ui_menu_id,
  DROP COLUMN IF EXISTS scope;

DROP TABLE IF EXISTS ui_menus;
