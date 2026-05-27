package catalog

import (
	"time"

	"github.com/mohit838/olario-platform-backend/internal/domain/tenant"
)

// ProductID is the identifier for a grocery product.
type ProductID int64

// Product is a sellable grocery item.
// StockQuantity is the current snapshot; inventory movements keep the stock
// history so product reads stay fast.
type Product struct {
	ID            ProductID
	TenantID      tenant.ID
	CategoryID    *CategoryID
	SKU           string
	Barcode       string
	Name          string
	Description   string
	Unit          string
	Price         Money
	CostPrice     *Money
	StockQuantity int
	AlertAtStock  int
	IsActive      bool
	IsDeleted     bool
	IsFeatured    bool
	IsDiscounted  bool
	DiscountRate  float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Money stores fixed-point currency values as cents.
// The database currently stores NUMERIC for readability; repository adapters
// can convert between NUMERIC and this domain representation later.
type Money int64
