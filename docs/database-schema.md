# Database Schema Notes

This schema is for a multi-tenant grocery backend. Every business table carries
`tenant_id` so data can be scoped to one shop or organization from the start.

## Tenant Boundary

- `tenants` is the root ownership table.
- Tenant-owned tables reference `tenants(id)` and use tenant-aware unique keys.
- Composite foreign keys such as `(tenant_id, product_id)` prevent linking one
  tenant's order to another tenant's product.

## Access Foundation

- `roles` keeps simple tenant-scoped roles such as owner, manager, cashier, or
  staff.
- `users` stores tenant-scoped application users.
- `users.password_hash` must contain only a password hash, never a plain text
  password.
- Fine-grained permissions are intentionally not included yet. Add them later
  only when role names are no longer enough.

## Catalog

- `categories` groups products and supports optional parent categories.
- `categories.code_prefix` is used to generate manual product codes such as
  `BEV-0001`.
- `products` stores grocery sellable items, price, stock, SKU, barcode, and
  simple soft-delete flags.
- `products.product_code` is the manual fallback when barcode scanning is not
  available.
- `stock_quantity` is the current stock snapshot for fast reads.
- `inventory_movements` is the stock history/audit trail.
- `vendors` stores suppliers or salesmen who provide products.
- `products.alert_at_stock` defaults to `5`, and users can change it per
  product.
- `low_stock_notifications` records low-stock events when stock reaches the
  product threshold.

## Checkout

- `customers` stores optional customer information for receipts and future
  loyalty features.
- `customers.email_opt_in` controls whether email receipts/messages should be
  sent later.
- `loyalty_transactions` is the points ledger.
- `orders` replaces the old `sales` table name because it is clearer for API
  design and can support draft, paid, cancelled, or refunded workflows.
- `order_items` stores the order lines.
- `order_items.product_name` is a snapshot so old receipts remain readable even
  if the product is renamed or removed.

## SaaS And Access Control

- `superadmins` stores platform-level admins outside tenant scope.
- `tenant_invitations` supports superadmin invite-only tenant onboarding.
- `subscription_plans` supports medium/business plan rules.
- `tenant_subscriptions` stores trial and paid subscription state.
- `tenant_deactivation_requests` lets a tenant admin request account
  deactivation for superadmin review.
- `ui_menus` stores flexible frontend menus and submenus for tenant and
  superadmin dashboards.
- `permissions` stores action-level access rules and can link each permission
  to a `ui_menus` row.
- `role_permissions` maps tenant role access to tenant permissions.
- `superadmin_permissions` maps platform superadmin access to superadmin
  permissions until full platform roles are added.
- `report_requests` queues monthly reports to run later, outside peak hours.
- `audit_logs` stores important business/security events in Postgres first.

## Tables Not Added Yet

- Email delivery tables are postponed until the email step.
- MinIO object metadata tables are postponed until storage behavior is clear.
- Kafka/outbox tables are postponed because v1 does not need async event
  streaming yet.

## Update Timestamp Rule

`updated_at` is stored on mutable tables, but automatic update triggers are not
added yet. For v1, the application layer should set `updated_at` when updating
rows. A shared Postgres trigger can be added later if repeated update logic
becomes noisy.
