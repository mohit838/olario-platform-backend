package seeders

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/mohit838/olario-platform-backend/internal/security/token"
)

const (
	DevInvitationEmail = "owner@example.test"
	DevInvitationToken = "dev-invite-token"
	DevTenantSlug      = "seed-grocery"
	DevSeedPassword    = "12345678"
	DevSuperadminEmail = "superadmin@example.test"
	DevAdminEmail      = "admin@seed-grocery.test"
	DevEmployeeOne     = "employee1@seed-grocery.test"
	DevEmployeeTwo     = "employee2@seed-grocery.test"
)

// Run executes every idempotent seeder in dependency order.
func Run(ctx context.Context, db *sql.DB) error {
	if err := seedPlans(ctx, db); err != nil {
		return err
	}
	if err := seedSuperadmin(ctx, db); err != nil {
		return err
	}
	if err := seedPermissions(ctx, db); err != nil {
		return err
	}
	if err := seedInvitation(ctx, db); err != nil {
		return err
	}
	if err := seedDemoTenant(ctx, db); err != nil {
		return err
	}
	return nil
}

func seedPlans(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO subscription_plans (code, name, billing_period, price, max_employees)
		VALUES
			('medium', 'Medium', 'monthly', 0, 3),
			('business', 'Business', 'monthly', 0, NULL)
		ON CONFLICT (code) DO UPDATE
		SET name = EXCLUDED.name,
			billing_period = EXCLUDED.billing_period,
			price = EXCLUDED.price,
			max_employees = EXCLUDED.max_employees,
			updated_at = NOW()
	`)
	if err != nil {
		return fmt.Errorf("seed plans: %w", err)
	}
	return nil
}

func seedPermissions(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO ui_menus (scope, menu_key, title, route_path, icon_name, sort_order)
		VALUES
			('superadmin', 'sa.dashboard', 'Dashboard', '/superadmin', 'layout-dashboard', 10),
			('superadmin', 'sa.tenants', 'Tenants', NULL, 'building-2', 20),
			('superadmin', 'sa.billing', 'Billing', NULL, 'credit-card', 30),
			('superadmin', 'sa.identity', 'Platform Users', NULL, 'shield', 40),
			('superadmin', 'sa.reports', 'Reports', NULL, 'file-bar-chart', 50),
			('superadmin', 'sa.system', 'System', NULL, 'settings', 60),
			('tenant', 'tenant.dashboard', 'Dashboard', '/app', 'layout-dashboard', 10),
			('tenant', 'tenant.catalog', 'Catalog', NULL, 'package', 20),
			('tenant', 'tenant.checkout', 'Checkout', NULL, 'shopping-cart', 30),
			('tenant', 'tenant.people', 'People', NULL, 'users', 40),
			('tenant', 'tenant.reports', 'Reports', NULL, 'file-bar-chart', 50),
			('tenant', 'tenant.settings', 'Settings', NULL, 'settings', 60)
		ON CONFLICT (menu_key) DO UPDATE
		SET scope = EXCLUDED.scope,
			title = EXCLUDED.title,
			route_path = EXCLUDED.route_path,
			icon_name = EXCLUDED.icon_name,
			sort_order = EXCLUDED.sort_order,
			is_active = true,
			updated_at = NOW()
	`)
	if err != nil {
		return fmt.Errorf("seed root menus: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO ui_menus (parent_id, scope, menu_key, title, route_path, icon_name, sort_order)
		VALUES
			((SELECT id FROM ui_menus WHERE menu_key = 'sa.tenants'), 'superadmin', 'sa.tenants.list', 'Tenant Accounts', '/superadmin/tenants', 'store', 10),
			((SELECT id FROM ui_menus WHERE menu_key = 'sa.tenants'), 'superadmin', 'sa.tenants.invitations', 'Invitations', '/superadmin/tenant-invitations', 'mail-plus', 20),
			((SELECT id FROM ui_menus WHERE menu_key = 'sa.tenants'), 'superadmin', 'sa.tenants.deactivation', 'Deactivation Requests', '/superadmin/deactivation-requests', 'power', 30),
			((SELECT id FROM ui_menus WHERE menu_key = 'sa.billing'), 'superadmin', 'sa.billing.plans', 'Plans', '/superadmin/plans', 'badge-dollar-sign', 10),
			((SELECT id FROM ui_menus WHERE menu_key = 'sa.billing'), 'superadmin', 'sa.billing.subscriptions', 'Subscriptions', '/superadmin/subscriptions', 'receipt', 20),
			((SELECT id FROM ui_menus WHERE menu_key = 'sa.identity'), 'superadmin', 'sa.identity.superadmins', 'Superadmins', '/superadmin/superadmins', 'shield-check', 10),
			((SELECT id FROM ui_menus WHERE menu_key = 'sa.reports'), 'superadmin', 'sa.reports.queue', 'Report Queue', '/superadmin/report-requests', 'clock', 10),
			((SELECT id FROM ui_menus WHERE menu_key = 'sa.system'), 'superadmin', 'sa.system.audit', 'Audit Logs', '/superadmin/audit-logs', 'scroll-text', 10),
			((SELECT id FROM ui_menus WHERE menu_key = 'sa.system'), 'superadmin', 'sa.system.settings', 'Platform Settings', '/superadmin/settings', 'sliders-horizontal', 20),
			((SELECT id FROM ui_menus WHERE menu_key = 'tenant.catalog'), 'tenant', 'tenant.catalog.categories', 'Categories', '/app/categories', 'tags', 10),
			((SELECT id FROM ui_menus WHERE menu_key = 'tenant.catalog'), 'tenant', 'tenant.catalog.products', 'Products', '/app/products', 'package-search', 20),
			((SELECT id FROM ui_menus WHERE menu_key = 'tenant.catalog'), 'tenant', 'tenant.catalog.inventory', 'Inventory', '/app/inventory', 'warehouse', 30),
			((SELECT id FROM ui_menus WHERE menu_key = 'tenant.catalog'), 'tenant', 'tenant.catalog.vendors', 'Vendors', '/app/vendors', 'truck', 40),
			((SELECT id FROM ui_menus WHERE menu_key = 'tenant.checkout'), 'tenant', 'tenant.checkout.pos', 'POS', '/app/pos', 'scan-barcode', 10),
			((SELECT id FROM ui_menus WHERE menu_key = 'tenant.checkout'), 'tenant', 'tenant.checkout.orders', 'Orders', '/app/orders', 'receipt-text', 20),
			((SELECT id FROM ui_menus WHERE menu_key = 'tenant.people'), 'tenant', 'tenant.people.customers', 'Customers', '/app/customers', 'contact', 10),
			((SELECT id FROM ui_menus WHERE menu_key = 'tenant.people'), 'tenant', 'tenant.people.employees', 'Employees', '/app/employees', 'user-cog', 20),
			((SELECT id FROM ui_menus WHERE menu_key = 'tenant.reports'), 'tenant', 'tenant.reports.revenue', 'Revenue', '/app/reports/revenue', 'chart-column', 10),
			((SELECT id FROM ui_menus WHERE menu_key = 'tenant.reports'), 'tenant', 'tenant.reports.low_stock', 'Low Stock', '/app/reports/low-stock', 'triangle-alert', 20),
			((SELECT id FROM ui_menus WHERE menu_key = 'tenant.settings'), 'tenant', 'tenant.settings.permissions', 'Roles & Permissions', '/app/settings/permissions', 'key-round', 10)
		ON CONFLICT (menu_key) DO UPDATE
		SET parent_id = EXCLUDED.parent_id,
			scope = EXCLUDED.scope,
			title = EXCLUDED.title,
			route_path = EXCLUDED.route_path,
			icon_name = EXCLUDED.icon_name,
			sort_order = EXCLUDED.sort_order,
			is_active = true,
			updated_at = NOW()
	`)
	if err != nil {
		return fmt.Errorf("seed submenus: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO permissions (code, menu_key, action_key, description, scope, ui_menu_id)
		VALUES
			('superadmin.dashboard.view', 'sa.dashboard', 'view', 'View platform dashboard', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.dashboard')),
			('superadmin.tenants.view', 'sa.tenants.list', 'view', 'View tenant accounts', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.tenants.list')),
			('superadmin.tenants.manage', 'sa.tenants.list', 'manage', 'Create, update, suspend, or reactivate tenant accounts', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.tenants.list')),
			('superadmin.invitations.view', 'sa.tenants.invitations', 'view', 'View tenant invitations', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.tenants.invitations')),
			('superadmin.invitations.manage', 'sa.tenants.invitations', 'manage', 'Create and revoke tenant invitations', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.tenants.invitations')),
			('superadmin.deactivation.view', 'sa.tenants.deactivation', 'view', 'View tenant deactivation requests', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.tenants.deactivation')),
			('superadmin.deactivation.review', 'sa.tenants.deactivation', 'review', 'Approve or reject tenant deactivation requests', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.tenants.deactivation')),
			('superadmin.plans.view', 'sa.billing.plans', 'view', 'View subscription plans', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.billing.plans')),
			('superadmin.plans.manage', 'sa.billing.plans', 'manage', 'Create and update subscription plans', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.billing.plans')),
			('superadmin.subscriptions.view', 'sa.billing.subscriptions', 'view', 'View tenant subscriptions', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.billing.subscriptions')),
			('superadmin.subscriptions.manage', 'sa.billing.subscriptions', 'manage', 'Adjust subscription status or billing period', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.billing.subscriptions')),
			('superadmin.superadmins.view', 'sa.identity.superadmins', 'view', 'View platform superadmin accounts', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.identity.superadmins')),
			('superadmin.superadmins.manage', 'sa.identity.superadmins', 'manage', 'Invite, activate, or deactivate superadmin accounts', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.identity.superadmins')),
			('superadmin.reports.view', 'sa.reports.queue', 'view', 'View report generation queue', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.reports.queue')),
			('superadmin.reports.manage', 'sa.reports.queue', 'manage', 'Retry or cancel report generation requests', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.reports.queue')),
			('superadmin.audit.view', 'sa.system.audit', 'view', 'View platform and tenant audit logs', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.system.audit')),
			('superadmin.settings.manage', 'sa.system.settings', 'manage', 'Manage platform settings', 'superadmin', (SELECT id FROM ui_menus WHERE menu_key = 'sa.system.settings')),
			('dashboard.view', 'tenant.dashboard', 'view', 'View tenant dashboard', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.dashboard')),
			('categories.view', 'tenant.catalog.categories', 'view', 'View categories', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.catalog.categories')),
			('categories.manage', 'tenant.catalog.categories', 'manage', 'Create and update categories', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.catalog.categories')),
			('products.view', 'tenant.catalog.products', 'view', 'View product catalog', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.catalog.products')),
			('products.manage', 'tenant.catalog.products', 'manage', 'Create and update products', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.catalog.products')),
			('inventory.view', 'tenant.catalog.inventory', 'view', 'View stock and movements', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.catalog.inventory')),
			('inventory.manage', 'tenant.catalog.inventory', 'manage', 'Adjust stock and thresholds', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.catalog.inventory')),
			('vendors.view', 'tenant.catalog.vendors', 'view', 'View vendors and salesmen', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.catalog.vendors')),
			('vendors.manage', 'tenant.catalog.vendors', 'manage', 'Create and update vendors', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.catalog.vendors')),
			('pos.use', 'tenant.checkout.pos', 'use', 'Use barcode/manual product checkout', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.checkout.pos')),
			('orders.view', 'tenant.checkout.orders', 'view', 'View orders and invoices', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.checkout.orders')),
			('orders.manage', 'tenant.checkout.orders', 'manage', 'Create, cancel, or refund orders', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.checkout.orders')),
			('customers.view', 'tenant.people.customers', 'view', 'View customer profiles', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.people.customers')),
			('customers.manage', 'tenant.people.customers', 'manage', 'Create and update customer profiles', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.people.customers')),
			('employees.view', 'tenant.people.employees', 'view', 'View employees', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.people.employees')),
			('employees.manage', 'tenant.people.employees', 'manage', 'Create, update, or deactivate employees', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.people.employees')),
			('revenue.view', 'tenant.reports.revenue', 'view', 'View revenue reports', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.reports.revenue')),
			('low_stock.view', 'tenant.reports.low_stock', 'view', 'View low-stock notifications', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.reports.low_stock')),
			('roles_permissions.manage', 'tenant.settings.permissions', 'manage', 'Manage role menu visibility and actions', 'tenant', (SELECT id FROM ui_menus WHERE menu_key = 'tenant.settings.permissions'))
		ON CONFLICT (code) DO UPDATE
		SET menu_key = EXCLUDED.menu_key,
			action_key = EXCLUDED.action_key,
			description = EXCLUDED.description,
			scope = EXCLUDED.scope,
			ui_menu_id = EXCLUDED.ui_menu_id
	`)
	if err != nil {
		return fmt.Errorf("seed permissions: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO superadmin_permissions (superadmin_id, permission_id)
		SELECT s.id, p.id
		FROM superadmins s
		CROSS JOIN permissions p
		WHERE s.email = $1
			AND p.scope = 'superadmin'
		ON CONFLICT (superadmin_id, permission_id) DO NOTHING
	`, DevSuperadminEmail)
	if err != nil {
		return fmt.Errorf("assign superadmin permissions: %w", err)
	}
	return nil
}

func seedInvitation(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO tenant_invitations (email, token_hash, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (email, status) DO UPDATE
		SET token_hash = EXCLUDED.token_hash,
			expires_at = EXCLUDED.expires_at,
			created_at = NOW()
	`, DevInvitationEmail, token.HashInvitationToken(DevInvitationToken), time.Now().UTC().Add(7*24*time.Hour))
	if err != nil {
		return fmt.Errorf("seed invitation: %w", err)
	}
	return nil
}

func seedSuperadmin(ctx context.Context, db *sql.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(DevSeedPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash superadmin seed password: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO superadmins (name, email, password_hash, two_fa_enabled)
		VALUES ('Seed Superadmin', $1, $2, false)
		ON CONFLICT (email) DO UPDATE
		SET password_hash = EXCLUDED.password_hash,
			two_fa_enabled = false,
			is_active = true,
			updated_at = NOW()
	`, DevSuperadminEmail, string(passwordHash))
	if err != nil {
		return fmt.Errorf("seed superadmin: %w", err)
	}
	return nil
}

func seedDemoTenant(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin demo tenant seed tx: %w", err)
	}
	defer tx.Rollback()

	var planID int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM subscription_plans WHERE code = 'medium'`).Scan(&planID); err != nil {
		return fmt.Errorf("find medium plan: %w", err)
	}

	var tenantID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO tenants (name, slug, status, trial_ends_at)
		VALUES ('Seed Grocery', $1, 'trial', $2)
		ON CONFLICT (slug) DO UPDATE SET updated_at = NOW()
		RETURNING id
	`, DevTenantSlug, time.Now().UTC().AddDate(0, 3, 0)).Scan(&tenantID); err != nil {
		return fmt.Errorf("seed tenant: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO tenant_subscriptions (tenant_id, plan_id, status, starts_at, trial_ends_at)
		SELECT $1, $2, 'trialing', NOW(), $3
		WHERE NOT EXISTS (SELECT 1 FROM tenant_subscriptions WHERE tenant_id = $1)
	`, tenantID, planID, time.Now().UTC().AddDate(0, 3, 0)); err != nil {
		return fmt.Errorf("seed tenant subscription: %w", err)
	}

	var ownerRoleID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO roles (tenant_id, name, code)
		VALUES ($1, 'Owner', 'owner')
		ON CONFLICT (tenant_id, code) DO UPDATE SET updated_at = NOW()
		RETURNING id
	`, tenantID).Scan(&ownerRoleID); err != nil {
		return fmt.Errorf("seed owner role: %w", err)
	}

	var employeeRoleID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO roles (tenant_id, name, code)
		VALUES ($1, 'Employee', 'employee')
		ON CONFLICT (tenant_id, code) DO UPDATE SET updated_at = NOW()
		RETURNING id
	`, tenantID).Scan(&employeeRoleID); err != nil {
		return fmt.Errorf("seed employee role: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(DevSeedPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash seed password: %w", err)
	}

	if err := upsertTenantUser(ctx, tx, tenantID, ownerRoleID, "Seed Admin", DevAdminEmail, string(passwordHash)); err != nil {
		return err
	}
	if err := upsertTenantUser(ctx, tx, tenantID, employeeRoleID, "Seed Employee One", DevEmployeeOne, string(passwordHash)); err != nil {
		return err
	}
	if err := upsertTenantUser(ctx, tx, tenantID, employeeRoleID, "Seed Employee Two", DevEmployeeTwo, string(passwordHash)); err != nil {
		return err
	}
	if err := assignTenantRolePermissions(ctx, tx, tenantID, ownerRoleID, []string{"tenant_all"}); err != nil {
		return err
	}
	if err := assignTenantRolePermissions(ctx, tx, tenantID, employeeRoleID, []string{
		"dashboard.view",
		"products.view",
		"inventory.view",
		"pos.use",
		"orders.view",
		"orders.manage",
		"customers.view",
		"customers.manage",
		"low_stock.view",
	}); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit demo tenant seed: %w", err)
	}
	return nil
}

func assignTenantRolePermissions(ctx context.Context, tx *sql.Tx, tenantID, roleID int64, permissionCodes []string) error {
	if len(permissionCodes) == 1 && permissionCodes[0] == "tenant_all" {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO role_permissions (tenant_id, role_id, permission_id)
			SELECT $1, $2, p.id
			FROM permissions p
			WHERE p.scope = 'tenant'
			ON CONFLICT (tenant_id, role_id, permission_id) DO NOTHING
		`, tenantID, roleID)
		if err != nil {
			return fmt.Errorf("assign all tenant permissions: %w", err)
		}
		return nil
	}

	for _, code := range permissionCodes {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO role_permissions (tenant_id, role_id, permission_id)
			SELECT $1, $2, p.id
			FROM permissions p
			WHERE p.code = $3
				AND p.scope = 'tenant'
			ON CONFLICT (tenant_id, role_id, permission_id) DO NOTHING
		`, tenantID, roleID, code)
		if err != nil {
			return fmt.Errorf("assign tenant permission %q: %w", code, err)
		}
	}
	return nil
}

func upsertTenantUser(ctx context.Context, tx *sql.Tx, tenantID, roleID int64, name, email, passwordHash string) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO users (tenant_id, role_id, name, email, password_hash, two_fa_enabled)
		VALUES ($1, $2, $3, $4, $5, false)
		ON CONFLICT (tenant_id, email) DO UPDATE
		SET role_id = EXCLUDED.role_id,
			name = EXCLUDED.name,
			password_hash = EXCLUDED.password_hash,
			two_fa_enabled = false,
			is_active = true,
			updated_at = NOW()
	`, tenantID, roleID, name, email, passwordHash)
	if err != nil {
		return fmt.Errorf("seed tenant user %q: %w", email, err)
	}
	return nil
}
