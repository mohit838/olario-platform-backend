CREATE TABLE tenants (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(120) NOT NULL,
  slug VARCHAR(120) NOT NULL UNIQUE,
  status VARCHAR(30) NOT NULL DEFAULT 'trial',
  trial_ends_at TIMESTAMPTZ,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (status IN ('invited', 'trial', 'active', 'past_due', 'suspended', 'deactivation_requested', 'deactivated'))
);

CREATE INDEX idx_tenants_is_active ON tenants(is_active);

COMMENT ON TABLE tenants IS 'Stores each independent shop or organization using the platform.';
COMMENT ON COLUMN tenants.id IS 'Internal tenant primary key used to scope all tenant-owned data.';
COMMENT ON COLUMN tenants.name IS 'Display name for the shop or organization.';
COMMENT ON COLUMN tenants.slug IS 'Public-safe unique tenant identifier used in URLs, headers, or admin tooling.';
COMMENT ON COLUMN tenants.status IS 'Tenant lifecycle status for invite-only onboarding, trial, paid, suspended, and deactivation flows.';
COMMENT ON COLUMN tenants.trial_ends_at IS 'End date for the initial free trial period.';
COMMENT ON COLUMN tenants.is_active IS 'Disables a tenant without deleting its historical data.';
