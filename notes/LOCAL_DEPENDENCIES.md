# Local Machine Dependencies

Last Updated: 2026-05-26

These dependencies are provided from your local machine Docker containers, not from this project compose stack.

## Current Host Ports (from your environment)
- Redis: `localhost:1236` (`local-redis-stack`)
- MinIO S3 API: `localhost:1239` (`local-minio`)
- MinIO Console: `localhost:1240`
- PostgreSQL: `local-postgres` container exists but host port mapping is not currently exposed
- Kafka-compatible broker: Redpanda containers exist but currently stopped

## Required Actions Before App Start
1. Start `local-postgres` and ensure host port is available (typically `5432`).
2. Start Redpanda (or another Kafka-compatible broker) and expose a bootstrap port (typically `9092`).
3. Keep MinIO and Redis running on the configured ports.

## Recommended `.env` Mapping
- `DB_HOST=localhost`
- `DB_PORT=5432` (or your mapped port)
- `REDIS_HOST=localhost`
- `REDIS_PORT=1236`
- `MINIO_HOST=localhost`
- `MINIO_API_PORT=1239`
- `MINIO_CONSOLE_PORT=1240`
- `S3_ENDPOINT=http://localhost:1239`
- `KAFKA_BOOTSTRAP_SERVERS=localhost:9092` (or your Redpanda mapped port)
