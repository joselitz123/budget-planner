---
# BP-x4sk
title: Change currency from USD to Philippine Peso
status: completed
type: task
priority: normal
tags:
    - frontend
    - currency
    - ui
created_at: 2025-12-28T06:06:54Z
updated_at: 2025-12-28T06:10:00Z
---

Update all currency displays from US Dollar ($) to Philippine Peso (₱).

## Changes Required:
- [x] Update formatCurrency function in frontend/src/lib/utils/format.ts
  - Change locale from 'en-US' to 'en-PH'
  - Change currency from 'USD' to 'PHP'
  - Update comment
- [x] Update backend test data in backend/internal/handlers/test_setup.go
  - Change test user currency from 'USD' to 'PHP' for consistency

## Files Modified:
- frontend/src/lib/utils/format.ts (lines 2, 5, 7)
- backend/internal/handlers/test_setup.go (line 83)

## Expected Outcome:
All currency displays will show ₱ symbol instead of $ across:
- Budget Overview page
- Transactions page
- Bills page
- Any other component using formatCurrency()

## Session Date: 2025-12-28