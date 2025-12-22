# Devcontainer Setup Guide

This document explains the devcontainer setup for the Budget Planner project and the fixes that have been applied.

## Overview

The devcontainer provides a complete development environment with:
- Go 1.23 backend with Chi framework
- Node.js LTS frontend with SvelteKit
- PostgreSQL 16 database
- VS Code with all necessary extensions
- Automated dependency installation

## Services

### App Service
- Base image: `mcr.microsoft.com/devcontainers/base:ubuntu-22.04`
- Features: Go, Node.js, Docker-in-Docker
- Exposed ports: 5173 (frontend), 8080 (backend)
- Dependencies: PostgreSQL

### PostgreSQL Service
- Image: `postgres:16-alpine`
- Database: `budgetdb`
- User: `budgetuser`
- Password: `budgetpass`
- Port: 5432
- Automatic schema initialization via `init-db.sql`

## Configuration Files

### devcontainer.json
- Defines VS Code extensions and settings
- Sets up port forwarding
- Configures post-create command
- Includes SQLTools for database management

### docker-compose.yml
- Orchestrates container services
- Manages volumes for persistence
- Sets up network isolation
- Removed unnecessary Redis and Mailhog services

### post-create.sh
- Installs Go dependencies (sqlc, migrate, air)
- Installs Node.js dependencies
- Creates environment files
- Waits for PostgreSQL to be ready
- Generates development certificates

### Dockerfile.backend
- Uses Go 1.23-alpine (fixed from 1.22)
- Installs PostgreSQL client and development tools
- Sets up air for live reloading
- Runs as non-root user

### Dockerfile.frontend
- Uses Node.js 20-alpine
- Installs Chromium for E2E testing
- Runs as non-root user

### sqlc.yaml
- Configures sqlc for SQL code generation
- Outputs to `internal/models`
- Uses pgx/v5 as SQL driver
- Generates Go structs and query functions

### .air.toml
- Configures Air for hot reloading
- Watches `.go` and `.sql` files
- Builds to `./tmp/main`
- Excludes test files and vendor directory

## Fixes Applied

### 1. Go Version Mismatch ✅
- **Issue**: Dockerfile.backend used Go 1.22, but go.mod specified Go 1.23
- **Fix**: Updated Dockerfile.backend to use `golang:1.23-alpine`

### 2. Port Conflict ✅
- **Issue**: Both app and postgres services exposed port 5432
- **Fix**: Removed port 5432 from app service (only postgres should expose it)

### 3. Unnecessary Services ✅
- **Issue**: Redis and Mailhog services were included but not required
- **Fix**: Removed both services and their volumes from docker-compose.yml

### 4. Missing sqlc Configuration ✅
- **Issue**: post-create.sh referenced sqlc.yaml which didn't exist
- **Fix**: Created sqlc.yaml with proper configuration for pgx/v5

### 5. Missing Air Configuration ✅
- **Issue**: Dockerfile.backend referenced .air.toml which didn't exist
- **Fix**: Created .air.toml with appropriate hot-reload settings

### 6. Redis References ✅
- **Issue**: post-create.sh waited for Redis and included Redis environment variables
- **Fix**: Removed Redis wait logic and Redis env vars from backend/.env template

### 7. Missing Dependencies ✅
- **Issue**: go.mod was missing PostgreSQL driver and JWT libraries
- **Fix**: Added pgx/v5, golang-jwt/jwt/v5, and clerk-sdk-go

## Environment Variables

### Backend (.env)
```
DB_HOST=postgres
DB_PORT=5432
DB_USER=budgetuser
DB_PASSWORD=budgetpass
DB_NAME=budgetdb
DB_SSL_MODE=disable
PORT=8080
ENVIRONMENT=development
ALLOWED_ORIGINS=http://localhost:5173
CLERK_SECRET_KEY=
CLERK_PUBLISHABLE_KEY=
CLERK_JWT_KEY=
SYNC_BATCH_SIZE=50
SYNC_RETRY_ATTEMPTS=3
SYNC_RETRY_DELAY=5s
LOG_LEVEL=debug
LOG_FORMAT=json
```

### Frontend (.env)
```
PUBLIC_CLERK_PUBLISHABLE_KEY=
PUBLIC_CLERK_SIGN_IN_URL=/sign-in
PUBLIC_CLERK_SIGN_UP_URL=/sign-up
PUBLIC_CLERK_AFTER_SIGN_IN_URL=/dashboard
PUBLIC_CLERK_AFTER_SIGN_UP_URL=/dashboard
PUBLIC_API_URL=http://localhost:8080/api
PUBLIC_APP_NAME="Budget Planner"
PUBLIC_APP_SHORT_NAME="Budget"
PUBLIC_APP_DESCRIPTION="Offline-first budget planning application"
PUBLIC_APP_THEME_COLOR="#3b82f6"
PUBLIC_APP_BACKGROUND_COLOR="#ffffff"
PUBLIC_SYNC_INTERVAL=30000
PUBLIC_OFFLINE_RETRY_DELAY=5000
PUBLIC_MAX_OFFLINE_OPERATIONS=100
```

## Development Workflow

### Starting the Devcontainer
1. Open the project in VS Code
2. Click "Reopen in Container" when prompted
3. Wait for post-create script to complete (~2-3 minutes)

### Running the Backend
```bash
cd backend
air
```
This will start the server with hot reload on port 8080.

### Running the Frontend
```bash
cd frontend
npm run dev
```
This will start the dev server with hot reload on port 5173.

### Generating SQL Code
```bash
cd backend
sqlc generate
```
This generates Go code from SQL queries in `sql/queries/`.

### Running Database Migrations
```bash
cd backend
migrate -path sql/schema -database "postgres://budgetuser:budgetpass@postgres:5432/budgetdb?sslmode=disable" up
```

## Database Access

### Via SQLTools Extension
- Pre-configured connection named "Budget Planner DB"
- Host: postgres
- Port: 5432
- Database: budgetdb
- User: budgetuser
- Password: budgetpass

### Via Command Line
```bash
psql -h postgres -U budgetuser -d budgetdb
```

## Troubleshooting

### Port Already in Use
If you see "port already in use" errors:
1. Check if another container is running: `docker ps`
2. Stop conflicting containers: `docker stop <container-name>`

### PostgreSQL Not Ready
If the app can't connect to PostgreSQL:
1. Check postgres health: `docker logs budget-planner-postgres`
2. Restart the devcontainer
3. Verify init-db.sql executed without errors

### Dependencies Not Installing
If Go or Node dependencies fail to install:
1. Check post-create.sh logs in terminal
2. Manually run installation commands
3. Check internet connectivity

### Hot Reload Not Working
If changes aren't reflected:
1. Ensure Air is running in the backend
2. Check that .air.toml is watching the right files
3. Verify file permissions are correct

## Next Steps

1. **Add Clerk API Keys**
   - Sign up at https://clerk.com
   - Create a new application
   - Copy keys to backend/.env and frontend/.env

2. **Create Initial SQL Queries**
   - Add query files to `backend/sql/queries/`
   - Run `sqlc generate` to create Go code

3. **Set Up Database Migrations**
   - Create migration files in `backend/sql/schema/`
   - Run migrations to apply changes

4. **Implement API Endpoints**
   - Create handlers in `backend/internal/handlers/`
   - Register routes in `backend/cmd/api/main.go`
   - Test endpoints with REST Client or Postman

5. **Build Frontend Components**
   - Use Shadcn-Svelte components
   - Set up IndexedDB for offline storage
   - Implement Clerk authentication
   - Create PWA manifest and service worker

## Additional Resources

- [Chi Router Documentation](https://github.com/go-chi/chi)
- [sqlc Documentation](https://docs.sqlc.dev/)
- [Shadcn-Svelte](https://www.shadcn-svelte.com)
- [SvelteKit Docs](https://kit.svelte.dev/)
- [Clerk SDK](https://clerk.com/docs)
- [Air Live Reload](https://github.com/cosmtrek/air)
