package inventory

import (
	"time"

	"github.com/mohit838/olario-platform-backend/internal/domain/catalog"
	"github.com/mohit838/olario-platform-backend/internal/domain/identity"
	"github.com/mohit838/olario-platform-backend/internal/domain/tenant"
)

// MovementID is the identifier for a stock movement event.
type MovementID int64

// MovementType explains why stock changed.
type MovementType string

const (
	MovementTypeStockIn    MovementType = "stock_in"
	MovementTypeStockOut   MovementType = "stock_out"
	MovementTypeAdjustment MovementType = "adjustment"
	MovementTypeSale       MovementType = "sale"
	MovementTypeReturn     MovementType = "return"
)

// Movement records a stock change for audit and debugging.
// Quantity is always positive; MovementType explains the reason for the change.
type Movement struct {
	ID        MovementID
	TenantID  tenant.ID
	ProductID catalog.ProductID
	Type      MovementType
	Quantity  int
	Note      string
	CreatedBy *identity.UserID
	CreatedAt time.Time
}
