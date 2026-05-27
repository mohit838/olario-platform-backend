CREATE TABLE report_requests (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  requested_by BIGINT,
  report_type VARCHAR(80) NOT NULL,
  status VARCHAR(30) NOT NULL DEFAULT 'queued',
  requested_for_month DATE,
  run_after TIMESTAMPTZ,
  output_object_key TEXT,
  error_message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  started_at TIMESTAMPTZ,
  completed_at TIMESTAMPTZ,
  FOREIGN KEY (tenant_id, requested_by) REFERENCES users(tenant_id, id) ON DELETE SET NULL (requested_by),
  CHECK (status IN ('queued', 'running', 'completed', 'failed', 'cancelled'))
);

CREATE INDEX idx_report_requests_tenant_id ON report_requests(tenant_id);
CREATE INDEX idx_report_requests_status ON report_requests(status);
CREATE INDEX idx_report_requests_run_after ON report_requests(run_after);

COMMENT ON TABLE report_requests IS 'Asynchronous report jobs, such as monthly revenue reports, scheduled outside peak hours.';
COMMENT ON COLUMN report_requests.run_after IS 'Earliest time the report worker should pick this request.';
COMMENT ON COLUMN report_requests.output_object_key IS 'Optional object storage key for generated report output.';
