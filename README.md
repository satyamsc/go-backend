# Go Backend Challenge

REST API for managing devices, built with Go, Gin, and GORM. Configuration via Viper. SQLite storage by default with Dockerized runtime and Swagger documentation.

## Features

- CRUD for `devices` with validation and filtering by `brand` and `state`
- Immutable `created_at` and restricted updates while `state` is `in-use`
- Consistent timestamp format `DD.MM.YYYY HH:mm:ss` in responses
- Centralized JSON error payloads with codes
- Swagger UI at `/docs` and spec at `/openapi.yaml`
- CORS and panic recovery middleware enabled

## Stack

- `Go` 1.25.x
- `Gin` for HTTP routing
- `GORM` ORM
- `Viper` for configuration (env + optional YAML)
- `SQLite` (pure-Go driver) by default

## Directory Structure

- `cmd/app/` entrypoint
- `config/` configuration (`config.go`, optional `config.yaml`)
- `internal/` models, repositories, services, handlers, routers, middlewares
- `database/` database connection
- `pkg/` shared utilities (error, logger, etc.)
- `docs/swagger` OpenAPI spec
- `test/` unit and integration tests

## Configuration

- `DB_PATH` path to SQLite database file (defaults to `./data/devices.db`)
- `SERVER_ADDR` server listen address (defaults to `:8080`)
- Optional file `config/config.yaml` can set the same keys; env vars override file values

## Run Locally

- `make build`
- `mkdir -p data`
- `DB_PATH=./data/devices.db SERVER_ADDR=:8080 ./app`
- Swagger UI: `http://localhost:8080/docs`

## Docker

- Build: `make docker-build`
- Run: `make docker-run`
- The container mounts `./data` to persist the SQLite file and exposes `8080`

## Docker Compose

- `docker compose up -d`
- Services:
  - `app`: the API server (uses SQLite by default via `DB_PATH`)
  - `postgres`: available for future use; the current app code uses SQLite
- To switch the app to Postgres, update the database driver and configuration in code to connect via DSN, then set appropriate env vars in compose.

## API Overview

- Base URL: `http://localhost:8080`
- Endpoints:
  - `GET /healthz` returns `200`
  - `GET /docs` Swagger UI
  - `GET /openapi.yaml` OpenAPI spec
  - `POST /devices`
  - `GET /devices?brand=...&state=...`
  - `GET /devices/:id`
  - `PUT /devices/:id` (cannot change `created_at`; restricted while `in-use`)
  - `PATCH /devices/:id` (cannot change `created_at`; name/brand blocked while `in-use`)
  - `DELETE /devices/:id` (blocked while `in-use`)

### Schemas

- `state` one of `available`, `in-use`, `inactive`
- Response `created_at` formatted as `DD.MM.YYYY HH:mm:ss`

### Error Payload

```json
{
  "code": "invalid_state",
  "message": "...",
  "details": { "field": "state" },
  "timestamp": "2025-12-14T20:01:13Z"
}
```

## Examples

- Create:

```sh
curl -sS -X POST http://localhost:8080/devices \
  -H 'Content-Type: application/json' \
  -d '{"name":"Phone X","brand":"Acme","state":"available"}'
```

- List:

```sh
curl -sS 'http://localhost:8080/devices?brand=Acme&state=available'
```

## Testing

- Run all tests: `make test`
- Integration tests hit HTTP handlers and routes in `test/integration`
- Unit tests cover services and models in `test/unit`


