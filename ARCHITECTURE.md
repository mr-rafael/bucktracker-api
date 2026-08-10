# Architecture

BuckTracker is a Go REST API for savings and loan payment calculations. Users can calculate plans without persisting them, or authenticate and save plans for later retrieval, update, and deletion.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP | Standard library (`net/http`, `ServeMux`) |
| Database | PostgreSQL 15 |
| DB driver / pool | [`pgx/v5`](https://github.com/jackc/pgx) |
| SQL codegen | [sqlc](https://sqlc.dev/) |
| Migrations | [Goose](https://github.com/pressly/goose) |
| Auth | JWT ([`golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt)), bcrypt ([`golang.org/x/crypto`](https://pkg.go.dev/golang.org/x/crypto)) |
| Decimal math | [`shopspring/decimal`](https://github.com/shopspring/decimal) |
| Config | Environment variables via [`godotenv`](https://github.com/joho/godotenv) |
| Local runtime | Docker Compose |

Money amounts are stored and exchanged as **integer cents**. Interest rates and similar percentages are stored as **text** and computed with decimal arithmetic.

## Design Pattern: Handler → Service → Repository

The app uses a layered **handler-service-repository** architecture (also called controller-service-repository).

```
HTTP Request
    │
    ▼
┌─────────────┐     DTO decode / map      ┌─────────────┐
│   Handler   │ ───────────────────────▶  │   Service   │
│ internal/api│ ◀───────────────────────  │internal/svc │
└─────────────┘     domain / DTO encode   └──────┬──────┘
                                                 │
                                          domain objects
                                                 │
                                                 ▼
                                          ┌─────────────┐
                                          │ Repository  │
                                          │internal/repo│
                                          └──────┬──────┘
                                                 │
                                            sqlc Queries
                                                 │
                                                 ▼
                                            PostgreSQL
```

Responsibilities by layer:

| Layer | Package | Responsibility |
|---|---|---|
| **Handler** | `internal/api` | Parse HTTP, decode DTOs, call service, map results to response DTOs, write status codes |
| **Service** | `internal/service` | Business logic: validation, amortization / savings calculations, auth orchestration |
| **Repository** | `internal/repository` | Persist and load domain data via sqlc-generated queries |
| **Mapper** | `internal/mapper` | Convert between DTO ↔ domain (and occasionally DB row shapes) |
| **Domain** | `internal/domain` | Core business objects and calculation methods |
| **DTO** | `internal/dto` | Request/response JSON shapes |

Services depend on **repository interfaces** defined in the service package (e.g. `LoansRepository`, `SavingsRepository`). Concrete repos in `internal/repository` implement those interfaces. This keeps business logic testable with mocks (`internal/service/mocks.go`).

Composition root wiring lives in `internal/app.go`: env → DB pool → sqlc `Queries` → repos → services → handlers → mux.

## File Structure

```
bucktracker-api/
├── cmd/server/main.go          # Entry point: internal.New().Run()
├── internal/
│   ├── app.go                  # Dependency wiring + route registration
│   ├── api/                    # HTTP handlers + auth middleware
│   ├── auth/                   # JWT token generation helpers (library)
│   ├── db/
│   │   ├── migrations/         # Goose SQL migrations (source of truth for schema)
│   │   ├── query/              # Hand-written sqlc query files
│   │   ├── *.sql.go            # sqlc-generated query methods
│   │   ├── models.go           # sqlc-generated DB models
│   │   └── db.go               # sqlc Queries wrapper
│   ├── domain/                 # Business objects + calculation methods
│   ├── dto/                    # Request/response bodies
│   ├── mapper/                 # Cross-layer conversions
│   ├── repository/             # Data access adapters over sqlc
│   ├── service/                # Business logic + repo interfaces + mocks
│   └── utils/                  # Small shared helpers
├── sqlc.yaml                   # sqlc config (schema + queries → internal/db)
├── docker-compose.yml          # Postgres + goose migrate + API
├── Dockerfile
├── .env.example
├── go.mod
└── README.md
```

### Package notes

- **`internal/api`** — Handlers for users, auth, savings, loans, and health. `middleware.go` validates Bearer access tokens and injects `userID` into request context. Shared response helpers live in `api_handlers.go`.
- **`internal/auth`** — Pure JWT helpers (`GenerateAccessToken`, `GenerateRefreshToken`). Not an application layer; used by `AuthService`.
- **`internal/db`** — Database boundary:
  - `migrations/` — Goose up/down SQL
  - `query/` — Named sqlc queries (`-- name: CreateLoan :one`, etc.)
  - Generated Go is committed under `internal/db/`
- **`internal/domain`** — Domain models used by services and repositories. Calculation methods (e.g. loan month steps) live on domain types.
- **`internal/dto`** — JSON-facing structs with `json` tags.
- **`internal/mapper`** — Explicit conversion functions between DTO and domain (keeps handlers thin).
- **`internal/repository`** — Translates domain objects to/from sqlc params and rows.
- **`internal/service`** — Use cases. Defines repository interfaces it depends on.

## Data Flow

### Unauthenticated calculate (savings / loans)

Used for “what-if” calculations that are not stored.

1. Handler decodes `dto.*RequestParams` from JSON.
2. Mapper converts DTO → domain input (`LoansInput` / `SavingsInput`).
3. Service validates inputs, initializes a domain object, runs the calculation loop.
4. Mapper converts domain result → response DTO.
5. Handler writes JSON (`200 OK`).

No repository call is required for calculate-only endpoints.

### Authenticated save / get / update / delete

1. `AuthMiddleware` extracts `Authorization: Bearer <token>`, validates it via `AuthService`, and stores `userID` in context.
2. Handler reads `userID` from context (and path `{id}` when needed).
3. Mapper builds domain input including `UserID`.
4. Service validates, calculates (for save/update), then calls the repository.
5. Repository maps domain → sqlc params, executes queries, maps rows → domain (or returns sqlc/DB types where the service expects them).
6. Handler maps the result to a response DTO and writes JSON.

### Auth flow

1. **Register** (`POST /app/users/create`) — hash password with bcrypt, insert user.
2. **Login** (`POST /app/login`) — verify password, issue access JWT (15m) + refresh JWT (7d). Refresh token hash (SHA-256) is stored in `refresh_tokens`.
3. **Refresh** (`POST /app/refresh`) — validate refresh cookie/token, issue a new access token.
4. Protected routes use the access token middleware.

## Domain Model Overview

### Users & auth

- `users` — email (unique), password hash, username
- `refresh_tokens` — hashed refresh token, expiry, revoked flag

### Savings

- `savings` — plan parameters and calculated aggregates (`total_deposited`, `total_interest_earnings`, rates of return, etc.)
- `savings_state` — one row per month in the amortization-style projection (`interest`, `tax`, `contribution`, `increase`, `capital`)

Domain type: `domain.SavingsPlan` with methods that advance month-by-month.

### Loans

Loans are modeled as a loan plus one or more payment plans:

| Table | Role |
|---|---|
| `loans` | Loan terms (principal, rates, monthly/escrow payment, start date) and `default_payment_plan` FK |
| `payment_plans` | Plan aggregates: `duration_months`, `total_expenditure`, `total_paid`, `cost_of_credit` |
| `loan_state` | Monthly schedule rows for a payment plan |
| `principal_payments` | Extra principal payments associated with a payment plan |

Domain types:

- `domain.Loan` — loan fields + working calculation state + `DefaultPaymentPlan`
- `domain.LoanPaymentPlan` — aggregates, monthly `Plan` statuses, `PrincipalPayments`
- Calculation methods (`PassMonth`, `GenerateInterest`, `ChargeEscrow`, `MakePayment`, `FinalCalculations`) live on `Loan` and mutate the default payment plan

`loans.default_payment_plan` and `payment_plans.loan_id` form a circular FK. Create order is: loan → payment plan → set default on loan. Delete clears `default_payment_plan` before removing the loan so cascades can proceed.

## Database & Code Generation

Schema source of truth: `internal/db/migrations/*.sql` (Goose).

Query source of truth: `internal/db/query/*.sql` (sqlc).

Regenerate Go after schema/query changes:

```bash
sqlc generate
```

`sqlc.yaml` points schema at migrations, queries at `internal/db/query`, and emits package `db` into `internal/db` using `pgx/v5`.

Docker Compose runs Goose against the migrations directory before starting the API.

## Routing

Routes are registered in `internal/app.go` on the standard library mux.

| Area | Examples |
|---|---|
| Health | `GET /api/healthz` |
| Users / auth | `POST /app/users/create`, `POST /app/login`, `POST /app/refresh` |
| Savings | calculate, save, list, get, patch, delete under `/app/savings...` |
| Loans | calculate, save, list, get, patch, delete under `/app/loans...` |

Save/list/get/update/delete for savings and loans require auth middleware. Calculate endpoints are registered without it (public calculation).

## Testing Approach

- **Service tests** use mock repositories (`internal/service/mocks.go`) to exercise validation and calculation logic without a database.
- **Repository / handler tests** may use a real DB connection via test helpers in `internal/repository`.
- Run package tests with:

```bash
go test ./internal/...
```

## Local Development

1. Copy `.env.example` → `.env` for host-side runs, or use Compose-injected env.
2. `docker compose up --build -d` starts Postgres, applies migrations, and runs the API on `:8080`.
3. Required env vars: `POSTGRES_CONNECTION_STRING`, `ACCESS_SECRET`, `REFRESH_SECRET` (plus optional `ALLOWED_ORIGIN`, `ENV`).
