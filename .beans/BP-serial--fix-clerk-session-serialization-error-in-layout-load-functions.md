# BP-serial: Fix Clerk Session Serialization Error

**Type:** bug  
**Status:** completed  
**Priority:** critical  
**Effort:** 15 minutes  
**Tags:** frontend, authentication, clerk, serialization, sveltekit  
**Created:** 2026-01-02  
**Completed:** 2026-01-02

## Problem

Application was failing to load with the following error:

```
Error: Data returned from `load` while rendering / is not serializable: Cannot stringify a function (data.session.getToken)
```

The issue occurred in both:

- `frontend/src/routes/+layout.server.ts`
- `frontend/src/routes/(auth)/+layout.server.ts`

## Root Cause

The Clerk session object was being returned from SvelteKit load functions. The Clerk session object contains non-serializable functions (like `getToken()`) that cannot be serialized by SvelteKit's data transport system.

## Solution

Modified both load functions to return only serializable data:

**Before:**

```typescript
return {
  session, // ❌ Contains non-serializable functions
  userId,
  user,
};
```

**After:**

```typescript
return {
  hasSession: !!session, // ✅ Serializable boolean
  userId,
  user,
};
```

## Files Modified

1. `frontend/src/routes/+layout.server.ts` - Changed return statement to use `hasSession` instead of `session`
2. `frontend/src/routes/(auth)/+layout.server.ts` - Changed return statement to use `hasSession` instead of `session`

## Authentication Flow Preserved

- ✅ Server-side JWT token verification still works
- ✅ Backend user fetching still works
- ✅ User data is still available in pages
- ✅ Clerk session can be accessed client-side via clerkProvider
- ✅ Session object remains available in `event.locals` for server-side operations

## Testing

To verify the fix:

1. Restart the application (stop and start services)
2. Navigate to http://localhost:5173
3. Verify application loads without serialization errors
4. Confirm authentication still works (user is logged in)

## Related Issues

- Related to Clerk authentication integration (BP-5qy6)
- Related to Clerk environment variables fixes (BP-rpwv, BP-pkey)
- Related to JWKS verification implementation (BP-jwks)

## Impact

**Critical:** This fix was blocking the entire application from loading. Without this fix, users could not access any pages due to the serialization error.

## Notes

The authentication flow was actually working correctly on the backend (JWT token verified successfully, user found in database). The issue was purely on the frontend serialization layer where SvelteKit could not serialize the Clerk session object for transport between server and client.
