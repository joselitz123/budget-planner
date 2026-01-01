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
updated_at: 2026-01-01T13:51:00Z
completed_at: 2026-01-01T13:51:00Z
---

## Problem

Clerk publishable key environment variable is not being loaded despite multiple fixes attempted.

### Error Message

### Root Causes Identified

1. **Malformed .env file**: First line was missing `#` comment prefix
2. **Inconsistent env var prefixes**: Code used `VITE_` prefix but .env had `PUBLIC_` prefix
3. **Incorrect envDir configuration**: `vite.config.ts` had `envDir: './'` pointing to wrong directory
4. **Missing Vite configuration**: `vite.config.ts` lacked explicit `envPrefix` and `envDir` settings
5. **DOM timing issue**: Clerk mounting before container DOM element was rendered

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

1. **Added explicit Vite configuration to `frontend/vite.config.ts`**:

   - Added `envPrefix: 'PUBLIC_'` to tell Vite which variables to load
   - Added `envDir: './'` to point to correct directory

2. **Fixed DOM timing in auth pages**:

   - Updated `frontend/src/routes/(auth)/sign-in/+page.svelte` to use `await tick()` before mounting Clerk
   - Updated `frontend/src/routes/(auth)/sign-up/+page.svelte` to use `await tick()` before mounting Clerk
   - This ensures container DOM element exists before Clerk tries to mount

3. Added `API_URL` to `frontend/.env` for server-side use

4. Updated `frontend/.env.example` with `API_URL` example

5. Updated `frontend/src/routes/+layout.server.ts` to use `process.env.API_URL`

6. Added comprehensive logging to `frontend/src/routes/+layout.svelte`

7. Fixed `frontend/src/routes/profile/+page.svelte` to use dynamic import (resolves SSR issues)

8. Added detailed logging to `frontend/src/lib/auth/clerkProvider.ts`

9. Fixed `frontend/src/test/setup.ts` env var prefix from `VITE_PUBLIC_` to `PUBLIC_`

10. Created `frontend/ENV_VAR_FIX_SUMMARY.md` with comprehensive documentation

11. Created `frontend/test-env.html` for env var testing

12. Created `frontend/verify-env.ts` for console verification

### Root Cause (Final - RESOLVED ✅)

The issue was **two-fold**:

1. **Missing Vite Configuration**: [`vite.config.ts`](frontend/vite.config.ts:7-9) lacked explicit environment variable configuration

   - Missing `envPrefix: 'PUBLIC_'` to tell Vite which variables to load
   - Missing `envDir: './'` to point to the frontend directory where `.env` is located

2. **DOM Timing Issue**: Clerk's `mountSignIn()` and `mountSignUp()` were called before the container DOM element was rendered
   - The `{#if loading}` block prevented the container from being in the DOM when Clerk tried to mount
   - Solution: Set `loading = false` before mounting, then use `await tick()` to wait for DOM update

### Final Fix Applied

**1. Updated `frontend/vite.config.ts`:**

```typescript
export default defineConfig({
    // Explicitly configure environment variable loading
    envPrefix: 'PUBLIC_',
    envDir: './',
    plugins: [
        sveltekit(),
```

**2. Updated `frontend/src/routes/(auth)/sign-in/+page.svelte`:**

- Imported `tick` from `svelte`
- Set `loading = false` before mounting Clerk
- Added `await tick()` to wait for DOM update before calling `clerk.mountSignIn()`

**3. Updated `frontend/src/routes/(auth)/sign-up/+page.svelte`:**

- Imported `tick` from `svelte`
- Set `loading = false` before mounting Clerk
- Added `await tick()` to wait for DOM update before calling `clerk.mountSignUp()`

### Verification

**Environment Variables Now Loading Correctly:**

```
✅ PUBLIC_CLERK_PUBLISHABLE_KEY
✅ PUBLIC_DEV_AUTH_PROVIDER
✅ PUBLIC_API_URL
✅ PUBLIC_APP_NAME
✅ PUBLIC_APP_SHORT_NAME
```

**Clerk Initialization Successful:**

```
✅ [Clerk] Successfully initialized
✅ Sign-in form renders correctly
✅ Sign-up form renders correctly
✅ No "empty element" errors
```

**Status**: ✅ **RESOLVED** - Clerk environment variables load correctly and sign-in/sign-up forms work

## Impact

**Critical** - Was blocking all authentication functionality

- Users can now sign in/sign up
- App can initialize Clerk
- Protected routes are accessible after authentication
- JWT tokens can be retrieved from Clerk session

**Resolution**: Authentication now works correctly

## Session Date

2026-01-01

## Files Modified (Session 2026-01-01 13:20-13:51)

- frontend/.env (rewritten to ensure correct format)
- frontend/vite.config.ts (added envPrefix and envDir configuration)
- frontend/src/routes/(auth)/sign-in/+page.svelte (added tick import and DOM timing fix)
- frontend/src/routes/(auth)/sign-up/+page.svelte (added tick import and DOM timing fix)
- frontend/test-env.html (created for env var testing)
- frontend/verify-env.ts (created for console verification)
- .beans/BP-rpwv--clerk-environment-variables-not-loading-despite-fi.md (updated)

## Files Modified (Session 2026-01-01 13:00-13:25)

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
