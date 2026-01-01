---
# BP-jwks
title: Fix JWT algorithm incompatibility - Implement JWKS verification for Clerk RS256 tokens
status: completed
type: bug
priority: critical
tags:
  - backend
  - authentication
  - jwt
  - clerk
  - jwks
  - critical
created_at: 2026-01-01T14:47:00Z
updated_at: 2026-01-01T14:47:00Z
completed_at: 2026-01-01T14:47:00Z
---

## Problem

Backend was using HMAC-SHA256 (HS256) to verify Clerk's RS256-signed tokens. These algorithms are incompatible - cannot verify RS256 tokens with HMAC secret. All API requests with Clerk tokens failed with 401 Unauthorized.

### Root Cause

The authentication middleware in `backend/internal/auth/clerk.go` was attempting to verify Clerk's RS256-signed JWT tokens using HMAC-SHA256 (HS256) algorithm with a JWT_SECRET environment variable. This is fundamentally incompatible because:

- **RS256**: Uses RSA private key for signing and public key for verification
- **HS256**: Uses shared secret for both signing and verification

Clerk signs tokens with RS256 using their private keys, but the backend was trying to verify them with HS256 using a secret key, which always fails.

### Impact

**Critical** - All authenticated API endpoints were returning 401 Unauthorized errors, blocking all user interactions with the backend after authentication.

## Solution

Implemented JWKS (JSON Web Key Set) based verification using the keyfunc library. Backend now fetches Clerk's public keys from their JWKS endpoint and verifies RS256 signatures correctly.

### Implementation Details

1. **Added keyfunc dependency** to `backend/go.mod`:

   - `github.com/MicahParks/keyfunc v1.9.0`

2. **Completely rewrote `backend/internal/auth/clerk.go`**:

   - Removed HS256 verification logic
   - Implemented JWKS-based verification
   - Added automatic key refresh from Clerk's JWKS endpoint
   - Proper error handling for key retrieval and verification

3. **Updated configuration** in `backend/internal/config/config.go`:

   - Added `CLERK_DOMAIN` environment variable
   - Removed dependency on `JWT_SECRET`

4. **Updated initialization** in `backend/cmd/api/main.go`:

   - Modified authentication initialization to use new JWKS verifier

5. **Updated environment files**:

   - `backend/.env` - Added CLERK_DOMAIN, removed JWT_SECRET
   - `backend/.env.example` - Updated documentation

6. **Updated documentation** in `backend/CLAUDE.md`

### Key Changes

**Before (HS256):**

```go
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(config.JWTSecret), nil
})
```

**After (JWKS with RS256):**

```go
jwks, err := keyfunc.Get(config.ClerkDomain + "/.well-known/jwks.json", keyfunc.Options{})
token, err := jwt.Parse(tokenString, jwks.Keyfunc)
```

## Files Modified

1. `backend/internal/auth/clerk.go` - Completely rewritten with JWKS verification
2. `backend/internal/config/config.go` - Added CLERK_DOMAIN configuration
3. `backend/cmd/api/main.go` - Updated authentication initialization
4. `backend/.env` - Added CLERK_DOMAIN, removed JWT_SECRET
5. `backend/.env.example` - Updated documentation
6. `backend/go.mod` - Added keyfunc v1.9.0 dependency
7. `backend/CLAUDE.md` - Updated authentication documentation

## Testing

### Verification Steps

1. ✅ Backend starts without errors
2. ✅ JWKS initialization successful (fetches Clerk's public keys)
3. ✅ Authentication middleware correctly validates RS256 tokens
4. ✅ Protected endpoints properly secured and accessible with valid tokens
5. ✅ Invalid tokens correctly rejected with 401 Unauthorized

### Test Results

- Backend successfully initializes JWKS key provider
- Clerk tokens from frontend are now verified correctly
- API endpoints return 200 OK for authenticated requests
- Invalid/expired tokens are properly rejected

## Impact

**Critical** - This fix unblocks all authenticated API functionality.

- Users can now successfully make authenticated requests to the backend
- All protected endpoints (transactions, budgets, categories, etc.) are accessible
- JWT token verification is cryptographically secure using RS256
- Automatic key refresh ensures long-term reliability

## Session Date

2026-01-01

## Effort Estimate

2 hours

## Related Issues

- This was a prerequisite for all other authentication-dependent features
- Required understanding of JWT algorithms and JWKS specification
- No breaking changes to API contracts or frontend code
