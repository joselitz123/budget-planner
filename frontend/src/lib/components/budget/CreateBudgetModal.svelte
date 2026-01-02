<script lang="ts">
  import CustomModal from "$lib/components/ui/CustomModal.svelte";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { currentMonth } from "$lib/stores/budgets";
  import { currency } from "$lib/stores/settings";
  import { parseMonthKey } from "$lib/utils/format";
  import { createEventDispatcher } from "svelte";
  import { get } from "svelte/store";

  export let isOpen = false;
  export let isCreating = false;

  const dispatch = createEventDispatcher();

  // Form state
  let name = "";
  let totalLimit = "";

  // Form errors
  let errors = {
    name: "",
    totalLimit: "",
  };

  // Computed: Form is valid
  let formValid = false;

  // Watch form changes to update validity
  $: formValid = validateForm(false);

  // Get current month display
  $: monthDisplay = (() => {
    const monthValue = get(currentMonth);
    if (!monthValue) return "";
    const date = parseMonthKey(monthValue);
    return date.toLocaleDateString("en-US", { month: "long", year: "numeric" });
  })();

  // Get currency symbol
  $: currencySymbol = (() => {
    const currencyValue = get(currency);
    return currencyValue === "PHP"
      ? "₱"
      : currencyValue === "USD"
        ? "$"
        : currencyValue === "EUR"
          ? "€"
          : currencyValue === "GBP"
            ? "£"
            : "¥";
  })();

  /**
   * Validate form
   * @param updateErrors - Whether to update error messages
   * @returns true if form is valid
   */
  function validateForm(updateErrors: boolean = true): boolean {
    let valid = true;

    if (updateErrors) {
      errors = { name: "", totalLimit: "" };
    }

    // Validate total limit (required)
    if (!totalLimit || totalLimit.trim() === "") {
      if (updateErrors) {
        errors.totalLimit = "Total limit is required";
      }
      valid = false;
    } else {
      const limitValue = parseFloat(totalLimit);
      if (isNaN(limitValue) || limitValue <= 0) {
        if (updateErrors) {
          errors.totalLimit = "Total limit must be a positive number";
        }
        valid = false;
      }
    }

    return valid;
  }

  function closeModal() {
    isOpen = false;
    resetForm();
    dispatch("close");
  }

  function resetForm() {
    name = "";
    totalLimit = "";
    errors = { name: "", totalLimit: "" };
  }

  function handleSubmit() {
    if (!validateForm(true)) {
      return;
    }

    const limitValue = parseFloat(totalLimit);
    const budgetData = {
      name: name.trim() || undefined,
      totalLimit: limitValue,
    };

    dispatch("submit", budgetData);
    closeModal();
  }

  // Handle keyboard shortcuts
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit();
    }
  }
</script>

<CustomModal
  bind:open={isOpen}
  title="Create Budget"
  description="Set up a new budget for the selected month"
  className="max-w-md"
>
  <form
    on:submit|preventDefault={handleSubmit}
    on:keydown={handleKeydown}
    class="space-y-6"
  >
    <!-- Month Display Info Box -->
    <div
      class="p-4 bg-paper-light dark:bg-paper-dark/50 rounded-lg border-2 border-line-light dark:border-line-dark notebook-lines"
    >
      <div class="flex items-center space-x-2">
        <span class="material-icons-outlined text-primary dark:text-white"
          >calendar_today</span
        >
        <span class="text-sm text-gray-600 dark:text-gray-400 font-medium"
          >Budget Period:</span
        >
      </div>
      <p
        class="mt-2 text-lg font-display font-bold text-primary dark:text-white"
      >
        {monthDisplay}
      </p>
    </div>

    <!-- Budget Name (Optional) -->
    <div class="space-y-2">
      <Label
        for="name"
        class="text-base font-medium text-primary dark:text-white"
      >
        Budget Name
      </Label>
      <Input
        id="name"
        bind:value={name}
        placeholder="e.g., Monthly Budget"
        class="font-handwriting text-lg border-2 border-line-light dark:border-line-dark focus:border-accent-gold focus:ring-accent-gold/20"
        maxlength="100"
        aria-describedby="name-description"
      />
      <p
        id="name-description"
        class="text-xs text-gray-500 dark:text-gray-400 italic"
      >
        Optional: Give your budget a descriptive name
      </p>
      {#if errors.name}
        <p
          class="text-sm text-red-500 dark:text-red-400 font-medium"
          role="alert"
        >
          {errors.name}
        </p>
      {/if}
    </div>

    <!-- Total Limit (Required) -->
    <div class="space-y-2">
      <Label
        for="totalLimit"
        class="text-base font-medium text-primary dark:text-white"
      >
        Total Limit <span class="text-accent-gold">*</span>
      </Label>
      <div class="relative">
        <span
          class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-600 dark:text-gray-400 font-handwriting text-xl"
        >
          {currencySymbol}
        </span>
        <Input
          id="totalLimit"
          type="number"
          step="0.01"
          min="0"
          bind:value={totalLimit}
          placeholder="e.g., 5000"
          class="pl-10 font-handwriting text-xl border-2 border-line-light dark:border-line-dark focus:border-accent-gold focus:ring-accent-gold/20"
          aria-describedby="totalLimit-description"
          aria-invalid={errors.totalLimit ? "true" : "false"}
          aria-required="true"
        />
      </div>
      <p
        id="totalLimit-description"
        class="text-xs text-gray-500 dark:text-gray-400 italic"
      >
        Enter the total spending limit for this month
      </p>
      {#if errors.totalLimit}
        <p
          class="text-sm text-red-500 dark:text-red-400 font-medium"
          role="alert"
        >
          {errors.totalLimit}
        </p>
      {/if}
    </div>

    <!-- Actions -->
    <div
      class="flex justify-end space-x-3 pt-6 border-t-2 border-line-light dark:border-line-dark"
    >
      <Button
        type="button"
        variant="outline"
        onclick={closeModal}
        disabled={isCreating}
        class="font-medium"
      >
        Cancel
      </Button>
      <Button
        type="submit"
        disabled={!formValid || isCreating}
        class="min-w-[120px] font-medium"
      >
        {#if isCreating}
          <span class="flex items-center space-x-2">
            <span class="material-icons-outlined animate-spin text-sm"
              >refresh</span
            >
            <span>Creating...</span>
          </span>
        {:else}
          <span class="flex items-center space-x-2">
            <span class="material-icons-outlined text-sm">add_circle</span>
            <span>Create Budget</span>
          </span>
        {/if}
      </Button>
    </div>
  </form>
</CustomModal>
