---
# BP-hhvi
title: Update CLAUDE.md files and clean up documentation
status: completed
type: task
priority: normal
tags:
    - documentation
    - tech-debt
    - cleanup
created_at: 2025-12-31T03:53:44Z
updated_at: 2025-12-31T04:00:00Z
---

## Overview

All CLAUDE.md files need to be reviewed and updated to reflect the current state of the application. References to deleted todo.md files should be removed, and the documentation should be consolidated to make it easy for Claude to understand the architecture.

## Problems to Fix

1. **Outdated references:** Multiple references to `backend/todo.md` and `frontend/todo.md` which no longer exist
2. **Status information:** Current status sections reference completed work that is now tracked in Beans
3. **Redundancy:** Some information is duplicated across files
4. **Missing context:** Key architectural decisions may not be well documented

## Files to Update

- `CLAUDE.md` (root) - Project overview, remove todo.md references from Key Reference Files section
- `backend/CLAUDE.md` - Backend guide, remove references to `todo.md.legacy` as primary source
- `frontend/CLAUDE.md` - Frontend guide, clean up status section

## Tasks

### Root CLAUDE.md
- [x] Remove `backend/todo.md` and `frontend/todo.md` from Key Reference Files section
- [x] Update Current Status section to point to Beans for task tracking
- [x] Ensure all essential architecture info is present
- [x] Verify WRAP-UP/KICK-START workflows are accurate

### Backend CLAUDE.md
- [x] Remove or minimize references to `todo.md.legacy`
- [x] Update Task Tracking section to emphasize Beans
- [x] Ensure all handler endpoints are documented
- [x] Verify sqlc patterns are clear

### Frontend CLAUDE.md
- [x] Remove or minimize references to `todo.md.legacy`
- [x] Update Implementation Status to reflect Beans tracking
- [x] Ensure all architectural patterns are documented
- [x] Verify component organization section

### Cleanup
- [x] Search for all remaining references to `todo.md` in the codebase
- [x] Update or remove those references
- [x] Consider if `todo.md.legacy` files should be removed or archived (kept for historical reference)

## Acceptance Criteria

- [x] All CLAUDE.md files are updated with current, accurate information
- [x] No broken references to deleted todo.md files
- [x] Architecture is clearly documented for Claude to understand
- [x] Beans is consistently referenced as the task tracking system

## Files Modified

- `/workspace/budget-planner/CLAUDE.md` - Updated Key Reference Files section
- `/workspace/budget-planner/backend/CLAUDE.md` - Updated Task Tracking and Important Files sections
- `/workspace/budget-planner/frontend/CLAUDE.md` - Updated Task Tracking, Important Files, and Project Structure sections

## Effort Estimate

1-2 hours

## Type

Documentation / Cleanup

## Session Date

2025-12-31
