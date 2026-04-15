# GitHub Copilot Instructions — gin-investment-tracker

## Stack

- **Language**: Go 1.25.7
- **Framework**: Gin (HTTP), pgx/v5 (PostgreSQL), golang-jwt/v5, bcrypt
- **Testing**: `testify` (`assert`, `require`, `mock`)
- **Docs**: swaggo/swag (Swagger)

---

## Project Layout

```
cmd/gin-investment-tracker/   ← entrypoint
internal/
  db/                         ← pgxpool connection
  dtos/                       ← request/response structs
  handlers/                   ← HTTP handlers (Gin)
  middlewares/                ← JWT middleware
  mocks/                      ← shared testify/mock structs
  models/                     ← DB models
  repositories/               ← repository interfaces + pgx implementations
  routes/                     ← route registration
  services/                   ← business logic + interfaces
  util/                       ← errors, jwt, bcrypt, validator, context
```

---

## Architecture Rules

- **Handlers** accept a **service interface** (not a concrete type) — enables mocking in tests.
- **Services** accept a **repository interface** — already defined in `internal/repositories/`.
- Interfaces for services live in `internal/services/<entity>_service_interface.go`.
- Error types use `util.AppError{Code, Message}` — `NewBadRequestError`, `NewNotFoundError`, `NewInternalError`.
- HTTP responses use `util.SendResponse` / `util.SendErrorResponse`.

---

## Validation Rules

**DTOs use `binding:` tags** (not `validate:`) so Gin's `ShouldBindJSON` handles all validation in one step — no separate `util.Validate.Struct()` call in the handler.

**Handler pattern** (mirror `user_handler.go`):

```go
var req dto.CreateFooRequest
if err := c.ShouldBindJSON(&req); err != nil {
    var ve validator.ValidationErrors
    if errors.As(err, &ve) {
        util.SendErrorResponse(c, http.StatusBadRequest, util.FormatValidationErrors(err))
        return
    }
    util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
    return
}
```

**Custom validators** are registered once with Gin's binding engine in `internal/util/validator.go`:

```go
func init() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("my_tag", func(fl validator.FieldLevel) bool { ... })
    }
}
```

Add a matching case to `getErrorMessage()` in `validator.go` for the human-readable error message.

Currently registered custom validators:

| Tag               | Valid values               |
| ----------------- | -------------------------- |
| `instrument_type` | `"stock"`, `"mutual_fund"` |

## Handler Helpers

Shared helpers live in `internal/handlers/helper.go`. **Always use these instead of inline parsing.**

| Helper                  | Signature                                          | Usage                                                                                |
| ----------------------- | -------------------------------------------------- | ------------------------------------------------------------------------------------ |
| `parseIntegerID`        | `(c *gin.Context, param string) (int64, error)`   | Parse a path parameter as int64                                                      |
| `parsePaginationParams` | `(c *gin.Context) (limit, offset int, err error)` | Parse `limit`/`offset` query params with defaults (50/0) and a silent max-cap of 200 |

**Usage pattern in a handler:**

```go
id, err := parseIntegerID(c, "entityId")
if err != nil {
    util.SendErrorResponse(c, http.StatusBadRequest, "invalid entity id")
    return
}

limit, offset, err := parsePaginationParams(c)
if err != nil {
    util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
    return
}
```

> `parsePaginationParams` caps `limit` at 200 silently — it does **not** error when limit > 200.

### File placement

| Layer   | Test file                                    |
| ------- | -------------------------------------------- |
| Service | `internal/services/<entity>_service_test.go` |
| Handler | `internal/handlers/<entity>_handler_test.go` |

### Package naming

- Service tests: `package service_test`
- Handler tests: `package handler_test`

### Test naming

```
Test<Type>_<Method>_<Scenario>

Examples:
  TestUserService_Login_EmailNotFound
  TestUserHandler_Create_PasswordTooShort
  TestUserHandler_Verify_NonBearerFormat
```

### Mock location

All `testify/mock` structs live in `internal/mocks/`:

- `mock_<entity>_repository.go` — implements the repository interface
- `mock_<entity>_service.go` — implements the service interface

### Handler test setup

```go
gin.SetMode(gin.TestMode)          // in init() or per-test
r := gin.New()                     // bare engine, no middleware noise
h := handler.NewUserHandler(svc)   // inject mock service
// register routes manually, matching routes/routes.go
```

### JWT in tests

```go
t.Setenv("JWT_SECRET", "test-secret-key")
token, _ := util.GenerateToken(userID, email)
req.Header.Set("Authorization", "Bearer "+token)
```

### Breaking scenarios to always cover

- Malformed JSON body → 400
- Each required field missing → 400 (one test per field)
- Field too short / invalid format → 400
- Service returns `BadRequestError` → 400
- Service returns `InternalError` → 500
- Protected endpoints: no auth header → 401, non-Bearer → 401, invalid token → 401
- Service returns `NotFoundError` → 404

---

## CI Gate (GitHub Actions)

Workflow: `.github/workflows/test.yml`  
Job name: **`unit-tests`** ← use this exact name in branch protection rules.

Runs on every `push` to `main` and every `pull_request` targeting `main`.  
Command: `go test -count=1 -race ./...`

### Enable branch protection on GitHub

1. Settings → Branches → Add rule for `main`
2. ✅ Require status checks to pass → add **`unit-tests`**
3. ✅ Require branches to be up to date before merging
4. ✅ Do not allow bypassing the above settings

---

## Local Git Hook (pre-push)

Blocks `git push origin main` if any test fails.

**First-time setup (run once after cloning):**

```bash
make install-hooks
```

Hook source: `scripts/hooks/pre-push`  
Installed to: `.git/hooks/pre-push`

---

## Useful Make Targets

| Target               | Purpose                   |
| -------------------- | ------------------------- |
| `make swagger`       | Regenerate Swagger docs   |
| `make test`          | Run all unit tests        |
| `make install-hooks` | Install pre-push git hook |

---

## Adding a New Module (e.g., `investment`)

1. Define model in `internal/models/investment.go`
2. Define DTOs in `internal/dtos/investment_dto.go`
3. Define repository interface in `internal/repositories/investment_repository.go`
4. Implement repository in `internal/repositories_impl/investment_repository.go`
5. Define service interface in `internal/services/investment_service_interface.go`
6. Implement service in `internal/services/investment_service.go`
7. Implement handler in `internal/handlers/investment_handler.go` — accept `InvestmentServiceInterface`
8. Register routes in `internal/routes/routes.go`
9. Add mock files `internal/mocks/mock_investment_repository.go` + `mock_investment_service.go`
10. Write `internal/services/investment_service_test.go` and `internal/handlers/investment_handler_test.go`
