package test

import (
	"context"
	"errors"
	"strings"

	"github.com/mohit838/olario-platform-backend/cmd/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidTestRequest = errors.New("invalid test request")

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
	username := strings.TrimSpace(req.Username)
	if username == "" || req.Password == "" {
		return dto.TestResponse{}, ErrInvalidTestRequest
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.TestResponse{}, err
	}

	return s.repository.CreateTest(ctx, dto.CreateTestRecord{
		Username:     username,
		PasswordHash: string(passwordHash),
		IsActive:     isActive,
	})
}
