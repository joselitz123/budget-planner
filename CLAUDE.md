# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Budget Planner is an **offline-first Progressive Web App (PWA)** with a SvelteKit frontend and Go backend. The app focuses on personal budget management with monthly reflections, budget sharing, and offline synchronization.

### Tech Stack

**Frontend:** SvelteKit PWA + Shadcn-Svelte + Tailwind CSS + IndexedDB (offline storage)
**Backend:** Go 1.23 + Chi router + PostgreSQL + sqlc (ORM)
**Auth:** Clerk SDK
**Development:** Docker Compose devcontainer with PostgreSQL 16

---

## Common Development Commands

### Backend (Go)

**Quick Start:**
```bash
cd backend

# Run with hot reload
air

# Run tests
DATABASE_URL="postgresql://budgetuser:budgetpass@localhost:5432/budgetdb?sslmode=disable" go test ./internal/handlers/ -v

# Build
go build -o ./tmp/main ./cmd/api
```

**For detailed backend documentation, see `backend/CLAUDE.md`** which includes:
- Complete architecture patterns (sqlc, type conversions, Chi router)
- Handler reference with all endpoints
- Testing guide with infrastructure details
- Type conversion cheat sheet
- Development workflow and debugging patterns

### Frontend (SvelteKit)

**Quick Start:**
```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm run dev

# Type checking
npm run check

# Build for production
npm run build
```

**For detailed frontend documentation, see `frontend/CLAUDE.md`** which includes:
- SvelteKit architecture and routing patterns
- PWA configuration and service workers
- IndexedDB for offline storage
- Shadcn-Svelte components
- State management with stores
- Clerk authentication integration
- Environment configuration

### Database

```bash
# Access PostgreSQL via psql (in devcontainer)
psql -h postgres -U budgetuser -d budgetdb

# Run migration
migrate -path backend/sql/schema -database "postgres://budgetuser:budgetpass@postgres:5432/budgetdb?sslmode=disable" up
```

---

## Architecture

### Offline-First Sync Architecture

The application is designed to work offline. Data flows between frontend IndexedDB and backend PostgreSQL:

1. **Frontend** stores all data in IndexedDB for offline access
2. **Sync queue** captures changes when offline
3. **Sync API** (`/api/sync/push`, `/api/sync/pull`) handles bidirectional sync
4. **Conflict resolution** uses timestamp-based and owner-priority strategies
5. **sync_operations table** tracks pending operations on the server

### Backend Architecture

**See `backend/CLAUDE.md` for detailed backend architecture including:**
- sqlc query generation workflow
- Database schema patterns
- Type conversion patterns (pgtype ↔ Go)
- Chi router path parameter handling
- Complete handler reference

**Quick Overview:**
- Go 1.23 + Chi v5 router + PostgreSQL + sqlc
- SQL queries written in `.sql` files → Go code generated via `sqlc generate`
- All 10 handlers implemented with 100% test coverage (48/48 tests passing)
- JWT-based authentication (Clerk-compatible)

### Frontend Architecture

**See `frontend/CLAUDE.md` for detailed frontend architecture including:**
- SvelteKit file-based routing patterns
- PWA configuration and service workers
- IndexedDB for offline storage
- Shadcn-Svelte component usage
- State management with Svelte stores
- Clerk authentication integration
- API integration patterns

**Quick Overview:**
- SvelteKit 2.0 + Svelte 5 + TypeScript 5
- Shadcn-Svelte components + Tailwind CSS
- IndexedDB (via `idb`) for offline storage
- Clerk for authentication
- PWA with service workers for offline functionality

### Frontend Structure

```
frontend/
├── src/
│   ├── lib/                # Utility functions, DB clients
│   ├── routes/             # SvelteKit file-based routing
│   └── app.html            # HTML template
├── static/                 # PWA assets, icons
└── svelte.config.js        # Svelte configuration
```

---

## Key Implementation Notes

### Authentication Flow

1. Frontend uses Clerk SDK for authentication
2. Backend validates JWT tokens via middleware
3. User context attached to requests
4. Clerk user ID stored in `users.clerk_user_id` for lookups

**For detailed backend authentication and permission model, see `backend/CLAUDE.md`**

### Migrations

The project uses numbered migration files (e.g., `001_initial_schema.up.sql`):
- Up migrations apply changes
- Down migrations roll back changes
- The init-db.sql in `.devcontainer/` contains the full initial schema

When adding migrations:
1. Create new file with next sequence number
2. Write both up and down versions
3. Test with `migrate up` and `migrate down`

### Hot Reload

- Backend uses Air (configured in `.air.toml`) - watches `.go` and `.sql` files
- Frontend uses Vite dev server - auto-reloads on save
- Running `sqlc generate` is manual - run after modifying `.sql` query files

---

## Testing

### Backend Testing

**Quick Start:**
```bash
cd backend

# Run all tests
DATABASE_URL="postgresql://budgetuser:budgetpass@localhost:5432/budgetdb?sslmode=disable" go test ./internal/handlers/ -v

# Run specific test
DATABASE_URL="postgresql://budgetuser:budgetpass@localhost:5432/budgetdb?sslmode=disable" go test ./internal/handlers/auth_test.go -v
```

**Current Status:** 48/48 tests passing (100%) ✅

**For detailed testing infrastructure, type conversions, and test patterns, see `backend/CLAUDE.md`**

---

## Environment Configuration

Required environment variables:

**Backend (.env):**
```
DATABASE_URL=postgresql://budgetuser:budgetpass@postgres:5432/budgetdb?sslmode=disable
PORT=8080
CLERK_SECRET_KEY=your_clerk_secret_key
CLERK_PUBLISHABLE_KEY=your_clerk_publishable_key
```

**Frontend (.env):**
```
PUBLIC_CLERK_PUBLISHABLE_KEY=your_clerk_publishable_key
PUBLIC_API_URL=http://localhost:8080/api
```

---

## Important Files

- `starting-point.md` - Full project specification with database schema, API endpoints, and workflows
- `backend/CLAUDE.md` - **Backend development guide** with architecture patterns, testing, and type conversions
- `frontend/CLAUDE.md` - **Frontend development guide** with SvelteKit patterns, PWA configuration, and IndexedDB usage
- `frontend/todo.md` - **Frontend implementation status & TODO** with detailed task list and technical debts
- `backend/TASK.md` - Backend development task log and session history
- `AGENTS.md` - AI agent specifications for automation (reference for future development)
- `.devcontainer/SETUP_GUIDE.md` - Devcontainer setup details and troubleshooting

---

## Recent Progress (Dec 27, 2025)

### ✅ Completed: Transaction Modal Implementation

**Frontend (SvelteKit) - 75% Complete**

The Transaction Modal feature has been successfully implemented, enabling users to add expenses through a functional form interface.

**What Was Built:**
- Custom UI components: Button, Input, Label, Textarea, Select, Badge, CustomModal
- AddExpenseModal with 8 form fields (amount, date, description, category, recurring, due date, notes)
- Budget auto-creation functionality (creates $2000 default budget if none exists)
- Form validation with inline error messages
- Toast notifications for success/error feedback
- FAB button integration to open modal

**Quality Assurance:**
- ✅ Type checking: 0 errors, 2 accessibility warnings (non-blocking)
- ✅ Build successful
- ✅ Manual testing complete

**Technical Notes:**
- Shadcn-Svelte CLI unavailable (TTY requirement) - components created manually
- bits-ui Dialog had type definition issues - created CustomModal instead
- Simplified type definitions used to avoid svelte/elements import conflicts
- See `frontend/todo.md` for detailed technical debts and workarounds

**Current Branch:** `iteration/1`

**Next Steps:**
1. Backend API Integration - Wire up actual API calls to Go backend
2. Mark Bill Paid Functionality - Implement button click handler
3. Month Navigation - Enable prev/next month buttons

**Files Modified:**
- `frontend/components.json` - Shadcn configuration
- `frontend/src/lib/components/ui/*` - UI components
- `frontend/src/routes/transactions/AddExpenseModal.svelte` - Transaction form
- `frontend/src/routes/transactions/+page.svelte` - FAB button wiring
- `frontend/src/lib/stores/budgets.ts` - Budget auto-creation function
- `frontend/package.json` - Added bits-ui, tailwind-variants dependencies

---
