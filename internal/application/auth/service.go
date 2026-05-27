package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/mohit838/olario-platform-backend/internal/security/token"
)

type User struct {
	TenantID     int64
	TenantSlug   string
	UserID       int64
	RoleID       int64
	RoleCode     string
	Name         string
	Email        string
	PasswordHash string
}

type UserRepository interface {
	FindActiveUserForLogin(ctx context.Context, tenantSlug, email string) (User, error)
	RegisterTenantAdminWithInvitation(ctx context.Context, input RegisterInput, passwordHash, invitationTokenHash string) (User, error)
}

type RefreshSession struct {
	TenantID int64  `json:"tenant_id"`
	UserID   int64  `json:"user_id"`
	RoleID   int64  `json:"role_id"`
	RoleCode string `json:"role_code"`
	Email    string `json:"email"`
}

type RefreshStore interface {
	StoreRefresh(ctx context.Context, tokenHash string, session RefreshSession, ttl time.Duration) error
	RotateRefresh(ctx context.Context, oldHash, newHash string, ttl time.Duration) (RefreshSession, error)
	DeleteRefresh(ctx context.Context, tokenHash string) error
}

type Service struct {
	users      UserRepository
	refreshes  RefreshStore
	tokens     *token.Manager
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type LoginInput struct {
	TenantSlug string
	Email      string
	Password   string
}

type RegisterInput struct {
	InvitationToken string
	TenantName      string
	TenantSlug      string
	AdminName       string
	Email           string
	Password        string
}

type TokenPair struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	TokenType             string    `json:"token_type"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func NewService(users UserRepository, refreshes RefreshStore, tokens *token.Manager, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		users:      users,
		refreshes:  refreshes,
		tokens:     tokens,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *Service) Login(ctx context.Context, input LoginInput) (TokenPair, error) {
	user, err := s.users.FindActiveUserForLogin(ctx, input.TenantSlug, input.Email)
	if err != nil {
		return TokenPair{}, fmt.Errorf("invalid login credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return TokenPair{}, fmt.Errorf("invalid login credentials")
	}

	return s.issuePair(ctx, RefreshSession{
		TenantID: user.TenantID,
		UserID:   user.UserID,
		RoleID:   user.RoleID,
		RoleCode: user.RoleCode,
		Email:    user.Email,
	})
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (TokenPair, error) {
	if input.InvitationToken == "" || input.TenantName == "" || input.TenantSlug == "" || input.AdminName == "" || input.Email == "" || input.Password == "" {
		return TokenPair{}, fmt.Errorf("invitation token, tenant name, tenant slug, admin name, email, and password are required")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return TokenPair{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.users.RegisterTenantAdminWithInvitation(ctx, input, string(passwordHash), token.HashInvitationToken(input.InvitationToken))
	if err != nil {
		return TokenPair{}, fmt.Errorf("register tenant admin: %w", err)
	}

	return s.issuePair(ctx, RefreshSession{
		TenantID: user.TenantID,
		UserID:   user.UserID,
		RoleID:   user.RoleID,
		RoleCode: user.RoleCode,
		Email:    user.Email,
	})
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (TokenPair, error) {
	oldHash := s.tokens.HashRefreshToken(refreshToken)
	newRaw, newHash, err := s.tokens.NewRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}

	session, err := s.refreshes.RotateRefresh(ctx, oldHash, newHash, s.refreshTTL)
	if err != nil {
		return TokenPair{}, fmt.Errorf("invalid refresh token")
	}

	accessRaw, claims, err := s.tokens.NewAccessToken(token.Claims{
		TenantID: session.TenantID,
		UserID:   session.UserID,
		RoleID:   session.RoleID,
		RoleCode: session.RoleCode,
	}, s.accessTTL)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:           accessRaw,
		RefreshToken:          newRaw,
		TokenType:             "Bearer",
		AccessTokenExpiresAt:  time.Unix(claims.ExpiresAt, 0).UTC(),
		RefreshTokenExpiresAt: time.Now().UTC().Add(s.refreshTTL),
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	return s.refreshes.DeleteRefresh(ctx, s.tokens.HashRefreshToken(refreshToken))
}

func (s *Service) issuePair(ctx context.Context, session RefreshSession) (TokenPair, error) {
	refreshRaw, refreshHash, err := s.tokens.NewRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}
	return s.issuePairWithRefresh(ctx, session, refreshRaw, refreshHash)
}

func (s *Service) issuePairWithRefresh(ctx context.Context, session RefreshSession, refreshRaw, refreshHash string) (TokenPair, error) {
	accessRaw, claims, err := s.tokens.NewAccessToken(token.Claims{
		TenantID: session.TenantID,
		UserID:   session.UserID,
		RoleID:   session.RoleID,
		RoleCode: session.RoleCode,
	}, s.accessTTL)
	if err != nil {
		return TokenPair{}, err
	}

	if err := s.refreshes.StoreRefresh(ctx, refreshHash, session, s.refreshTTL); err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:           accessRaw,
		RefreshToken:          refreshRaw,
		TokenType:             "Bearer",
		AccessTokenExpiresAt:  time.Unix(claims.ExpiresAt, 0).UTC(),
		RefreshTokenExpiresAt: time.Now().UTC().Add(s.refreshTTL),
	}, nil
}
