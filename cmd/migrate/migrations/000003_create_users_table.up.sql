CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  role_id BIGINT,
  name VARCHAR(120) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password_hash TEXT NOT NULL,
  two_fa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (tenant_id, id),
  UNIQUE (tenant_id, email),
  FOREIGN KEY (tenant_id, role_id) REFERENCES roles(tenant_id, id) ON DELETE SET NULL (role_id)
);

CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_is_active ON users(is_active);

COMMENT ON TABLE users IS 'Tenant-scoped platform users who can operate a grocery shop.';
COMMENT ON COLUMN users.tenant_id IS 'Scopes the user to one tenant.';
COMMENT ON COLUMN users.role_id IS 'Optional tenant-safe role assignment for access control.';
COMMENT ON COLUMN users.email IS 'Tenant-unique login/contact email.';
COMMENT ON COLUMN users.password_hash IS 'Stores a password hash only, never a plain text password.';
COMMENT ON COLUMN users.two_fa_enabled IS 'Shows whether this employee has enabled two-factor authentication.';
COMMENT ON COLUMN users.is_active IS 'Disables login or operation access without deleting the user.';
