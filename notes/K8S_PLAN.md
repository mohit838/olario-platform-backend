# Kubernetes Adoption Plan (Future)

Last Updated: 2026-05-26

## Decision

Kubernetes is a good fit for this project, but not for initial development.

## When to Adopt K8s

Move to Kubernetes after:

- Core services are stable locally
- Kong routing and service health checks are stable
- Basic migrations and deployment flow are repeatable

## Suggested Phases

1. Phase 0 (Now): local dependencies + minimal service bootstrapping
2. Phase 1: VPS deployment with docker compose or systemd containers
3. Phase 2: Kubernetes migration for scale and reliability

## Planned K8s Structure (Later)

```text
infra/k8s/
├── base/
│   ├── namespace.yaml
│   ├── configmaps/
│   ├── secrets-template/
│   ├── cms-service/
│   ├── catalog-service/
│   ├── heavy-service/
│   ├── kong/
│   └── frontend/
└── overlays/
    ├── dev/
    ├── staging/
    └── prod/
```

## Why K8s Later

    ├── staging/
    └── prod/

```

## Why K8s Later

- Better rollout/rollback
- Self-healing and autoscaling
- Cleaner production operations for many services
```
