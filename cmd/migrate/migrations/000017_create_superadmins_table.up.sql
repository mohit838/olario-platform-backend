CREATE TABLE superadmins (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(120) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  two_fa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_superadmins_is_active ON superadmins(is_active);

COMMENT ON TABLE superadmins IS 'Platform-level admins who manage tenant invitations, billing review, and tenant deactivation requests.';
COMMENT ON COLUMN superadmins.email IS 'Unique login email for a platform superadmin.';
COMMENT ON COLUMN superadmins.password_hash IS 'Stores a password hash only, never a plain text password.';
COMMENT ON COLUMN superadmins.two_fa_enabled IS '2FA flag for superadmin accounts; seeded local accounts keep it disabled for easier testing.';
