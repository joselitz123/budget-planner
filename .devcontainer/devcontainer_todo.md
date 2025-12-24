# Devcontainer Setup and Fix Log

This document tracks the devcontainer setup fixes and configurations for the Budget Planner project.

**Date Created:** 2025-12-24
**Status:** Fixed and Ready for Rebuild

---

## Summary of Fixes Applied

### 1. Git Configuration Issue (Fixed)

**Problem:**
- `/home/vscode/.gitconfig` was created as a directory instead of a file
- This broke all git operations with error: "warning: unable to access '/home/vscode/.gitconfig': Is a directory"
- Caused by docker volume mount: `~/.gitconfig.user:/home/vscode/.gitconfig:ro` in docker-compose.yml

**Solution:**
- Removed the problematic `.gitconfig.user` volume mount from docker-compose.yml
- Added `remoteEnv` to devcontainer.json to set proper git config environment variables:
  ```json
  "remoteEnv": {
    "GIT_CONFIG_GLOBAL": "${containerEnv:HOME}/.gitconfig_container",
    "GIT_CONFIG_NOSYSTEM": "1"
  }
  ```
- Updated post-create.sh to:
  - Remove any incorrect `.gitconfig` directory
  - Create proper `.gitconfig_container` file
  - Use `GIT_CONFIG_GLOBAL` and `GIT_CONFIG_NOSYSTEM` for all git commands

**Files Modified:**
- `.devcontainer/docker-compose.yml` - Removed .gitconfig.user volume mount
- `.devcontainer/devcontainer.json` - Added remoteEnv for git config
- `.devcontainer/post-create.sh` - Added git config fix logic

---

### 2. GLM Coding Helper Integration (Added)

**Purpose:**
- Enable Claude Code integration with Z.AI's GLM Coding Plan
- Provides cheaper alternative to Anthropic's Claude API
- Uses `@z_ai/coding-helper` npm package for setup

**Implementation:**
- Added GLM Coding Helper installation to post-create.sh
- Installs globally via npm: `npm install -g @z_ai/coding-helper`
- Provides user with setup instructions in post-create output

**Usage After Setup:**
```bash
# Run the coding helper to configure Claude Code
npx @z_ai/coding-helper
# or
chelper
```

**Documentation References:**
- [Z.AI Claude Code Setup](https://docs.z.ai/devpack/tool/claude)
- [Coding Tool Helper](https://docs.z.ai/devpack/tool/coding-tool-helper)
- [Quick Start Guide](https://docs.z.ai/devpack/quick-start)
- [NPM Package](https://www.npmjs.com/package/@z_ai/coding-helper)

**Files Modified:**
- `.devcontainer/post-create.sh` - Added GLM Coding Helper installation

---

## Configuration Files Reference

### devcontainer.json
Key settings:
- Base image: `mcr.microsoft.com/devcontainers/base:ubuntu-22.04`
- Features: Go, Node.js, Docker-in-Docker
- Remote Environment Variables for git config fix
- VS Code extensions for Go, Svelte, SQL tools, etc.

### docker-compose.yml
Services:
- **app**: Main development container
- **postgres**: PostgreSQL 16 database
- **Volumes**: postgres-data, docker-data, node_modules, go-pkg, go-mod

### post-create.sh
Automated setup tasks:
1. Fix git configuration issue
2. Install GLM Coding Helper
3. Install Go dependencies (sqlc, migrate, air)
4. Install Node.js dependencies
5. Create environment files (.env)
6. Wait for PostgreSQL to be ready
7. Generate development certificates

---

## Environment Variables

### Git Config (Fixed via remoteEnv)
```bash
GIT_CONFIG_GLOBAL=/home/vscode/.gitconfig_container
GIT_CONFIG_NOSYSTEM=1
```

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

---

## Next Steps After Rebuild

### 1. Configure GLM Coding Plan (Claude Code)
```bash
npx @z_ai/coding-helper
```
Follow the prompts to set up your Z.AI API key.

### 2. Configure Clerk Authentication
- Sign up at https://clerk.com
- Create a new application
- Copy keys to:
  - `backend/.env`: `CLERK_SECRET_KEY`, `CLERK_PUBLISHABLE_KEY`, `CLERK_JWT_KEY`
  - `frontend/.env`: `PUBLIC_CLERK_PUBLISHABLE_KEY`

### 3. Verify Git is Working
```bash
git remote -v
git status
```
Should show no errors about .gitconfig.

### 4. Start Development Servers
```bash
# Backend (in /workspace/backend)
air

# Frontend (in /workspace/frontend)
npm run dev
```

### 5. Access Database
```bash
psql -h postgres -U budgetuser -d budgetdb
```
Or use the SQLTools extension in VS Code.

---

## Troubleshooting

### Git Issues
**Problem:** "unable to access '/home/vscode/.gitconfig': Is a directory"

**Solution:** The devcontainer rebuild should fix this. The `.gitconfig` directory is removed and replaced with a proper file at `$HOME/.gitconfig_container`.

**Workaround (if needed):**
```bash
GIT_CONFIG_GLOBAL=/home/vscode/.gitconfig_container GIT_CONFIG_NOSYSTEM=1 git <command>
```

### GLM Coding Helper Issues
**Problem:** `chelper` command not found

**Solution:**
```bash
npm install -g @z_ai/coding-helper
```

**Documentation:** https://docs.z.ai/devpack/tool/coding-tool-helper

### PostgreSQL Connection Issues
**Problem:** Cannot connect to PostgreSQL

**Solution:**
1. Check if postgres container is running: `docker ps`
2. Check logs: `docker logs budget-planner-postgres`
3. Verify connection string in backend/.env

---

## Additional Resources

### Project Documentation
- `CLAUDE.md` - Project overview and development commands
- `starting-point.md` - Full project specification
- `SETUP_GUIDE.md` - Detailed setup documentation

### External Documentation
- [Z.AI DevPack Quick Start](https://docs.z.ai/devpack/quick-start)
- [Claude Code Documentation](https://docs.z.ai/devpack/tool/claude)
- [Chi Router](https://github.com/go-chi/chi)
- [sqlc](https://docs.sqlc.dev/)
- [Shadcn-Svelte](https://www.shadcn-svelte.com)
- [SvelteKit](https://kit.svelte.dev/)

---

## Rebuild Instructions

To apply all fixes, rebuild the devcontainer:

1. In VS Code: `F1` → "Dev Containers: Rebuild Container"
2. Or press `Ctrl+Shift+P` → "Dev Containers: Rebuild Container"

The rebuild will:
- Remove the problematic `.gitconfig` directory
- Create proper git configuration at `.gitconfig_container`
- Install GLM Coding Helper
- Set up all dependencies automatically

---

## Change History

| Date | Description | Files Modified |
|------|-------------|----------------|
| 2025-12-24 | Initial devcontainer fixes | docker-compose.yml, devcontainer.json, post-create.sh |
| 2025-12-24 | Added GLM Coding Helper integration | post-create.sh |
| 2025-12-24 | Created this documentation | devcontainer_todo.md |
