package dev

import (
	"context"
	"fmt"
)

// Repository stores the durable part of the demo flow in Postgres.
type Repository interface {
	CreateFullCircleDemo(ctx context.Context) (FullCircleResult, error)
}

// Cache stores fast, non-authoritative demo counters in Redis.
type Cache interface {
	IncrementFullCircleRuns(ctx context.Context) (int64, error)
	RememberLastOrder(ctx context.Context, orderNumber string) error
}

// Service demonstrates the complete application flow.
// It exists as a teaching endpoint: HTTP calls the service, the service uses a
// repository for Postgres, and Redis is used for a fast counter/cache write.
type Service struct {
	repo  Repository
	cache Cache
}

// FullCircleResult is returned by the dev flow as an invoice-style response.
type FullCircleResult struct {
	TenantID          int64  `json:"tenant_id"`
	TenantSlug        string `json:"tenant_slug"`
	EmployeeID        int64  `json:"employee_id"`
	EmployeeName      string `json:"employee_name"`
	CustomerID        int64  `json:"customer_id"`
	CustomerName      string `json:"customer_name"`
	ProductID         int64  `json:"product_id"`
	ProductName       string `json:"product_name"`
	ProductCode       string `json:"product_code"`
	Barcode           string `json:"barcode"`
	OrderID           int64  `json:"order_id"`
	OrderNumber       string `json:"order_number"`
	Quantity          int    `json:"quantity"`
	TotalAmount       string `json:"total_amount"`
	StockAfterOrder   int    `json:"stock_after_order"`
	LowStockCreated   bool   `json:"low_stock_created"`
	LoyaltyPoints     int    `json:"loyalty_points"`
	RedisRunCount     int64  `json:"redis_run_count"`
	InvoicePrintReady bool   `json:"invoice_print_ready"`
}

// NewService creates the dev service after its infrastructure is available.
func NewService(repo Repository, cache Cache) *Service {
	return &Service{repo: repo, cache: cache}
}

// RunFullCircle creates seed data, builds one order, records audit data, and
// touches Redis so a new developer can see one complete request path.
func (s *Service) RunFullCircle(ctx context.Context) (FullCircleResult, error) {
	result, err := s.repo.CreateFullCircleDemo(ctx)
	if err != nil {
		return FullCircleResult{}, err
	}

	count, err := s.cache.IncrementFullCircleRuns(ctx)
	if err != nil {
		return FullCircleResult{}, fmt.Errorf("increment redis full-circle counter: %w", err)
	}

	if err := s.cache.RememberLastOrder(ctx, result.OrderNumber); err != nil {
		return FullCircleResult{}, fmt.Errorf("remember last order in redis: %w", err)
	}

	result.RedisRunCount = count
	result.InvoicePrintReady = true
	return result, nil
}
