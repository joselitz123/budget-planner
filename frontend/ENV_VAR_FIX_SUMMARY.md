# Clerk Environment Variable Loading Fix - Summary

## Bug ID: BP-rpwv

**Type:** Critical Bug  
**Status:** Fixed  
**Impact:** Was blocking all authentication functionality

---

## Root Cause Analysis

### Primary Issue: Incorrect `envDir` Configuration in Vite

The [`vite.config.ts`](frontend/vite.config.ts) file had the following configuration:

```typescript
export default defineConfig({
	envDir: './',  // ❌ INCORRECT
	plugins: [...]
});
```

**Problem:** This told Vite to look for `.env` files in the **current working directory** (`/workspace/budget-planner`), but the actual `.env` file was located in the **`frontend/` subdirectory**.

**Result:** Vite could not find the `.env` file, so `PUBLIC_CLERK_PUBLISHABLE_KEY` was undefined in client-side code.

### Secondary Issues Found

1. **Server-side env var access pattern**: Server-side code in [`+layout.server.ts`](frontend/src/routes/+layout.server.ts:28) was using `process.env.PUBLIC_API_URL`, which is incorrect. Server-side code should not use `PUBLIC_` prefixed env vars.

2. **Static import causing SSR issues**: [`profile/+page.svelte`](frontend/src/routes/profile/+page.svelte:6) was using a static import for Clerk SDK instead of dynamic import, which could cause SSR issues.

3. **Test setup using old prefix**: [`src/test/setup.ts`](frontend/src/test/setup.ts:22) was still using the old `VITE_PUBLIC_` prefix instead of `PUBLIC_`.

---

## Fixes Implemented

### 1. Removed `envDir` from Vite Configuration

**File:** [`frontend/vite.config.ts`](frontend/vite.config.ts)

**Change:** Removed the `envDir: './'` line, allowing Vite to use its default behavior of looking for `.env` files in the project root (where `vite.config.ts` is located).

```diff
export default defineConfig({
-	envDir: './',
	plugins: [
```

### 2. Added Server-Side API URL Environment Variable

**Files:**

- [`frontend/.env`](frontend/.env)
- [`frontend/.env.example`](frontend/.env.example)

**Change:** Added `API_URL` (without `PUBLIC_` prefix) for server-side use:

```diff
# API Configuration
PUBLIC_API_URL=http://localhost:8080/api
+API_URL=http://localhost:8080/api
```

### 3. Updated Server-Side Code to Use Correct Env Var

**File:** [`frontend/src/routes/+layout.server.ts`](frontend/src/routes/+layout.server.ts)

**Change:** Changed from `process.env.PUBLIC_API_URL` to `process.env.API_URL`:

```diff
- const response = await event.fetch(`${process.env.PUBLIC_API_URL}/users/me`, {
+ const response = await event.fetch(`${process.env.API_URL}/users/me`, {
```

```diff
- const onboardingResponse = await event.fetch(`${process.env.PUBLIC_API_URL}/auth/onboarding`, {
+ const onboardingResponse = await event.fetch(`${process.env.API_URL}/auth/onboarding`, {
```

### 4. Fixed Clerk SDK Import Pattern in Profile Page

**File:** [`frontend/src/routes/profile/+page.svelte`](frontend/src/routes/profile/+page.svelte)

**Change:** Changed from static import to dynamic import and added logging:

```diff
- import clerkPkg from '@clerk/clerk-js';
- const { Clerk } = clerkPkg;

  onMount(async () => {
  	try {
+ 		// Dynamic import to prevent SSR issues
+ 		const clerkPkg = await import('@clerk/clerk-js');
+ 		const { Clerk } = clerkPkg;
+
  		const publishableKey = import.meta.env.PUBLIC_CLERK_PUBLISHABLE_KEY;
+
+ 		console.log('[Profile] Clerk env var check:', {
+ 			hasKey: !!publishableKey,
+ 			keyPrefix: publishableKey ? publishableKey.substring(0, 10) + '...' : 'N/A'
+ 		});
```

### 5. Added Comprehensive Logging for Debugging

**Files:**

- [`frontend/src/lib/auth/clerkProvider.ts`](frontend/src/lib/auth/clerkProvider.ts)
- [`frontend/src/routes/+layout.svelte`](frontend/src/routes/+layout.svelte)
- [`frontend/src/routes/profile/+page.svelte`](frontend/src/routes/profile/+page.svelte)

**Change:** Added detailed console logging to track env var loading and Clerk initialization:

```typescript
console.log("[Clerk] Environment check:", {
  hasKey: !!publishableKey,
  keyPrefix: publishableKey ? publishableKey.substring(0, 10) + "..." : "N/A",
  envMode: import.meta.env.MODE,
  allEnvVars: Object.keys(import.meta.env).filter((k) =>
    k.startsWith("PUBLIC_")
  ),
});
```

### 6. Fixed Test Setup Environment Variables

**File:** [`frontend/src/test/setup.ts`](frontend/src/test/setup.ts)

**Change:** Updated from `VITE_PUBLIC_` to `PUBLIC_` prefix:

```diff
vi.stubGlobal('import.meta', {
	env: {
-		VITE_PUBLIC_API_URL: 'http://localhost:8080/api',
-		VITE_PUBLIC_CLERK_PUBLISHABLE_KEY: 'pk_test_mock_key'
+		PUBLIC_API_URL: 'http://localhost:8080/api',
+		PUBLIC_CLERK_PUBLISHABLE_KEY: 'pk_test_mock_key'
	}
});
```

---

## Environment Variable Best Practices

### Client-Side (Browser)

- **Prefix:** `PUBLIC_` (required by SvelteKit/Vite)
- **Access:** `import.meta.env.PUBLIC_*`
- **Example:** `import.meta.env.PUBLIC_CLERK_PUBLISHABLE_KEY`
- **Usage:** Safe to expose to client (non-sensitive data)

### Server-Side (Node.js)

- **Prefix:** No prefix required (or `PRIVATE_` for clarity)
- **Access:** `process.env.*` or SvelteKit's `$env/static/private`
- **Example:** `process.env.CLERK_SECRET_KEY` or `process.env.API_URL`
- **Usage:** Never expose `PUBLIC_` prefixed vars to server code

### Environment Variable Hierarchy (Vite)

1. `.env.local` - Highest priority (gitignored)
2. `.env.[mode].local` - Mode-specific local (gitignored)
3. `.env.[mode]` - Mode-specific (committed)
4. `.env` - Default (committed)

---

## Verification Steps

After applying the fixes:

1. **Restart the dev server:**

   ```bash
   cd frontend
   npm run dev
   ```

2. **Check browser console for logs:**
   - Look for `[Clerk] Environment check:` logs
   - Verify `hasKey: true` and correct key prefix
   - Look for `[Clerk] Successfully initialized` or `[App] Clerk successfully initialized`

3. **Test authentication flow:**
   - Navigate to `/sign-in` - Should load Clerk sign-in form
   - Navigate to `/sign-up` - Should load Clerk sign-up form
   - After signing in, should be redirected to `/`
   - Profile page should load correctly

4. **Verify env vars are loaded:**
   - Open browser DevTools Console
   - Type: `import.meta.env.PUBLIC_CLERK_PUBLISHABLE_KEY`
   - Should return the actual key value (not undefined)

---

## Files Modified

1. [`frontend/vite.config.ts`](frontend/vite.config.ts) - Removed `envDir` configuration
2. [`frontend/.env`](frontend/.env) - Added `API_URL` for server-side use
3. [`frontend/.env.example`](frontend/.env.example) - Added `API_URL` example
4. [`frontend/src/routes/+layout.server.ts`](frontend/src/routes/+layout.server.ts) - Updated to use `API_URL`
5. [`frontend/src/routes/+layout.svelte`](frontend/src/routes/+layout.svelte) - Added logging
6. [`frontend/src/routes/profile/+page.svelte`](frontend/src/routes/profile/+page.svelte) - Fixed import pattern, added logging
7. [`frontend/src/lib/auth/clerkProvider.ts`](frontend/src/lib/auth/clerkProvider.ts) - Added logging
8. [`frontend/src/test/setup.ts`](frontend/src/test/setup.ts) - Fixed env var prefix

---

## Recommendations

1. **Environment Variable Documentation:** Consider creating a dedicated `ENV_VARS.md` file documenting all environment variables, their purposes, and where they should be used.

2. **Type Safety:** Add TypeScript type definitions for environment variables:

   ```typescript
   // frontend/src/lib/env.d.ts
   interface ImportMetaEnv {
     readonly PUBLIC_CLERK_PUBLISHABLE_KEY: string;
     readonly PUBLIC_API_URL: string;
     readonly PUBLIC_DEV_AUTH_PROVIDER: string;
     readonly PUBLIC_APP_NAME: string;
     readonly PUBLIC_APP_SHORT_NAME: string;
   }

   interface ImportMeta {
     readonly env: ImportMetaEnv;
   }
   ```

3. **Validation:** Add validation at startup to ensure required env vars are present:

   ```typescript
   const requiredEnvVars = ["PUBLIC_CLERK_PUBLISHABLE_KEY", "PUBLIC_API_URL"];
   const missing = requiredEnvVars.filter((key) => !import.meta.env[key]);
   if (missing.length > 0) {
     console.error("Missing required environment variables:", missing);
   }
   ```

4. **Consistent Naming:** Ensure all env vars follow consistent naming conventions across the project.

---

## Conclusion

The root cause was a simple misconfiguration in [`vite.config.ts`](frontend/vite.config.ts) that prevented Vite from finding the `.env` file. By removing the `envDir` configuration and fixing related issues with server-side env var access and import patterns, the Clerk environment variable loading issue has been resolved.

The authentication flow should now work correctly, and the added logging will help diagnose any future issues with environment variable loading.
