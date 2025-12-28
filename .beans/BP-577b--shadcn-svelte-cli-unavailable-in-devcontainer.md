---
# BP-577b
title: Shadcn-Svelte CLI Unavailable in Devcontainer
status: todo
type: bug
priority: normal
tags:
    - frontend
    - tech-debt
created_at: 2025-12-28T05:37:49Z
updated_at: 2025-12-28T05:37:49Z
---

## Problem
Shadcn-Svelte CLI requires interactive TTY which is unavailable in devcontainer. All UI components were created manually.

## Workaround
- Created components.json manually
- Manually created all UI components
- Installed dependencies with --legacy-peer-deps

## Acceptance Criteria
- [ ] Set up interactive terminal in devcontainer OR accept manual approach
- [ ] Document manual component creation process
- [ ] Create component templates for future use

## Impact
Medium - Development is slower but not blocked. Manual approach is sustainable.

## Effort Estimate
2 hours

## Type: Technical Debt
## Session Date
2025-12-27

## Files Modified
- frontend/components.json
- frontend/src/lib/components/ui/*

## Migration Notes
Migrated from frontend/todo.md Technical Debts section