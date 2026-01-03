---
# BP-gc60
title: Create Budget button appears disabled when Total Limit field is filled
status: completed
type: bug
priority: high
created_at: 2026-01-03T04:06:49Z
updated_at: 2026-01-03T04:17:33Z
---

When navigating to the Budget Overview page and creating a budget, a modal appears. After filling in the Total Limit field and clicking the "Create Budget" button, the form doesn't submit - the button appears to be disabled.

**Location:** CreateBudgetModal component at frontend/src/lib/components/budget/CreateBudgetModal.svelte

**Priority:** High (blocks core functionality)

**Type:** Bug

**Context:** This is a critical UX issue preventing users from creating budgets in the application.