# 🤖 Agent Guide (AGENTS.md)

Welcome! This repository, **XYZ Football API**, follows a specific set of architectural patterns and coding conventions. To maintain consistency and high-quality code, please follow these guidelines when making modifications.

---

## 🏗️ Architecture Overview

The project is structured using **Modular Clean Architecture**. Each business domain is encapsulated within a module in the `internal/modules/` directory.

### Module Structure
Each module (e.g., `player`, `team`, `match`) consists of:
-   `entity.go`: Core domain models, GORM tags, and Request/Response DTOs.
-   `repository/postgres/`: Data access layer for PostgreSQL (using GORM).
-   `usecase/`: Business logic layer, often split into individual files (e.g., `create.go`, `get.go`).
-   `delivery/http/`: Gin-based HTTP handlers and routing.

---

## 📂 Key Directory Map

-   `cmd/`: Entry points for the `api` and `migrate` (CLI) tools.
-   `internal/config/`: Configuration management (DB connections, Env loading).
-   `internal/middleware/`: Global and route-specific Gin middleware (e.g., Auth).
-   `internal/pkg/`: Shared utility packages (Errors, Response formatting, Pagination).
-   `migrations/`: SQL migration files (`.up.sql` and `.down.sql`).
-   `uploads/`: Local storage for uploaded files.

---

## 🛠️ Tech Stack & Patterns

-   **Language:** Go (Golang)
-   **Web Framework:** [Gin Gonic](https://github.com/gin-gonic/gin)
-   **ORM:** [GORM](https://gorm.io/) with PostgreSQL
-   **Authentication:** JWT (JSON Web Tokens)
-   **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)
-   **ID Strategy:** UUID (v4) for all primary keys.

---

## 📏 Coding Standards

### 1. Error Handling
Always use the custom `apperror` package for business logic and validation errors. It supports bilingual messages (English and Indonesian).
```go
// Example
return apperror.New(http.StatusBadRequest, "Invalid team ID", "ID tim tidak valid")
```

### 2. Validation
Use Gin's `binding` tags in request structs.
```go
type CreatePlayerRequest struct {
    Name     string `json:"name" binding:"required"`
    Position string `json:"position" binding:"required,oneof=forward midfielder defender goalkeeper"`
}
```

### 3. Response Formatting
Standardize responses using the `response` package.
```go
response.Success(c, http.StatusOK, "Player created successfully", data)
```

### 4. Naming Conventions
-   **Files:** `snake_case.go` (e.g., `player_handler.go`)
-   **Variables/Functions:** `CamelCase` (public) or `lowerCamelCase` (private).
-   **Interfaces:** Usually defined in the `usecase` or `repository` folder to enable dependency injection.

---

## 🚀 Development Workflow

### Adding a New Module
1.  Define the `Entity` and `Request/Response` structs.
2.  Create the `Repository` interface and its Postgres implementation.
3.  Implement the `Usecase` logic.
4.  Write the `Handler` and register routes in `main.go`.
5.  Add any necessary SQL migrations.

### Migrations
Always create a migration for schema changes:
```bash
# Example
# migrations/000002_add_email_to_admins.up.sql
# migrations/000002_add_email_to_admins.down.sql
```

---

## 🤖 Instructions for AI Agents
-   **Surgical Edits:** When modifying a file, preserve existing indentation and comments.
-   **Modular Scope:** Stay within the relevant module directory unless modifying `internal/pkg` or `cmd/`.
-   **Self-Correction:** Always run `go build ./...` or `go test ./...` after changes to ensure type safety.
-   **Documentation:** If adding a new public function, include a brief Go doc comment.

---
