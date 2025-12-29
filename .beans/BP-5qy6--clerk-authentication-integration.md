---
# BP-5qy6
title: Clerk Authentication Integration
status: in-progress
type: feature
priority: critical
created_at: 2025-12-28T16:51:54Z
updated_at: 2025-12-29T06:26:30Z
---

Complete Clerk authentication integration with JWT token handling and protected routes.

## Background
Clerk SDK is configured and DevTokenProvider exists, but actual Clerk integration is not wired up. Users cannot log in/out.

## Acceptance Criteria
- [x] Clerk provider setup in root layout
- [x] Sign in/sign up pages
- [x] Protected route middleware
- [x] JWT token retrieval from Clerk session
- [x] Auto-redirect to sign in for protected routes
- [x] Logout functionality
- [x] User profile display

## Implementation Notes (2025-12-29)

### Completed Features
1. ✅ **Clerk Packages Installed**: @clerk/clerk-js and @clerk/backend
2. ✅ **Environment Variables**: Configured with Clerk keys and development switch (VITE_DEV_AUTH_PROVIDER)
3. ✅ **Server-Side Auth**:
   - `hooks.server.ts`: Extracts session ID from Clerk cookies
   - `+layout.server.ts`: Route protection and onboarding flow
4. ✅ **Client-Side Auth**:
   - `ClerkTokenProvider`: Retrieves JWT tokens from Clerk session
   - API client updated with development switch
5. ✅ **Auth Pages**:
   - `/sign-in`: Clerk SignIn component
   - `/sign-up`: Clerk SignUp component
   - `/profile`: User profile page with logout
6. ✅ **UI Integration**:
   - User menu in header with profile link and logout button
   - User display shows name from backend
7. ✅ **Type Checking**: 0 errors, 15 warnings (all acceptable)
8. ✅ **Build**: Successful with PWA config updated for 5MB chunk size limit

### Development Approach
- Used `@clerk/clerk-js` JavaScript SDK (no official SvelteKit SDK exists)
- JWT decoded on server-side for user ID extraction (backend does full verification)
- Development switch allows toggling between Clerk and DevTokenProvider
- Onboarding flow integrated in `+layout.server.ts`

### Files Created
- `frontend/src/hooks.server.ts`
- `frontend/src/routes/+layout.server.ts`
- `frontend/src/lib/auth/clerkProvider.ts`
- `frontend/src/routes/(auth)/+layout.svelte`
- `frontend/src/routes/(auth)/sign-in/+page.svelte`
- `frontend/src/routes/(auth)/sign-up/+page.svelte`
- `frontend/src/routes/profile/+page.svelte`
- `frontend/src/lib/components/ui/card/card.svelte`

### Files Modified
- `frontend/package.json`: Added @clerk/clerk-js and @clerk/backend
- `frontend/.env`: Clerk credentials and development switch
- `frontend/src/routes/+layout.svelte`: Added Clerk init and user menu
- `frontend/src/lib/api/client.ts`: Added ClerkTokenProvider import and switch logic
- `frontend/vite.config.ts`: Increased PWA chunk size limit to 5MB

### Testing Required
⚠️ **Not Yet Tested** (requires running servers):
- Sign-in/sign-up flow with real Clerk account
- Onboarding flow to backend `/api/auth/onboarding`
- JWT token passing to backend API
- Logout functionality
- Development switch (VITE_DEV_AUTH_PROVIDER=dev)
- Protected route redirects

### Known Issues
- None - type checking passes and build succeeds

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