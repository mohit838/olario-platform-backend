package ordering

import (
	"time"

	"github.com/mohit838/olario-platform-backend/internal/domain/catalog"
	"github.com/mohit838/olario-platform-backend/internal/domain/customer"
	"github.com/mohit838/olario-platform-backend/internal/domain/identity"
	"github.com/mohit838/olario-platform-backend/internal/domain/tenant"
)

// OrderID is the identifier for a checkout order.
type OrderID int64

// Status describes the order lifecycle.
type Status string

const (
	StatusDraft     Status = "draft"
	StatusConfirmed Status = "confirmed"
	StatusPaid      Status = "paid"
	StatusCancelled Status = "cancelled"
	StatusRefunded  Status = "refunded"
)

// Order represents one grocery checkout record.
// Totals are snapshots of the checkout math so historical orders do not change
// when product prices change later.
type Order struct {
	ID             OrderID
	TenantID       tenant.ID
	OrderNumber    string
	CustomerID     *customer.ID
	CreatedBy      *identity.UserID
	Status         Status
	SubtotalAmount catalog.Money
	DiscountAmount catalog.Money
	TaxAmount      catalog.Money
	TotalAmount    catalog.Money
	PaidAmount     catalog.Money
	PaymentMethod  string
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
