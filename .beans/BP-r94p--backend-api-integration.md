---
# BP-r94p
title: Backend API Integration
status: todo
type: feature
priority: critical
tags:
    - frontend
    - api
    - sync
    - backend
created_at: 2025-12-28T05:36:47Z
updated_at: 2025-12-28T05:36:47Z
---

Wire up actual API calls to Go backend for all CRUD operations.

## Acceptance Criteria
- [ ] Budgets API integration (GET/POST /api/budgets)
- [ ] Transactions API integration (GET/POST/PUT /api/transactions)
- [ ] Categories API integration (GET /api/categories)
- [ ] Sync API integration (POST /api/sync/push, POST /api/sync/pull)
- [ ] Error handling with user-friendly messages
- [ ] Loading states during API calls
- [ ] JWT token integration from Clerk

## Files to Modify
- frontend/src/lib/api/budgets.ts (create)
- frontend/src/lib/api/transactions.ts (create)
- frontend/src/lib/api/sync.ts (create)
- frontend/src/lib/api/client.ts (extend)

## Effort Estimate
3 hours

## Migration Notes
Migrated from frontend/todo.md Priority 1, Task #1