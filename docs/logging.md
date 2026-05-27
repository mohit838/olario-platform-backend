# Logging And Audit Logs

Runtime logs currently use Go `log/slog` and write to stdout.

File rotation, daily compression, and MinIO upload are intentionally delayed
until the observability step. For now, container platforms and local terminals
can collect stdout logs.

## Example Runtime Log

```text
time=2026-05-28T10:00:00+06:00 level=INFO msg="http server starting" addr=:8080 env=local
time=2026-05-28T10:00:05+06:00 level=INFO msg="http request completed" method=POST path=/api/v1/dev/full-circle status=201 bytes=512 duration=20ms remote_addr=127.0.0.1:50000
```

## Audit Logs

Business/security audit logs are stored in Postgres first:

```text
audit_logs
```

Use audit logs for important events:

- tenant invited
- employee deactivated
- order created
- report requested
- subscription changed
- permission changed

Do not use audit logs for noisy request logging.
