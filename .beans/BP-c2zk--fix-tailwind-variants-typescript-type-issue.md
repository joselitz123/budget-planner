---
# BP-c2zk
title: Fix tailwind-variants TypeScript Type Issue
status: todo
type: bug
priority: low
tags:
    - frontend
    - tech-debt
    - typescript
created_at: 2025-12-28T15:11:09Z
updated_at: 2025-12-28T15:11:09Z
---

tailwind-variants tv() function doesn't properly type the `class` parameter, causing TypeScript errors when trying to pass className to variants.

## Problem
When using tailwind-variants' tv() function, passing the `class` parameter causes TypeScript errors because the type definitions don't include it in the VariantProps.

## Workaround Used
String concatenation instead of passing `class` to tv():
```typescript
spinnerVariants({ size, color }) + (className ? ` ${className}` : '')
```

## Files Affected
- frontend/src/lib/components/ui/spinner/spinner.svelte
- frontend/src/lib/components/ui/skeleton/skeleton.svelte

## Proper Solutions
1. Create cn() utility function (like shadcn/ui) for className merging
2. Create proper TypeScript type definitions for tailwind-variants
3. Use clsx or classnames library for className merging

## Impact
Minor code quality issue - workaround functions correctly and components work as expected.

## Priority
Low - no functional impact, workaround is clean and maintainable

## Session Date
2025-12-28

## Discovered During
BP-r94p implementation (loading states work) - encountered when creating Spinner and Skeleton components.