package customer

import (
	"time"

	"github.com/mohit838/olario-platform-backend/internal/domain/tenant"
)

// ID is the identifier for a grocery customer.
type ID int64

// Customer stores optional buyer information for orders and future loyalty
// features. Orders can still exist without a customer.
type Customer struct {
	ID        ID
	TenantID  tenant.ID
	Name      string
	Email     string
	Phone     string
	Address   string
	Notes     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
