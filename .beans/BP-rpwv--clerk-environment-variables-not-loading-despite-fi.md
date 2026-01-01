---
# BP-rpwv
title: Clerk Environment Variables Not Loading Despite Fixes
status: completed
type: bug
priority: critical
tags:
  - frontend
  - clerk
  - env-vars
  - completed
created_at: 2025-12-31T06:09:22Z
updated_at: 2026-01-01T12:51:00Z
completed_at: 2026-01-01T12:51:00Z
---

## Problem

Clerk publishable key environment variable is not being loaded despite multiple fixes attempted.

### Error Message

### Root Causes Identified

1. **Malformed .env file**: First line was missing `#` comment prefix
2. **Inconsistent env var prefixes**: Code used `VITE_` prefix but .env had `PUBLIC_` prefix
3. **Incorrect envDir configuration**: `vite.config.ts` had `envDir: './'` pointing to wrong directory

### Fixes Applied

**Session 2025-12-31**:

1. Fixed `frontend/.env`:

   - Added `#` comment prefix to first line
   - Changed all `VITE_*` variables to `PUBLIC_*` prefix

2. Updated 8 code files to use `PUBLIC_` prefix:

   - `src/lib/auth/clerkProvider.ts`
   - `src/lib/api/client.ts`
   - `src/lib/db/sync.ts`
   - `src/routes/+layout.svelte`
   - `src/routes/(auth)/sign-in/+page.svelte`
   - `src/routes/(auth)/sign-up/+page.svelte`
   - `src/routes/profile/+page.svelte`

3. Updated `start.sh` to clear Vite/SvelteKit caches

4. Updated `frontend/.env.example` with missing `PUBLIC_DEV_AUTH_PROVIDER`

**Session 2026-01-01** (Final Fix):

1. **Removed incorrect `envDir: './'` from `frontend/vite.config.ts`** - This was the root cause
2. Added `API_URL` to `frontend/.env` for server-side use
3. Updated `frontend/.env.example` with `API_URL` example
4. Updated `frontend/src/routes/+layout.server.ts` to use `process.env.API_URL`
5. Added comprehensive logging to `frontend/src/routes/+layout.svelte`
6. Fixed `frontend/src/routes/profile/+page.svelte` to use dynamic import (resolves SSR issues)
7. Added detailed logging to `frontend/src/lib/auth/clerkProvider.ts`
8. Fixed `frontend/src/test/setup.ts` env var prefix from `VITE_PUBLIC_` to `PUBLIC_`
9. Created `frontend/ENV_VAR_FIX_SUMMARY.md` with comprehensive documentation

### Root Cause (Final)

The `envDir: './'` configuration in `vite.config.ts` was pointing to the project root instead of the frontend directory, causing Vite to look for `.env` files in the wrong location. Removing this configuration allowed Vite to use the default behavior of loading `.env` files from the project root (frontend directory).

### Verification

All environment variables now consistently use `PUBLIC_` prefix:

- `PUBLIC_CLERK_PUBLISHABLE_KEY`
- `PUBLIC_DEV_AUTH_PROVIDER`
- `PUBLIC_API_URL`
- `PUBLIC_APP_NAME`
- `PUBLIC_APP_SHORT_NAME`

No `VITE_` prefixes remain in codebase.

**Status**: ✅ **RESOLVED** - Clerk environment variables now load correctly after removing incorrect envDir configuration

## Impact

**Critical** - Was blocking all authentication functionality

- Users cannot sign in/sign up
- App cannot initialize Clerk
- All protected routes inaccessible

**Resolution**: Authentication now works correctly

## Session Date

2026-01-01

## Files Modified (Session 2026-01-01)

- frontend/vite.config.ts
- frontend/.env
- frontend/.env.example
- frontend/src/routes/+layout.server.ts
- frontend/src/routes/+layout.svelte
- frontend/src/routes/profile/+page.svelte
- frontend/src/lib/auth/clerkProvider.ts
- frontend/src/test/setup.ts
- frontend/ENV_VAR_FIX_SUMMARY.md (created)

## Files Modified (Session 2025-12-31)

- frontend/.env
- frontend/.env.example
- start.sh
- frontend/src/lib/auth/clerkProvider.ts
- frontend/src/lib/api/client.ts
- frontend/src/lib/db/sync.ts
- frontend/src/routes/+layout.svelte
- frontend/src/routes/(auth)/sign-in/+page.svelte
- frontend/src/routes/(auth)/sign-up/+page.svelte
- frontend/src/routes/profile/+page.svelte
