---
# BP-olz1
title: Offline Sync Testing
status: todo
type: task
priority: normal
created_at: 2025-12-28T17:11:03Z
updated_at: 2025-12-28T17:11:03Z
---

Plan and execute comprehensive testing of the offline sync functionality to ensure it works correctly in all scenarios.

## Background
The offline sync implementation (BP-iyr6) is complete but needs thorough testing to verify:
- Background sync queue processing works correctly
- Exponential backoff retry logic functions as expected
- Conflict resolution handles edge cases properly
- Sync status indicators display accurately
- Manual sync button works correctly
- Offline→online flow transitions smoothly

## Test Scenarios to Cover

### 1. Basic Sync Flow
- [ ] Create transaction while online → verify sync happens
- [ ] Create transaction while offline → go online → verify sync happens
- [ ] Multiple pending operations → verify FIFO processing

### 2. Retry Logic Testing
- [ ] Simulate sync failure → verify exponential backoff works
- [ ] Verify retry attempts increment correctly (up to 5)
- [ ] Verify max retry behavior → permanent failure with toast
- [ ] Test jitter to prevent thundering herd

### 3. Conflict Resolution Testing
- [ ] Same record edited locally and on server → verify timestamp-based resolution
- [ ] Test server wins when timestamps are equal
- [ ] Test local wins when local timestamp is newer
- [ ] Verify IndexedDB updates correctly after conflict resolution

### 4. Pull from Server Testing
- [ ] Server has new data → verify local IndexedDB updates
- [ ] Server has updated data → verify conflict resolution applies
- [ ] Server has deleted data → verify local deletion handling

### 5. UI/UX Testing
- [ ] Verify sync status indicator shows correct states (syncing/synced/error/pending/offline)
- [ ] Verify pending count updates in real-time
- [ ] Verify "Sync Now" button enables/disables correctly
- [ ] Verify toast notifications appear for all sync events
- [ ] Verify last sync timestamp updates correctly

### 6. Background Sync Testing
- [ ] Verify background sync runs every 30 seconds
- [ ] Verify background sync stops when going offline
- [ ] Verify background sync resumes when coming back online
- [ ] Verify background sync stops on page unload

### 7. Edge Cases
- [ ] Rapid network changes (online/offline/online) → verify stability
- [ ] Large number of pending operations (100+) → verify performance
- [ ] Network timeout during sync → verify retry logic
- [ ] Server returns 500 error → verify retry behavior
- [ ] Server returns 401 error → verify auth handling

### 8. Integration with Backend
- [ ] Test with actual backend /api/sync/push endpoint
- [ ] Test with actual backend /api/sync/pull endpoint
- [ ] Verify request/response format matches backend expectations
- [ ] Test with JWT authentication tokens

## Test Plan

### Phase 1: Unit Testing (without backend)
- Mock fetch API to simulate server responses
- Test sync queue processing logic
- Test exponential backoff calculation
- Test conflict resolution function
- Test status indicator state changes

### Phase 2: Integration Testing (with mock backend)
- Set up local mock server with sync endpoints
- Test full sync flow end-to-end
- Test error scenarios (network failures, server errors)
- Test retry logic with actual delays

### Phase 3: Manual Testing (with real backend)
- Start backend server
- Start frontend dev server
- Use DevTools Network throttling to simulate offline
- Test all scenarios from "Test Scenarios" list
- Document any bugs found

### Phase 4: Edge Case Testing
- Test with large datasets
- Test with rapid state changes
- Test with browser tab background/foreground
- Test with multiple tabs open

## Test Data Setup
- Create test user account
- Prepare test budgets, categories, transactions
- Set up different conflict scenarios
- Create scripts for automated data generation

## Success Criteria
- All test scenarios pass
- No console errors during normal operation
- Sync status accurately reflects reality
- Conflicts resolve correctly
- UI updates smoothly without jank
- No memory leaks in background sync

## Dependencies
- BP-iyr6 (Offline Sync Implementation) - completed
- Backend sync endpoints must be functional
- Test backend server must be running

## Effort Estimate
2-3 hours

## Session Date
2025-12-28