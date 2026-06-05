package test

import (
	"context"

	"github.com/mohit838/olario-platform-backend/cmd/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) ListTests(ctx context.Context) ([]dto.TestResponse, error) {
	return s.repository.ListTests(ctx)
}

func (s *Service) CreateTest(ctx context.Context, req dto.TestRequest) (dto.TestResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.TestResponse{}, err
	}

	return s.repository.CreateTest(ctx, dto.CreateTestRecord{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		IsActive:     req.IsActive,
	})
}
