---
# BP-480r
title: Budget Sharing Feature
status: open
type: feature
priority: normal
created_at: 2025-12-28T16:51:45Z
updated_at: 2025-12-28T16:51:45Z
---

Implement complete budget sharing workflow with invitations and permissions.

## Background
Backend has full sharing API (invitations, access control, permissions) but frontend UI is not implemented.

## Acceptance Criteria
- [ ] Share button in budget overview
- [ ] Share invitation dialog (email input, permission dropdown)
- [ ] Invitations list (pending/accepted/declined)
- [ ] Shared budgets view (see budgets shared with me)
- [ ] Permission-based UI (view vs edit access)
- [ ] Accept/decline invitation flow

## Technical Details

### UI Components to Create
- ShareBudgetDialog.svelte (enter email, select permission)
- InvitationList.svelte (show pending invitations)
- SharedBudgetsPage.svelte (/shared route)

### User Flow
1. Owner clicks "Share Budget"
2. Enters recipient email + selects permission (view/edit)
3. POST /api/shares creates invitation
4. Recipient gets email (or share link)
5. Recipient accepts via PUT /api/shares/invitations/:id
6. Budget appears in recipient's "Shared with Me" view

### Permission Handling
- View: Can see budget, transactions, analytics
- Edit: Can add/edit transactions, not budget settings
- Owner: Full control including sharing

### Files to Create
- frontend/src/lib/components/sharing/ShareBudgetDialog.svelte
- frontend/src/lib/components/sharing/InvitationList.svelte
- frontend/src/routes/shared/+page.svelte

### Files to Modify
- frontend/src/routes/+page.svelte (add Share button)
- frontend/src/lib/stores/budgets.ts (add shared budgets)
- frontend/src/lib/api/shares.ts (create if not exists)

## Effort Estimate
2-3 hours

## Dependencies
- BP-r94p (Backend API Integration) - completed

## Session Date
2025-12-28