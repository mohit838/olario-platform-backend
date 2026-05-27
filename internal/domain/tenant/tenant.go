package tenant

import "time"

// ID is the tenant identifier used to scope business data.
// Keeping it as a domain type makes tenant boundaries visible in function
// signatures instead of passing plain int64 values everywhere.
type ID int64

// Tenant represents one shop or organization using Olario.
// It is the root ownership concept for grocery data, users, orders, and stock.
type Tenant struct {
	ID        ID
	Name      string
	Slug      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
