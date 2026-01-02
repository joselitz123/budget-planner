---
# BP-pkey
title: Clerk publishable key missing in server-side authentication
status: completed
type: bug
priority: critical
tags:
  - frontend
  - authentication
  - clerk
  - env-vars
created_at: 2026-01-02T04:36:00Z
updated_at: 2026-01-02T04:36:00Z
---

## Problem

Server-side Clerk client initialization is missing the PUBLISHABLE_KEY, causing all server-side authentication to fail. The `authenticateRequest()` function requires both SECRET_KEY and PUBLISHABLE_KEY to be provided to `createClerkClient()`.

### Error Message

```
[Clerk] Session verification failed: Error: Publishable key is missing. Ensure that your publishable key is correctly configured. Double-check your environment configuration for your keys, or access them here: https://dashboard.clerk.com/last-active?path=api-keys
    at parsePublishableKey (file:///workspace/budget-planner/frontend/node_modules/@clerk/shared/dist/runtime/keys-YNv6yjKk.mjs:53:36)
    at assertValidPublishableKey (file:///workspace/budget-planner/frontend/node_modules/@clerk/backend/dist/chunk-TTM76E4X.mjs:208:3)
    at AuthenticateContext.initPublishableKeyValues (file:///workspace/budget-planner/frontend/node_modules/@clerk/backend/dist/chunk-TTM76E4X.mjs:343:5)
```

## Context

The current implementation in `frontend/src/hooks.server.ts` only passes the SECRET_KEY to `createClerkClient()`, but Clerk's server-side authentication requires both SECRET_KEY and PUBLISHABLE_KEY to be configured. This causes any server-side authentication request to fail with the error above.

The environment variable for the publishable key exists as `CLERK_PUBLISHABLE_KEY` but is not being passed to the Clerk client initialization.

## Reproduction Steps

1. Start the application with Clerk configured
2. Attempt to access a protected route or make an authenticated server-side request
3. Observe the error: "Publishable key is missing"
4. All server-side authentication fails

## Expected Behavior

The Clerk client should be initialized with both SECRET_KEY and PUBLISHABLE_KEY, allowing server-side authentication to work correctly.

## Actual Behavior

The Clerk client is initialized with only SECRET_KEY, causing all server-side authentication attempts to fail with a missing publishable key error.

## Root Cause

In `frontend/src/hooks.server.ts`, the `createClerkClient()` function is called with only the `secretKey` parameter:

```typescript
const clerk = createClerkClient({
  secretKey: process.env.CLERK_SECRET_KEY,
});
```

However, Clerk's server-side authentication requires both keys:

```typescript
const clerk = createClerkClient({
  secretKey: process.env.CLERK_SECRET_KEY,
  publishableKey: process.env.CLERK_PUBLISHABLE_KEY, // Missing!
});
```

The `authenticateRequest()` method internally validates that both keys are present and throws an error if the publishable key is missing.

## Proposed Solution

Update `frontend/src/hooks.server.ts` to include the publishable key when creating the Clerk client:

1. Add `CLERK_PUBLISHABLE_KEY` to the Clerk client initialization
2. Ensure the environment variable is properly loaded from the backend `.env` file
3. Test server-side authentication flows

Example fix:

```typescript
const clerk = createClerkClient({
  secretKey: process.env.CLERK_SECRET_KEY,
  publishableKey: process.env.CLERK_PUBLISHABLE_KEY,
});
```

Note: The publishable key needs to be available in the server environment. This may require:

- Adding `CLERK_PUBLISHABLE_KEY` to `backend/.env`
- Loading it in the server environment (not just frontend)
- Ensuring it's passed to the SvelteKit server hooks

## Files Affected

- `frontend/src/hooks.server.ts` - Add publishableKey to createClerkClient() call
- `backend/.env` - May need to add CLERK_PUBLISHABLE_KEY if not present
- `backend/.env.example` - Add CLERK_PUBLISHABLE_KEY example

## Related Issues

- BP-rpwv: Clerk environment variables not loading despite fixes (resolved)
- BP-jwks: Fix JWT algorithm incompatibility (related to Clerk authentication)

## Impact

**Critical** - Blocks all server-side authentication functionality, including:

- Protected route access
- Server-side API authentication
- Session verification on the server
- Any server-side Clerk operations

This completely breaks the authentication flow for server-side operations.

## Effort Estimate

15 minutes

## Session Date

2026-01-02

## Fix Applied

**Status:** Completed on 2026-01-02

### Files Modified

- [`frontend/src/hooks.server.ts`](frontend/src/hooks.server.ts)

### Fix Description

Fixed SvelteKit import to use both `$env/dynamic/private` and `$env/dynamic/public` environment modules, and updated `createClerkClient()` to use the correct env objects for both SECRET_KEY and PUBLISHABLE_KEY.

### Implementation

```typescript
import { env as envPrivate } from "$env/dynamic/private";
import { env as envPublic } from "$env/dynamic/public";

const clerk = createClerkClient({
  secretKey: envPrivate.CLERK_SECRET_KEY,
  publishableKey: envPublic.CLERK_PUBLISHABLE_KEY,
});
```

The fix ensures that:

1. The SECRET_KEY is loaded from private environment variables (backend/server-side)
2. The PUBLISHABLE_KEY is loaded from public environment variables (client-side accessible)
3. Both keys are properly passed to the Clerk client initialization
4. Server-side authentication now works correctly with both keys configured
