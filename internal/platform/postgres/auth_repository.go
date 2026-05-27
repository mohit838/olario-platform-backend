package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	authapp "github.com/mohit838/olario-platform-backend/internal/application/auth"
)

// AuthRepository reads users needed for login.
// It keeps SQL out of the auth service while enforcing tenant scoping in the
// query itself.
type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindActiveUserForLogin(ctx context.Context, tenantSlug, email string) (authapp.User, error) {
	var user authapp.User
	err := r.db.QueryRowContext(ctx, `
		SELECT
			t.id,
			t.slug,
			u.id,
			COALESCE(r.id, 0),
			COALESCE(r.code, ''),
			u.name,
			u.email,
			u.password_hash
		FROM users u
		JOIN tenants t ON t.id = u.tenant_id
		LEFT JOIN roles r ON r.tenant_id = u.tenant_id AND r.id = u.role_id
		WHERE t.slug = $1
			AND lower(u.email) = lower($2)
			AND t.is_active = true
			AND u.is_active = true
			AND t.status IN ('trial', 'active')
		LIMIT 1
	`, tenantSlug, email).Scan(
		&user.TenantID,
		&user.TenantSlug,
		&user.UserID,
		&user.RoleID,
		&user.RoleCode,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		return authapp.User{}, fmt.Errorf("find active user for login: %w", err)
	}
	return user, nil
}

func (r *AuthRepository) RegisterTenantAdminWithInvitation(ctx context.Context, input authapp.RegisterInput, passwordHash, invitationTokenHash string) (authapp.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return authapp.User{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var invitationID int64
	if err := tx.QueryRowContext(ctx, `
		SELECT id
		FROM tenant_invitations
		WHERE email = $1
			AND token_hash = $2
			AND status = 'pending'
			AND expires_at > NOW()
		LIMIT 1
	`, input.Email, invitationTokenHash).Scan(&invitationID); err != nil {
		return authapp.User{}, fmt.Errorf("valid invitation not found: %w", err)
	}

	var planID int64
	if err := tx.QueryRowContext(ctx, `
		SELECT id
		FROM subscription_plans
		WHERE code = 'medium' AND is_active = true
		LIMIT 1
	`).Scan(&planID); err != nil {
		return authapp.User{}, fmt.Errorf("default medium plan not found: %w", err)
	}

	trialEndsAt := time.Now().UTC().AddDate(0, 3, 0)
	var tenantID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO tenants (name, slug, status, trial_ends_at)
		VALUES ($1, $2, 'trial', $3)
		RETURNING id
	`, input.TenantName, input.TenantSlug, trialEndsAt).Scan(&tenantID); err != nil {
		return authapp.User{}, fmt.Errorf("insert tenant: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO tenant_subscriptions (tenant_id, plan_id, status, starts_at, trial_ends_at)
		VALUES ($1, $2, 'trialing', NOW(), $3)
	`, tenantID, planID, trialEndsAt); err != nil {
		return authapp.User{}, fmt.Errorf("insert tenant subscription: %w", err)
	}

	var roleID int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO roles (tenant_id, name, code)
		VALUES ($1, 'Owner', 'owner')
		RETURNING id
	`, tenantID).Scan(&roleID); err != nil {
		return authapp.User{}, fmt.Errorf("insert owner role: %w", err)
	}

	var user authapp.User
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO users (tenant_id, role_id, name, email, password_hash, two_fa_enabled)
		VALUES ($1, $2, $3, $4, $5, true)
		RETURNING id, name, email, password_hash
	`, tenantID, roleID, input.AdminName, input.Email, passwordHash).Scan(&user.UserID, &user.Name, &user.Email, &user.PasswordHash); err != nil {
		return authapp.User{}, fmt.Errorf("insert tenant admin: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE tenant_invitations
		SET status = 'accepted', accepted_at = NOW()
		WHERE id = $1
	`, invitationID); err != nil {
		return authapp.User{}, fmt.Errorf("mark invitation accepted: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO audit_logs (tenant_id, actor_user_id, action, entity_type, entity_id, metadata)
		VALUES ($1, $2, 'tenant.registered', 'tenants', $1, '{}'::jsonb)
	`, tenantID, user.UserID); err != nil {
		return authapp.User{}, fmt.Errorf("insert registration audit log: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return authapp.User{}, fmt.Errorf("commit registration: %w", err)
	}

	user.TenantID = tenantID
	user.TenantSlug = input.TenantSlug
	user.RoleID = roleID
	user.RoleCode = "owner"
	return user, nil
}
