# Service Responsibilities

Last Updated: 2026-05-26

## cms-service (Java 21 LTS)
- CMS/admin operations
- business orchestration
- tenant/user/order coordination

## catalog-service (Go 1.26)
- products/categories/variants domain
- read-heavy catalog APIs

## heavy-service (C++ Drogon)
- pricing/discount/scoring/report calculations
- compute-heavy async workflows

## Boundary Rule
Each service owns its own data, logic, and migrations.
