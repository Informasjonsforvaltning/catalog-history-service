# AGENTS.md

This file provides guidance to AI agents when working with code in this repository.

## Project Overview

Catalog History Service is a Go-based REST API that tracks changes (JSON Patch operations) to datasets, dataservices, concepts, and services. It uses Gin web framework with MongoDB for persistence and Keycloak for JWT authentication.

## Development Commands

```bash
# Install dependencies
go get

# Run locally (requires MongoDB via docker compose)
docker compose up -d
go run main.go

# Run tests (uses dockertest to spin up MongoDB container)
go test ./test

# Run tests with coverage
go test -v -race -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./test

# Build Docker image
docker build -t catalog-history-service .
```

## Architecture

**Layered architecture:**

- `handlers/` - HTTP request handlers (Gin handlers)
- `service/` - Business logic layer
- `repository/` - MongoDB data access with input validation
- `model/` - Data structures (Update, Person, JsonPatchOperation)
- `config/` - Environment variables, MongoDB connection, JWT security, router setup
- `logging/` - Logrus-based structured logging

**API Routes (all prefixed by catalogId/resourceId):**

- `POST /:catalogId/:resourceId/updates` - Create update (requires write permission)
- `GET /:catalogId/:resourceId/updates` - Get paginated updates (requires read permission)
- `GET /:catalogId/:resourceId/updates/:updateId` - Get specific update (requires read permission)
- `GET /ping` and `GET /ready` - Health/readiness probes

## Security

**Authorization model:**

- JWT tokens validated against Keycloak
- Token audience: `catalog-history-service`
- Admin pattern: `system:root:admin`
- Organization roles: `organization:{catalogId}:{admin|write|read}`

**Input validation (repository/validation.go):**

- NoSQL injection prevention - IDs reject `$`, `{}`, `[]`, quotes
- Sort field whitelist: `datetime`, `name`, `email`
- Pagination limits: max page 10,000, max size 100

## Testing

Tests are in `/test` directory and use:

- `testify` for assertions
- `dockertest` to spin up MongoDB containers
- Mock JWT store for authentication testing

Test database initialization: `test/init-mongo/init-mongo.js`

## Deployment

Kubernetes deployments via Kustomize in `deploy/`:

- `base/` - Common templates
- `staging/`, `prod/`, `demo/` - Environment-specific overlays

CI/CD runs via GitHub Actions - PRs deploy to staging, merges to main deploy to prod then demo.
