---
# BP-p808
title: Backend port 8080 already in use when starting
status: completed
type: bug
priority: critical
tags:
  - backend
  - startup
  - port-conflict
created_at: 2026-01-02T04:36:00Z
updated_at: 2026-01-02T04:36:00Z
---

## Problem

Backend fails to start because port 8080 is already in use when running the `start.sh` script.

### Error Message

```
2026/01/02 04:32:44 Server starting on :8080
2026/01/02 04:32:44 Server error: listen tcp :8080: bind: address already in use
exit status 1
```

## Context

This issue occurs when running the `start.sh` script to start the application. The previous backend process may not have been killed properly, leaving port 8080 occupied. The `start.sh` script needs to check for and kill any existing process on port 8080 before attempting to start a new backend instance.

## Reproduction Steps

1. Run `./start.sh` to start the application
2. Stop the application (e.g., with Ctrl+C)
3. Run `./start.sh` again without manually killing the backend process
4. Observe the error: "listen tcp :8080: bind: address already in use"

## Expected Behavior

The `start.sh` script should:

1. Check if port 8080 is already in use
2. If a process is using port 8080, kill it gracefully
3. Start the new backend instance successfully

## Actual Behavior

The `start.sh` script attempts to start the backend without checking for existing processes, resulting in a port binding error and startup failure.

## Root Cause

The `start.sh` script does not include logic to:

- Detect if port 8080 is already in use
- Kill existing processes on port 8080 before starting a new backend instance

This is a common issue when developing locally where processes may not be properly terminated.

## Proposed Solution

Update the `start.sh` script to add port checking and cleanup logic before starting the backend:

1. Check if port 8080 is in use using `lsof -i :8080` or `netstat -tlnp | grep 8080`
2. If a process is found, extract the PID and kill it
3. Wait a moment to ensure the port is released
4. Proceed with starting the backend

Example bash logic:

```bash
# Kill existing backend process on port 8080
if lsof -i :8080 > /dev/null 2>&1; then
    echo "Killing existing process on port 8080..."
    lsof -ti :8080 | xargs kill -9 2>/dev/null || true
    sleep 1
fi
```

## Files Affected

- `start.sh` - Add port checking and cleanup logic before starting backend

## Related Issues

- None

## Impact

**Critical** - Prevents the application from starting when a previous backend process is still running. This blocks development and testing workflows.

## Effort Estimate

30 minutes

## Session Date

2026-01-02

## Fix Applied

**Status:** Completed on 2026-01-02

### Files Modified

- [`start.sh`](start.sh)

### Fix Description

Added port 8080 check and kill logic before starting backend. The `start.sh` script now:

1. Checks if port 8080 is already in use using `lsof -i :8080`
2. If a process is found, extracts the PID and kills it with `kill -9`
3. Waits 1 second to ensure the port is released
4. Proceeds with starting the backend

This prevents the "address already in use" error when restarting the application.

### Implementation

```bash
# Kill existing backend process on port 8080
if lsof -i :8080 > /dev/null 2>&1; then
    echo "Killing existing process on port 8080..."
    lsof -ti :8080 | xargs kill -9 2>/dev/null || true
    sleep 1
fi
```
