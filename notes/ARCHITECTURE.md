# Architecture Overview

Last Updated: 2026-05-26

## High-Level Flow
Angular Frontend -> Kong Gateway -> CMS/Catalog/Heavy Services

## Service Communication
- External: REST via Kong
- Internal: gRPC
- Async events: Kafka

## Dependency Source
- PostgreSQL, Kafka, and MinIO are provided from local machine services.
- Project services connect using local endpoints from `.env`.

## Core Rules
- Database per service
- No direct cross-service DB queries
- Services are independently deployable
