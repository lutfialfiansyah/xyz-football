# XYZ Football Management API

A robust RESTful API backend service built for XYZ Company to manage its amateur football teams. This application handles everything from team and player management to complex match scheduling, event tracking, and dynamic win/loss reporting.

## 🚀 Key Features

*   **Admin Authentication**: Secure JWT-based Login and Role-Based Access Control (RBAC).
*   **Team Management**: Comprehensive CRUD for football teams with soft-delete support preventing data loss.
*   **Player Management**: Detailed player tracking ensuring unique jersey numbers per team.
*   **Match Scheduling**: Prevents duplicate match schedules and self-matches (Team A vs Team A).
*   **Match Events Tracking**: Track Goals, Red/Yellow Cards, and Own Goals.
*   **Auto-Scoring System**: Match scores are dynamically updated when a `GOAL` or `OWN_GOAL` event is registered.
*   **Dynamic Reporting**: High-performance match summary generation outlining match status, aggregated total wins, and the top scorer of each match, powered by raw SQL queries to prevent DB layer congestion.

## 🛠️ Technology Stack

*   **Language**: [Go (Golang)](https://golang.org/)
*   **Framework**: [Gin Web Framework](https://gin-gonic.com/)
*   **Database**: [PostgreSQL](https://www.postgresql.org/)
*   **ORM**: [GORM](https://gorm.io/)
*   **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
*   **Security**: bcrypt, golang-jwt

## 🏗️ Architecture

This project strictly follows **Modular Clean Architecture**. 
Unlike standard layered architecture (which groups files by technical concern like all `repositories` in one folder), this approach groups code by **Feature Domains** (`admin`, `team`, `player`, `match`, `report`). 

This ensures:
*   High maintainability and separation of concerns.
*   Easier navigation when debugging specific business logic.
*   A clear path to transitioning into Microservices if XYZ Company expands in the future.

## 🐳 Quick Start with Docker (Recommended)

The fastest way to run the entire stack. **No Go or PostgreSQL installation required.**

### Prerequisites
*   Docker & Docker Compose

### Run
```bash
# Clone the repo
git clone <repository_url>
cd xyz-football-api

# Start everything (API + PostgreSQL)
make docker-up
```
The API will be available at `http://localhost:8080`.

### Stop & Data Management
```bash
# Stop containers but KEEP the database data
make docker-down

# Stop containers AND DELETE all data (fresh start)
make docker-destroy
```

---

## 💻 Local Development Setup


### Prerequisites
*   Go 1.20 or newer
*   PostgreSQL 13 or newer
*   Git

### Step 1: Clone the repository
```bash
git clone <repository_url>
cd xyz-football-api
```

### Step 2: Configure Environment Variables
Navigate to the config directory and create an `.env` file from the example:

```bash
cd internal/config/env
touch .env
```

Add the following configuration to the `.env` file (adjust the DB credentials according to your local Postgres setup):
```ini
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=xyz_football
JWT_SECRET=supersecretkey_change_in_production
JWT_EXPIRATION_SECONDS=86400
PORT=8080
UPLOAD_STORAGE=local
APP_URL=http://localhost:8080
```
Return to the project root directory: `cd ../../../`

### Step 3: Database Preparation
Create an empty database in your PostgreSQL instance:
```sql
CREATE DATABASE xyz_football;
```

### Step 4: Run Database Migrations
We provide a custom CLI tool to run Data Definition Language (DDL) queries. This sets up all the tables, strict constraints, and partial indexes correctly.

```bash
# To apply migrations (create tables)
make migrate-up

# To rollback migrations (drop tables) - USE CAUTION
# make migrate-down
```

### Step 5: Start the Server
```bash
make run
```
The server will start on port `8080` (or the port defined in your `.env`).

### Step 6: Initial Admin Login
On the first successful run, the application automatically seeds a default **Admin** account so you can begin testing immediately.

*   **Username**: `admin`
*   **Password**: `password`

Use the `/api/v1/login` endpoint with the above credentials to receive your JWT token. Add this token as a `Bearer` token in the `Authorization` header for all protected endpoints.

---

## 📡 API Endpoints Summary

### Public
*   `POST /api/v1/login` - Authenticate admin and receive JWT token

### Protected (Requires Bearer Token)
**Teams**
*   `POST /api/v1/teams/` - Register a new team
*   `GET /api/v1/teams/` - Get all active teams
*   `GET /api/v1/teams/:id` - Get specific team
*   `PUT /api/v1/teams/:id` - Update team info
*   `DELETE /api/v1/teams/:id` - Soft delete team

**Players**
*   `POST /api/v1/players/` - Add player to a team
*   `GET /api/v1/players/` - Get all active players
*   `GET /api/v1/players/:id` - Get specific player
*   `PUT /api/v1/players/:id` - Update player info
*   `DELETE /api/v1/players/:id` - Soft delete player

**Matches & Events**
*   `POST /api/v1/matches/` - Schedule a new match
*   `GET /api/v1/matches/` - List matches
*   `GET /api/v1/matches/:id` - Get match details
*   `PUT /api/v1/matches/:id/status` - Change match status (e.g., scheduled -> ongoing)
*   `DELETE /api/v1/matches/:id` - Delete a scheduled match
*    --
*   `POST /api/v1/matches/:id/events` - Submit a match event (GOAL, YELLOW_CARD, etc). *Note: `GOAL` and `OWN_GOAL` will automatically update the main scoreline.*
*   `GET /api/v1/matches/:id/events` - View all events for a match

**Reporting**
*   `GET /api/v1/reports/matches` - Generates a complex dynamic report detailing final match results, top scorers, and accumulated win statistics continuously for both Home and Away teams.
