package test

import (
	"context"
	"database/sql"

	"github.com/mohit838/olario-platform-backend/cmd/internal/dto"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListTests(ctx context.Context) ([]dto.TestResponse, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id::text, username, is_active, created_at, updated_at
		FROM test
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []dto.TestResponse
	for rows.Next() {
		test, err := scanTest(rows)
		if err != nil {
			return nil, err
		}

		tests = append(tests, test)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tests, nil
}

func (r *Repository) CreateTest(ctx context.Context, req dto.CreateTestRecord) (dto.TestResponse, error) {
	var test dto.TestResponse
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO test (username, password_hash, is_active)
		VALUES ($1, $2, $3)
		RETURNING id::text, username, is_active, created_at, updated_at
	`, req.Username, req.PasswordHash, req.IsActive).Scan(
		&test.ID,
		&test.Username,
		&test.IsActive,
		&test.CreatedAt,
		&test.UpdatedAt,
	)
	if err != nil {
		return dto.TestResponse{}, err
	}

	return test, nil
}

type testScanner interface {
	Scan(dest ...any) error
}

func scanTest(scanner testScanner) (dto.TestResponse, error) {
	var test dto.TestResponse

	err := scanner.Scan(
		&test.ID,
		&test.Username,
		&test.IsActive,
		&test.CreatedAt,
		&test.UpdatedAt,
	)
	if err != nil {
		return dto.TestResponse{}, err
	}

	return test, nil
}
