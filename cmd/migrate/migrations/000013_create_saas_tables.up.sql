CREATE TABLE subscription_plans (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(50) NOT NULL UNIQUE,
  name VARCHAR(100) NOT NULL,
  billing_period VARCHAR(20) NOT NULL,
  price NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (price >= 0),
  max_employees INT,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (billing_period IN ('monthly', 'yearly')),
  CHECK (max_employees IS NULL OR max_employees > 0)
);

CREATE TABLE tenant_subscriptions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  plan_id BIGINT NOT NULL REFERENCES subscription_plans(id),
  status VARCHAR(30) NOT NULL DEFAULT 'trialing',
  starts_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  trial_ends_at TIMESTAMPTZ,
  current_period_ends_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (status IN ('trialing', 'active', 'past_due', 'cancelled', 'expired'))
);

CREATE TABLE tenant_invitations (
  id BIGSERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL,
  invited_by_superadmin_id BIGINT,
  status VARCHAR(30) NOT NULL DEFAULT 'pending',
  token_hash TEXT NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  accepted_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (email, status),
  CHECK (status IN ('pending', 'accepted', 'expired', 'cancelled'))
);

CREATE TABLE tenant_deactivation_requests (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  requested_by BIGINT,
  reason TEXT,
  status VARCHAR(30) NOT NULL DEFAULT 'pending',
  reviewed_by_superadmin_id BIGINT,
  reviewed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  FOREIGN KEY (tenant_id, requested_by) REFERENCES users(tenant_id, id) ON DELETE SET NULL (requested_by),
  CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled'))
);

CREATE INDEX idx_tenant_subscriptions_tenant_id ON tenant_subscriptions(tenant_id);
CREATE INDEX idx_tenant_invitations_email ON tenant_invitations(email);
CREATE INDEX idx_tenant_deactivation_requests_tenant_id ON tenant_deactivation_requests(tenant_id);
CREATE INDEX idx_tenant_deactivation_requests_status ON tenant_deactivation_requests(status);

COMMENT ON TABLE subscription_plans IS 'Available SaaS plans such as medium or business.';
COMMENT ON COLUMN subscription_plans.max_employees IS 'Employee limit for a plan; NULL means unlimited.';
COMMENT ON TABLE tenant_subscriptions IS 'Tenant subscription lifecycle, including trial and paid periods.';
COMMENT ON TABLE tenant_invitations IS 'Invite-only tenant onboarding controlled by superadmin email invitations.';
COMMENT ON COLUMN tenant_invitations.token_hash IS 'Stores invitation token hash only, never the raw token.';
COMMENT ON TABLE tenant_deactivation_requests IS 'Tenant admin request to deactivate the tenant account, reviewed by superadmin.';
