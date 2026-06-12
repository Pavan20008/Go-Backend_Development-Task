# User Age API

A small, production-style RESTful API written in **Go** that manages users
(`name` + `dob`) and **calculates each user's age dynamically** from their date
of birth using Go's `time` package.

## Tech Stack

| Concern        | Choice                                   |
| -------------- | ---------------------------------------- |
| HTTP framework | [GoFiber](https://gofiber.io/) (v2)      |
| Database       | PostgreSQL                               |
| DB access      | [SQLC](https://sqlc.dev/) (`pgx/v5`)     |
| Logging        | [Uber Zap](https://github.com/uber-go/zap) |
| Validation     | [go-playground/validator](https://github.com/go-playground/validator) |

## Features

- CRUD endpoints for users
- **Age computed on the fly** from `dob` (never stored)
- Input validation (`name` required, `dob` must be a valid past `YYYY-MM-DD` date)
- Structured request logging with **request duration** and a **request id**
- `X-Request-ID` middleware (reuses an incoming id or generates one)
- Pagination on `GET /users` (`?page=` & `?limit=`)
- Docker + docker-compose for one-command startup
- Unit tests for the age calculation
- Graceful shutdown and DB connection retry

## Project Structure

```
.
├── cmd/server/main.go        # entrypoint: config, DI wiring, server lifecycle
├── config/                   # env-based configuration
├── db/
│   ├── migrations/           # SQL migrations (*.up.sql / *.down.sql)
│   ├── query/                # SQLC query definitions
│   ├── sqlc/                 # generated DB access layer (DO NOT EDIT)
│   └── embed.go              # embeds + applies migrations on startup
└── internal/
    ├── handler/              # HTTP handlers, validation, error rendering
    ├── repository/           # data access (wraps sqlc)
    ├── service/              # business logic + age calculation
    ├── routes/               # route registration
    ├── middleware/           # request id + request logging
    ├── models/               # request/response DTOs
    └── logger/               # Zap logger setup
```

## Database Schema

```sql
CREATE TABLE users (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    dob  DATE NOT NULL
);
```

## Getting Started

### Option A — Docker (recommended)

Requires Docker. Starts both Postgres and the API:

```bash
docker compose up --build
# API is available on http://localhost:8080
```

### Option B — Run locally

Requires Go 1.23+ and a running PostgreSQL.

```bash
# 1. Start a Postgres (example with Docker)
docker run --name pg -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=user_age_api -p 5432:5432 -d postgres:16-alpine

# 2. Configure env (defaults already match the command above)
cp .env.example .env

# 3. Run — migrations are applied automatically on startup
make run            # or: go run ./cmd/server
```

Configuration is read from environment variables (see `.env.example`). You can
provide a full `DATABASE_URL`, or the individual `DB_*` parts.

## Regenerating the DB layer (SQLC)

After editing `db/query/*.sql` or `db/migrations/*.sql`:

```bash
make sqlc   # or: sqlc generate
```

## Migrations

Migrations live in `db/migrations` and are **applied automatically on startup**
(they are idempotent). For an explicit workflow you can also use the
[golang-migrate](https://github.com/golang-migrate/migrate) CLI:

```bash
make migrate-up
make migrate-down
```

## API

Base URL: `http://localhost:8080`

| Method | Path         | Description            | Success status |
| ------ | ------------ | ---------------------- | -------------- |
| POST   | `/users`     | Create a user          | `201 Created`  |
| GET    | `/users/:id` | Get a user (with age)  | `200 OK`       |
| GET    | `/users`     | List users (with age)  | `200 OK`       |
| PUT    | `/users/:id` | Update a user          | `200 OK`       |
| DELETE | `/users/:id` | Delete a user          | `204 No Content` |
| GET    | `/health`    | Health check           | `200 OK`       |

### Create User

```bash
curl -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","dob":"1990-05-10"}'
```

```json
{ "id": 1, "name": "Alice", "dob": "1990-05-10" }
```

### Get User by ID

```bash
curl http://localhost:8080/users/1
```

```json
{ "id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35 }
```

### List Users (with pagination)

```bash
curl "http://localhost:8080/users?page=1&limit=10"
```

```json
[
  { "id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35 }
]
```

Pagination metadata is returned via response headers: `X-Total-Count`,
`X-Page`, `X-Limit`.

### Update User

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice Updated","dob":"1991-03-15"}'
```

```json
{ "id": 1, "name": "Alice Updated", "dob": "1991-03-15" }
```

### Delete User

```bash
curl -i -X DELETE http://localhost:8080/users/1
# HTTP/1.1 204 No Content
```

### Error format

Validation and other errors return a consistent JSON envelope:

```json
{ "error": "validation failed", "fields": { "Dob": "must be a valid date (YYYY-MM-DD) and not in the future" } }
```

## Testing

```bash
make test   # or: go test ./...
```

The age calculation is covered by table-driven unit tests in
`internal/service/age_test.go` (including leap-year and boundary cases).

## Make targets

Run `make help` to see all available targets (`run`, `build`, `test`, `vet`,
`sqlc`, `migrate-up`, `docker-up`, ...).
