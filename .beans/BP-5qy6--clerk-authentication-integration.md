---
# BP-5qy6
title: Clerk Authentication Integration
status: open
type: feature
priority: critical
created_at: 2025-12-28T16:51:54Z
updated_at: 2025-12-28T16:51:54Z
---

Complete Clerk authentication integration with JWT token handling and protected routes.

## Background
Clerk SDK is configured and DevTokenProvider exists, but actual Clerk integration is not wired up. Users cannot log in/out.

## Acceptance Criteria
- [ ] Clerk provider setup in root layout
- [ ] Sign in/sign up pages
- [ ] Protected route middleware
- [ ] JWT token retrieval from Clerk session
- [ ] Auto-redirect to sign in for protected routes
- [ ] Logout functionality
- [ ] User profile display

## Technical Details

### Clerk Integration Points
1. Wrap app with ClerkProvider
2. Use clerkClient.session.getToken() for JWT
3. Create ClerkAuthTokenProvider implementing AuthTokenProvider
4. Protect routes with +layout.server.ts or hooks

### Auth Flow
- Unauthenticated users: redirect to /sign-in
- After sign in: redirect to /dashboard
- Get JWT token from Clerk session
- Pass token to ApiClient for backend requests

### Routes to Create
- /sign-in (+page.svelte with Clerk SignIn component)
- /sign-up (+page.svelte with Clerk SignUp component)
- /profile (user profile page)

### Files to Create
- frontend/src/routes/(auth)/+layout.svelte (Clerk provider)
- frontend/src/routes/(auth)/sign-in/+page.svelte
- frontend/src/routes/(auth)/sign-up/+page.svelte
- frontend/src/lib/auth/clerkProvider.ts (Clerk token provider)

### Files to Modify
- frontend/src/routes/+layout.svelte (wrap with Clerk provider)
- frontend/src/lib/api/client.ts (use ClerkAuthTokenProvider)
- frontend/src/hooks.server.ts (add Clerk middleware)

### Environment Variables
- PUBLIC_CLERK_PUBLISHABLE_KEY (already set)
- CLERK_SECRET_KEY (backend, already set)

## Effort Estimate
2 hours

## Dependencies
- BP-r94p (Backend API Integration) - completed

## Session Date
2025-12-28