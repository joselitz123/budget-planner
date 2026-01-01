---
# BP-cors
title: Fix CORS configuration bug - getEnvSlice() not splitting comma-separated origins
status: completed
type: bug
priority: normal
tags:
  - backend
  - cors
  - config
  - normal
created_at: 2026-01-01T14:47:00Z
updated_at: 2026-01-01T14:47:00Z
completed_at: 2026-01-01T14:47:00Z
---

## Problem

The `getEnvSlice()` function in `backend/internal/config/config.go` was returning the entire comma-separated string as a single element instead of properly splitting it into individual origins. This caused CORS configuration to fail when multiple origins were specified in the environment variable.

### Root Cause

The `getEnvSlice()` function was not using `strings.Split()` to parse comma-separated values from the environment variable. Instead, it was treating the entire comma-separated string as a single origin.

Example:

- Environment variable: `CORS_ORIGINS=http://localhost:5173,http://localhost:3000`
- Expected result: `["http://localhost:5173", "http://localhost:3000"]`
- Actual result: `["http://localhost:5173,http://localhost:3000"]` (single element)

### Impact

**Normal Priority** - CORS would only work if a single origin was specified. Multiple origins configuration was broken, which could prevent legitimate frontend requests from different origins during development.

## Solution

Updated the `getEnvSlice()` function to use `strings.Split()` to properly parse comma-separated values from environment variables.

### Implementation Details

1. **Added `strings` import** to `backend/internal/config/config.go`:

   ```go
   import "strings"
   ```

2. **Updated `getEnvSlice()` function** at line 127:
   ```go
   func getEnvSlice(key string, defaultVal []string) []string {
       if value := os.Getenv(key); value != "" {
           return strings.Split(value, ",")
       }
       return defaultVal
   }
   ```

### Key Changes

**Before (incorrect):**

```go
func getEnvSlice(key string, defaultVal []string) []string {
    if value := os.Getenv(key); value != "" {
        return []string{value}  // Returns entire string as single element
    }
    return defaultVal
}
```

**After (correct):**

```go
func getEnvSlice(key string, defaultVal []string) []string {
    if value := os.Getenv(key); value != "" {
        return strings.Split(value, ",")  // Properly splits into multiple elements
    }
    return defaultVal
}
```

## Files Modified

1. `backend/internal/config/config.go` - Added `strings` import and fixed `getEnvSlice()` function

## Testing

### Verification Steps

1. ✅ CORS origins are now properly split from comma-separated environment variable
2. ✅ Multiple origins work correctly
3. ✅ OPTIONS requests are successful for all configured origins
4. ✅ Single origin configuration still works as expected
5. ✅ Default values are used when environment variable is not set

### Test Results

- Environment variable `CORS_ORIGINS=http://localhost:5173,http://localhost:3000` now correctly parses to two separate origins
- Both origins can successfully make requests to the backend
- CORS headers are properly set for each origin
- Pre-flight OPTIONS requests succeed for all configured origins

## Impact

**Normal Priority** - This fix ensures proper CORS configuration for development environments with multiple frontend origins.

- Developers can now configure multiple frontend origins (e.g., development and production)
- CORS configuration works as documented in `.env.example`
- No breaking changes to existing single-origin configurations

## Session Date

2026-01-01

## Effort Estimate

30 minutes

## Related Issues

- Discovered during testing of the JWT authentication fix
- Simple one-line fix with `strings.Split()`
- No impact on other parts of the codebase
