package identity

import (
	"time"

	"github.com/mohit838/olario-platform-backend/internal/domain/tenant"
)

// RoleID is the identifier for a tenant-scoped role.
type RoleID int64

// Role describes a user's high-level responsibility inside one tenant.
// Fine-grained permissions are intentionally delayed until roles are not enough.
type Role struct {
	ID        RoleID
	TenantID  tenant.ID
	Name      string
	Code      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
