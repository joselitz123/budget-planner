---
# BP-x0f1
title: Install and configure beans CLI tool
status: completed
type: bug
created_at: 2026-01-02T13:24:01Z
updated_at: 2026-01-02T13:24:01Z
---

## Problem

The beans CLI tool described in CLAUDE.md was not available in the system. When trying to run `beans` commands, the system reported that the command was not found.

## Root Cause

The beans CLI tool (github.com/hmans/beans) was not installed on the system. The project documentation references beans for task management, but the CLI tool itself was never installed.

## Solution

1. Installed beans CLI via Go:
   ```bash
   export GOPATH=$HOME/go
   mkdir -p $GOPATH/bin
   go install github.com/hmans/beans@latest
   ```

2. Added to PATH permanently:
   ```bash
   echo 'export PATH=$HOME/go/bin:$PATH' >> ~/.bashrc
   echo 'export GOPATH=$HOME/go' >> ~/.bashrc
   ```

3. Verified installation:
   - `beans version` shows: beans dev (unknown) built unknown
   - `beans list` successfully displays all beans
   - `beans prime` outputs agent instructions

## Files Modified

- ~/.bashrc - Added PATH and GOPATH configuration

## Verification

✅ beans CLI is now installed and functional
✅ All beans commands work (list, create, update, prime, etc.)
✅ Project can now use beans for task management as documented in CLAUDE.md

## Notes

The beans tool is installed at: /home/vscode/go/bin/beans
Version installed: dev (latest from main branch)
