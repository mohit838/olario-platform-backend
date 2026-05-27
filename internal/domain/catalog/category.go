package catalog

import (
	"time"

	"github.com/mohit838/olario-platform-backend/internal/domain/tenant"
)

// CategoryID is the identifier for a product category.
type CategoryID int64

// Category groups grocery products for browsing, reporting, and management.
// ParentID allows simple nested categories without creating a separate tree
// system too early.
type Category struct {
	ID          CategoryID
	TenantID    tenant.ID
	ParentID    *CategoryID
	Name        string
	Description string
	SortOrder   int
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
