---
# BP-bind-value
title: Fix Input and Textarea components two-way binding
status: completed
type: bug
priority: critical
tags:
  - frontend
  - ui
  - input
  - textarea
  - bug
  - completed
created_at: 2026-01-02T14:25:00Z
updated_at: 2026-01-02T14:25:00Z
completed_at: 2026-01-02T14:25:00Z
---

## Problem

The Input and Textarea components were using one-way binding (`{value}`) instead of two-way binding (`bind:value`), which caused the Create Budget Modal to fail:

1. **Total Limit field showing "Total limit is required" error** - Even when a value was entered
2. **Create Budget button doing nothing** - Button remained disabled

### Error Details

**Files:**

- [`frontend/src/lib/components/ui/input/input.svelte`](frontend/src/lib/components/ui/input/input.svelte:1)
- [`frontend/src/lib/components/ui/textarea/textarea.svelte`](frontend/src/lib/components/ui/textarea/textarea.svelte:1)

**Previous (incorrect) implementation:**

```svelte
<script lang="ts">
  export let className = '';
  export let value: string = '';
</script>

<input
  {value}
  class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 {className}"
  {...$$restProps}
/>
```

### Root Cause

The components were using `{value}` (one-way binding) instead of `bind:value` (two-way binding).

When a user typed in the Total Limit field:

- The native input element updated its own value
- But the parent's `totalLimit` variable was NEVER updated
- Validation failed because `totalLimit` remained empty
- The submit button stayed disabled

This happened because:

- `{value}` only sets the initial value from parent to child (one-way)
- `bind:value` creates a two-way binding that updates the parent when the child changes
- Without `bind:value`, changes to the input don't propagate back to the parent

### Impact

**Critical** - Users could not create budgets because form values were never captured:

- Form submission was completely broken
- Input fields appeared to work but values were not captured
- Validation always failed even with valid input
- This blocked the entire budget creation workflow

## Solution

Changed from one-way binding to two-way binding in both Input and Textarea components.

### Implementation Details

**Files Modified:**

1. [`frontend/src/lib/components/ui/input/input.svelte`](frontend/src/lib/components/ui/input/input.svelte:1)
2. [`frontend/src/lib/components/ui/textarea/textarea.svelte`](frontend/src/lib/components/ui/textarea/textarea.svelte:1)

**Changes Made:**

Changed from `{value}` to `bind:value` on the native elements:

**Before (incorrect):**

```svelte
<input
  {value}
  class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 {className}"
  {...$$restProps}
/>
```

**After (correct):**

```svelte
<input
  bind:value
  class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 {className}"
  {...$$restProps}
/>
```

Same change applied to Textarea component.

### Why This Works

- The `bind:value` directive on the native element ensures two-way binding
- When parent uses `bind:value={totalLimit}`, the component's `bind:value` on the native input element ensures changes flow both ways
- User typing in input → native input updates → parent's `totalLimit` updates → validation passes → submit button enables
- This is the correct way to implement two-way binding in Svelte wrapper components

## Files Modified

- [`frontend/src/lib/components/ui/input/input.svelte`](frontend/src/lib/components/ui/input/input.svelte:1)

  - Changed `{value}` to `bind:value` on input element
  - Maintained all other props and styling

- [`frontend/src/lib/components/ui/textarea/textarea.svelte`](frontend/src/lib/components/ui/textarea/textarea.svelte:1)
  - Changed `{value}` to `bind:value` on textarea element
  - Maintained all other props and styling

## Testing

### Verification Steps

1. ✅ Input component uses two-way binding
2. ✅ Textarea component uses two-way binding
3. ✅ Form submission works when using Input with `bind:value`
4. ✅ Validation correctly reads input values
5. ✅ TypeScript compilation successful
6. ✅ No breaking changes to existing Input/Textarea usage

### Test Results

- CreateBudgetModal form submission now works correctly
- Total Limit input properly captures and validates user input
- Form validation triggers appropriately
- Submit button enables when form is valid
- All forms using Input/Textarea components now work correctly

## Impact

**Critical** - This fix restores full functionality to all forms using Input and Textarea components:

- Create Budget Modal now works
- Add Expense Modal (uses Input for amount, date, description, dueDate)
- Payment Method Form (uses Input for name, lastFour, brand, creditLimit, currentBalance)
- Share Budget Dialog (uses Input for email)
- Textarea components (used for notes in AddExpenseModal)

All forms using these components will now work correctly with two-way binding.

## Related Issues

- **BP-input-value** - Previous attempt to fix Input component by adding `value` prop and `{value}`. This was incomplete because it only enabled one-way binding, not two-way binding.

- **BP-btn-type** - Similar fix for Button component's `type` attribute forwarding. All three issues stem from same root cause: wrapper components not properly implementing Svelte's binding directives.

## Effort Estimate

15 minutes

## Session Date

2026-01-02

## Technical Details

### Svelte Binding Directives

Svelte provides several binding directives:

- `bind:value` - Two-way binding for input/textarea/select elements
- `bind:checked` - Two-way binding for checkbox/radio elements
- `bind:group` - Two-way binding for radio button groups
- `bind:this` - Bind component to DOM element reference

When creating wrapper components, you must:

1. Declare the prop (e.g., `export let value: string = ''`)
2. Use the same binding directive on the native element (e.g., `bind:value`)
3. This creates a "binding chain" that propagates changes both ways

### Binding Chain Example

```svelte
<!-- Parent component -->
<script>
  let totalLimit = '';
</script>

<Input bind:value={totalLimit} />

<!-- Input component -->
<script>
  export let value: string = '';
</script>

<input bind:value />

<!-- Native input element -->
```

When user types:

1. Native input updates its value
2. `bind:value` on native input updates Input component's `value` prop
3. `bind:value` on Input component updates parent's `totalLimit` variable
4. Parent's reactive statements run with new value

### Type Safety

Both components maintain type safety:

```typescript
export let value: string = "";
```

This ensures:

- Only string values can be passed
- TypeScript provides autocomplete and validation
- Default empty string maintains backward compatibility
- Clear documentation of accepted value type

### Backward Compatibility

The fix maintains full backward compatibility:

- Existing Input/Textarea usage without `bind:value` still works
- Default empty string is HTML5 default behavior
- No changes required to existing Input/Textarea component consumers
- $$restProps still spreads any additional props correctly

### Integration with Shadcn-Svelte

This fix aligns with Shadcn-Svelte patterns:

- Proper use of Svelte binding directives
- Type-safe prop definitions
- Backward-compatible defaults
- Clean separation of concerns

## Notes

This is the correct implementation of two-way binding in Svelte wrapper components. The key insight is that when a parent uses `bind:value={variable}`, the wrapper component must also use `bind:value` on the native element to complete the binding chain.

The previous fix (BP-input-value) was incomplete because it only added the `value` prop and used `{value}` (one-way binding). This is a common mistake when implementing wrapper components.

## Browser Caching

After this fix, users should perform a hard refresh (Ctrl+Shift+R or Cmd+Shift+R) to ensure the updated JavaScript bundle is loaded, as browsers may cache the old bundle.

## Similar Issues

This is third critical attribute/binding bug discovered in this project:

1. **BP-btn-type** - Button component not forwarding `type` attribute
2. **BP-input-value** - Input component using one-way binding instead of two-way binding (incomplete fix)
3. **BP-bind-value** - Input/Textarea components using one-way binding instead of two-way binding (correct fix)

All issues stem from improper implementation of Svelte's binding and attribute forwarding in wrapper components.
