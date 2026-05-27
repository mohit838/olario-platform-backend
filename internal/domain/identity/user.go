package identity

import (
	"time"

	"github.com/mohit838/olario-platform-backend/internal/domain/tenant"
)

// UserID is the identifier for an application user.
type UserID int64

// User represents a tenant-scoped person who can operate the grocery system.
// PasswordHash must contain a hash produced by an auth component, never a plain
// text password.
type User struct {
	ID           UserID
	TenantID     tenant.ID
	RoleID       *RoleID
	Name         string
	Email        string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
