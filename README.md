# go-todo-app

Translations: 
- [日本語](README_ja.md)

## Overview

Learning Go by building a backend-only to-do API. Everything required to run the API and a Postgres database locally is containerized, so you only need Docker installed on your machine.

## Requirements

- Docker (24.x or newer recommended)
- Docker Compose plugin
- Two free ports on your host: `8080` (API) and `5432` (Postgres)

## Quick start

```bash
docker compose up --build
```

- The API is served from <http://localhost:8080>
- Postgres is available on `localhost:5432` with `todo` / `todo` credentials
- On first boot Postgres runs the SQL files in `db/init/` to create the `todos` table

Stop the stack with `CTRL+C`. Use `docker compose down -v` when you also want to wipe local database files.

### Example requests

```bash
# Health check
curl http://localhost:8080/healthz

# Create a todo
curl -X POST http://localhost:8080/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"Learn Go with Docker"}'

# List todos
curl http://localhost:8080/todos
```

## Project layout

- `cmd/server/`: program entrypoint (wires config, DB, HTTP server)
- `internal/config`: environment variable parsing
- `internal/database`: DB connection helper
- `internal/httpserver`: HTTP server + middleware
- `internal/todo`: domain model, repository, and HTTP handler
- `db/init/`: SQL files that Postgres executes on bootstrap

## Developing inside Docker

The `api` service in `docker-compose.yml` uses the `dev` stage of the `Dockerfile`. Source files are bind-mounted, and [`air`](https://github.com/air-verse/air) automatically rebuilds/restarts the server on file changes.

Useful commands:

```bash
# Follow API logs
docker compose logs -f api

# Run Go tests inside the container
docker compose exec api go test ./...

# Rebuild the API container after changing Go dependencies
docker compose build api
```

> 💡 The first build will fetch Go modules inside the container and create `go.sum` in your local workspace automatically. If you add dependencies later, run `docker compose exec api go mod tidy` to keep modules tidy.

## Configuration

Environment variables (all optional) with their defaults:

| Variable        | Default                                                     | Description                      |
| --------------- | ----------------------------------------------------------- | -------------------------------- |
| `PORT`          | `8080`                                                      | HTTP port exposed by the API     |
| `DATABASE_URL`  | `postgres://todo:todo@localhost:5432/todo?sslmode=disable`  | Connection string for Postgres   |
| `POSTGRES_USER` | `todo` (docker compose only)                                | DB username                      |
| `POSTGRES_DB`   | `todo` (docker compose only)                                | Database name                    |

Override them in `docker-compose.yml`, via `docker compose --env-file`, or in your shell when running the Go binary directly.
