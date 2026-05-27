package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	devapp "github.com/mohit838/olario-platform-backend/internal/application/dev"
)

// DevRepository is a teaching repository for one complete local API flow.
// Production repositories should be split by domain area once real endpoints are
// added; this one intentionally stays together so the end-to-end flow is easy
// to study.
type DevRepository struct {
	db *sql.DB
}

func NewDevRepository(db *sql.DB) *DevRepository {
	return &DevRepository{db: db}
}

func (r *DevRepository) CreateFullCircleDemo(ctx context.Context) (devapp.FullCircleResult, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	suffix := now.Format("20060102150405")
	tenantSlug := "demo-grocery-" + suffix
	orderNumber := "INV-" + suffix
	productCode := "BEV-0001-" + suffix
	barcode := "880" + now.Format("0102150405")

	var tenantID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO tenants (name, slug, status, trial_ends_at)
		VALUES ($1, $2, 'trial', $3)
		RETURNING id
	`, "Demo Grocery "+suffix, tenantSlug, now.AddDate(0, 3, 0)).Scan(&tenantID); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert tenant: %w", err)
	}

	var planID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO subscription_plans (code, name, billing_period, price, max_employees)
		VALUES ('medium', 'Medium', 'monthly', 0, 3)
		ON CONFLICT (code) DO UPDATE SET updated_at = NOW()
		RETURNING id
	`).Scan(&planID); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("upsert subscription plan: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO tenant_subscriptions (tenant_id, plan_id, status, starts_at, trial_ends_at)
		VALUES ($1, $2, 'trialing', NOW(), $3)
	`, tenantID, planID, now.AddDate(0, 3, 0)); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert tenant subscription: %w", err)
	}

	var roleID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO roles (tenant_id, name, code)
		VALUES ($1, 'Cashier', 'cashier')
		RETURNING id
	`, tenantID).Scan(&roleID); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert role: %w", err)
	}

	var employeeID int64
	employeeName := "Demo Cashier"
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO users (tenant_id, role_id, name, email, password_hash, two_fa_enabled)
		VALUES ($1, $2, $3, $4, 'dev-password-hash-placeholder', true)
		RETURNING id
	`, tenantID, roleID, employeeName, "cashier+"+suffix+"@example.test").Scan(&employeeID); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert employee: %w", err)
	}

	var categoryID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO categories (tenant_id, name, code_prefix, description)
		VALUES ($1, 'Beverages', 'BEV', 'Drinks and bottled grocery items')
		RETURNING id
	`, tenantID).Scan(&categoryID); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert category: %w", err)
	}

	var vendorID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO vendors (tenant_id, name, contact_person, phone)
		VALUES ($1, 'Demo Vendor', 'Demo Salesman', '+880000000000')
		RETURNING id
	`, tenantID).Scan(&vendorID); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert vendor: %w", err)
	}

	var productID int64
	productName := "Demo Juice"
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO products (
			tenant_id, category_id, primary_vendor_id, product_code, barcode, name,
			unit, price, cost_price, stock_quantity, alert_at_stock
		)
		VALUES ($1, $2, $3, $4, $5, $6, 'pcs', 120.00, 80.00, 6, 5)
		RETURNING id
	`, tenantID, categoryID, vendorID, productCode, barcode, productName).Scan(&productID); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert product: %w", err)
	}

	var customerID int64
	customerName := "Demo Customer"
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO customers (tenant_id, name, phone, email, email_opt_in)
		VALUES ($1, $2, '+880111111111', $3, true)
		RETURNING id
	`, tenantID, customerName, "customer+"+suffix+"@example.test").Scan(&customerID); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert customer: %w", err)
	}

	quantity := 2
	totalAmount := "240.00"
	loyaltyPoints := 2

	var orderID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO orders (
			tenant_id, order_number, customer_id, created_by, status,
			subtotal_amount, total_amount, paid_amount, loyalty_points_earned, payment_method
		)
		VALUES ($1, $2, $3, $4, 'paid', 240.00, 240.00, 240.00, $5, 'cash')
		RETURNING id
	`, tenantID, orderNumber, customerID, employeeID, loyaltyPoints).Scan(&orderID); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert order: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO order_items (
			tenant_id, order_id, product_id, product_name, quantity, unit_price, line_total
		)
		VALUES ($1, $2, $3, $4, $5, 120.00, 240.00)
	`, tenantID, orderID, productID, productName, quantity); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert order item: %w", err)
	}

	stockAfterOrder := 4
	if _, err := tx.ExecContext(ctx, `
		UPDATE products
		SET stock_quantity = $1, updated_at = NOW()
		WHERE tenant_id = $2 AND id = $3
	`, stockAfterOrder, tenantID, productID); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("update stock: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO inventory_movements (tenant_id, product_id, movement_type, quantity, note, created_by)
		VALUES ($1, $2, 'sale', $3, 'Dev full-circle order reduced stock', $4)
	`, tenantID, productID, quantity, employeeID); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert inventory movement: %w", err)
	}

	lowStockCreated := stockAfterOrder <= 5
	if lowStockCreated {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO low_stock_notifications (tenant_id, product_id, threshold_quantity, current_quantity)
			VALUES ($1, $2, 5, $3)
		`, tenantID, productID, stockAfterOrder); err != nil {
			return devapp.FullCircleResult{}, fmt.Errorf("insert low-stock notification: %w", err)
		}
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE customers
		SET loyalty_points_balance = loyalty_points_balance + $1, updated_at = NOW()
		WHERE tenant_id = $2 AND id = $3
	`, loyaltyPoints, tenantID, customerID); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("update loyalty balance: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO loyalty_transactions (tenant_id, customer_id, order_id, points, reason)
		VALUES ($1, $2, $3, $4, 'order_purchase')
	`, tenantID, customerID, orderID, loyaltyPoints); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert loyalty transaction: %w", err)
	}

	metadata, err := json.Marshal(map[string]any{
		"order_number":        orderNumber,
		"product_code":        productCode,
		"barcode":             barcode,
		"low_stock_created":   lowStockCreated,
		"invoice_print_ready": true,
	})
	if err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("marshal audit metadata: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO audit_logs (tenant_id, actor_user_id, action, entity_type, entity_id, metadata)
		VALUES ($1, $2, 'dev.full_circle', 'orders', $3, $4)
	`, tenantID, employeeID, orderID, metadata); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("insert audit log: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return devapp.FullCircleResult{}, fmt.Errorf("commit tx: %w", err)
	}

	return devapp.FullCircleResult{
		TenantID:        tenantID,
		TenantSlug:      tenantSlug,
		EmployeeID:      employeeID,
		EmployeeName:    employeeName,
		CustomerID:      customerID,
		CustomerName:    customerName,
		ProductID:       productID,
		ProductName:     productName,
		ProductCode:     productCode,
		Barcode:         barcode,
		OrderID:         orderID,
		OrderNumber:     orderNumber,
		Quantity:        quantity,
		TotalAmount:     totalAmount,
		StockAfterOrder: stockAfterOrder,
		LowStockCreated: lowStockCreated,
		LoyaltyPoints:   loyaltyPoints,
	}, nil
}
