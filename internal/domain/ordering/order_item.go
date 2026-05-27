package ordering

import (
	"time"

	"github.com/mohit838/olario-platform-backend/internal/domain/catalog"
	"github.com/mohit838/olario-platform-backend/internal/domain/tenant"
)

// OrderItemID is the identifier for an order line.
type OrderItemID int64

// OrderItem captures one product line inside an order.
// ProductName is copied from the product at checkout time so old receipts remain
// readable even if the product is renamed or deleted.
type OrderItem struct {
	ID             OrderItemID
	TenantID       tenant.ID
	OrderID        OrderID
	ProductID      *catalog.ProductID
	ProductName    string
	Quantity       int
	UnitPrice      catalog.Money
	DiscountAmount catalog.Money
	LineTotal      catalog.Money
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
