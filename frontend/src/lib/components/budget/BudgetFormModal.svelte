<script lang="ts">
  import type { Budget } from "$lib/api/budgets";
  import CustomModal from "$lib/components/ui/CustomModal.svelte";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { currentMonth } from "$lib/stores/budgets";
  import { currency, getCurrencyLocale } from "$lib/stores/settings";
  import { getCurrencySymbol, parseMonthKey } from "$lib/utils/format";
  import { createEventDispatcher } from "svelte";
  import { get } from "svelte/store";

  export let isOpen = false;
  export let mode: "create" | "edit" = "create";
  export let budgetData: Budget | undefined = undefined;
  export let isSaving = false;
  export let currentSpending = 0;

  const dispatch = createEventDispatcher();

  // Form state
  let name = "";
  let totalLimit = "";
  let totalIncome = "";

  // Store original values for reset functionality (edit mode)
  let originalName = "";
  let originalTotalLimit = "";
  let originalTotalIncome = "";

  // Form errors
  let errors = {
    name: "",
    totalLimit: "",
    totalIncome: "",
  };

  // Computed: Form is valid
  $: formValid = isFormValid();

  // Computed: Month display
  $: monthDisplay = getMonthDisplay();

  // Get currency symbol using utility function
  $: currencySymbol = getCurrencySymbol(get(currency));

  // Computed: Savings calculation
  $: savings = calculateSavings(totalLimit, totalIncome);

  // Computed: Spending percentage (edit mode)
  $: spendingPercentage = calculateSpendingPercentage();

  // Computed: Budget health status (edit mode)
  $: budgetHealth = getBudgetHealth();

  // Computed: Warnings
  $: warnings = getWarnings();

  /**
   * Get current month display with locale-aware formatting
   */
  function getMonthDisplay(): string {
    const monthValue = get(currentMonth);

    if (!monthValue) {
      return "";
    }

    try {
      const date = parseMonthKey(monthValue);
      if (isNaN(date.getTime())) {
        console.warn(`Invalid month key format: ${monthValue}`);
        return "";
      }

      const locale = getCurrencyLocale(get(currency));
      return date.toLocaleDateString(locale, {
        month: "long",
        year: "numeric",
      });
    } catch (error) {
      console.error(
        `Error formatting month display for key "${monthValue}":`,
        error
      );
      return "";
    }
  }

  /**
   * Calculate savings (income - limit)
   */
  function calculateSavings(limit: string, income: string): number | null {
    if (!limit || !income) return null;
    const limitVal = Number(limit);
    const incomeVal = Number(income);
    if (isNaN(limitVal) || isNaN(incomeVal)) return null;
    return incomeVal - limitVal;
  }

  /**
   * Calculate spending percentage
   */
  function calculateSpendingPercentage(): number {
    if (!totalLimit || currentSpending === 0) return 0;
    const limitVal = Number(totalLimit);
    if (isNaN(limitVal) || limitVal === 0) return 0;
    return Math.min((currentSpending / limitVal) * 100, 100);
  }

  /**
   * Get budget health status
   */
  function getBudgetHealth(): "healthy" | "warning" | "danger" {
    if (!totalLimit) return "healthy";
    const limitVal = Number(totalLimit);
    if (isNaN(limitVal) || limitVal === 0) return "healthy";

    const percentage = (currentSpending / limitVal) * 100;
    if (percentage >= 100) return "danger";
    if (percentage >= 80) return "warning";
    return "healthy";
  }

  /**
   * Get warning messages
   */
  function getWarnings(): string[] {
    const warnings: string[] = [];

    if (mode === "edit" && totalLimit) {
      const limitVal = Number(totalLimit);
      if (!isNaN(limitVal) && currentSpending > limitVal) {
        warnings.push(
          `Current spending (${formatCurrency(currentSpending)}) exceeds new limit (${formatCurrency(limitVal)})`
        );
      }
    }

    if (totalLimit && totalIncome) {
      const limitVal = Number(totalLimit);
      const incomeVal = Number(totalIncome);
      if (
        !isNaN(limitVal) &&
        !isNaN(incomeVal) &&
        incomeVal < currentSpending
      ) {
        warnings.push(
          `Current spending (${formatCurrency(currentSpending)}) exceeds income (${formatCurrency(incomeVal)})`
        );
      }
    }

    if (savings !== null && savings < 0) {
      warnings.push(
        `Budget exceeds income by ${formatCurrency(Math.abs(savings))}`
      );
    }

    return warnings;
  }

  /**
   * Format currency value
   */
  function formatCurrency(value: number): string {
    const locale = getCurrencyLocale(get(currency));
    return new Intl.NumberFormat(locale, {
      style: "currency",
      currency: get(currency),
    }).format(value);
  }

  /**
   * Check if form is valid
   */
  function isFormValid(): boolean {
    const limitValid = totalLimit && Number(totalLimit) > 0;
    const incomeValid = !totalIncome || Number(totalIncome) >= 0;
    return Boolean(limitValid && incomeValid);
  }

  /**
   * Validate form and update error messages
   */
  function validateForm(): boolean {
    let valid = true;
    errors = { name: "", totalLimit: "", totalIncome: "" };

    // Validate total limit (required)
    const limitStr = String(totalLimit || "").trim();
    if (!limitStr) {
      errors.totalLimit = "Total limit is required";
      valid = false;
    } else {
      const limitValue = Number(limitStr);
      if (isNaN(limitValue) || limitValue <= 0) {
        errors.totalLimit = "Total limit must be a positive number";
        valid = false;
      }
    }

    // Validate total income (optional but must be non-negative if provided)
    const incomeStr = String(totalIncome || "").trim();
    if (incomeStr) {
      const incomeValue = Number(incomeStr);
      if (isNaN(incomeValue) || incomeValue < 0) {
        errors.totalIncome = "Total income cannot be negative";
        valid = false;
      }
    }

    return valid;
  }

  /**
   * Reset form to original values (edit mode)
   */
  function resetToOriginal() {
    name = originalName;
    totalLimit = originalTotalLimit;
    totalIncome = originalTotalIncome;
    errors = { name: "", totalLimit: "", totalIncome: "" };
  }

  /**
   * Initialize form when budgetData changes
   */
  $: if (budgetData && mode === "edit") {
    originalName = budgetData.name || "";
    originalTotalLimit = String(budgetData.totalLimit);
    originalTotalIncome = budgetData.totalIncome
      ? String(budgetData.totalIncome)
      : "";

    name = originalName;
    totalLimit = originalTotalLimit;
    totalIncome = originalTotalIncome;
  }

  /**
   * Reset form when modal closes (only in create mode)
   */
  $: if (!isOpen && mode === "create") {
    resetForm();
  }

  function resetForm() {
    name = "";
    totalLimit = "";
    totalIncome = "";
    errors = { name: "", totalLimit: "", totalIncome: "" };
  }

  function closeModal() {
    isOpen = false;
    resetForm();
    dispatch("close");
  }

  function handleSubmit() {
    if (!validateForm()) {
      return;
    }

    const budgetData = {
      name: name.trim() || undefined,
      totalLimit: Number(totalLimit),
      totalIncome: totalIncome ? Number(totalIncome) : undefined,
    };

    dispatch("submit", budgetData);
    closeModal();
  }

  // Handle keyboard shortcuts
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit();
    } else if (e.key === "Escape") {
      closeModal();
    }
  }
</script>

<CustomModal
  bind:open={isOpen}
  title={mode === "create" ? "Create Budget" : "Edit Budget"}
  description={mode === "create"
    ? "Set up a new budget for the selected month"
    : "Update your budget settings"}
  className="max-w-md"
>
  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
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

    <!-- Financial Overview Section -->
    <div
      class="space-y-4 p-4 bg-paper-light dark:bg-paper-dark/30 rounded-lg border-2 border-line-light dark:border-line-dark"
    >
      <h3 class="text-base font-semibold text-primary dark:text-white mb-3">
        Financial Overview
      </h3>

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

      <!-- Total Income (Optional) -->
      <div class="space-y-2">
        <Label
          for="totalIncome"
          class="text-base font-medium text-primary dark:text-white"
        >
          Total Income
        </Label>
        <div class="relative">
          <span
            class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-600 dark:text-gray-400 font-handwriting text-xl"
          >
            {currencySymbol}
          </span>
          <Input
            id="totalIncome"
            type="number"
            step="0.01"
            min="0"
            bind:value={totalIncome}
            placeholder="e.g., 6000"
            class="pl-10 font-handwriting text-xl border-2 border-line-light dark:border-line-dark focus:border-accent-gold focus:ring-accent-gold/20"
            aria-describedby="totalIncome-description"
            aria-invalid={errors.totalIncome ? "true" : "false"}
          />
        </div>
        <p
          id="totalIncome-description"
          class="text-xs text-gray-500 dark:text-gray-400 italic"
        >
          Optional: Enter your total income for this month
        </p>
        {#if errors.totalIncome}
          <p
            class="text-sm text-red-500 dark:text-red-400 font-medium"
            role="alert"
          >
            {errors.totalIncome}
          </p>
        {/if}
      </div>

      <!-- Savings Calculation -->
      {#if savings !== null}
        <div
          class="p-4 rounded-lg {savings >= 0
            ? 'bg-green-50 dark:bg-green-900/20'
            : 'bg-red-50 dark:bg-red-900/20'}"
          role="status"
          aria-live="polite"
        >
          <div class="flex justify-between items-center">
            <div class="flex items-center space-x-2">
              <span
                class="material-icons-outlined {savings >= 0
                  ? 'text-green-600'
                  : 'text-red-600'}"
              >
                {savings >= 0 ? "savings" : "warning"}
              </span>
              <span class="font-medium text-primary dark:text-white">
                {savings >= 0 ? "Projected Savings" : "Over Budget"}
              </span>
            </div>
            <span
              class="text-lg font-bold {savings >= 0
                ? 'text-green-600'
                : 'text-red-600'}"
            >
              {currencySymbol}{Math.abs(savings).toFixed(2)}
            </span>
          </div>
        </div>
      {/if}
    </div>

    <!-- Edit Mode: Budget Health Status -->
    {#if mode === "edit" && budgetData}
      <div
        class="space-y-3 p-4 bg-paper-light dark:bg-paper-dark/30 rounded-lg border-2 border-line-light dark:border-line-dark"
      >
        <h3 class="text-base font-semibold text-primary dark:text-white mb-3">
          Budget Health
        </h3>

        <!-- Progress Bar -->
        <div class="space-y-2">
          <div class="flex justify-between text-sm">
            <span class="text-gray-600 dark:text-gray-400">Spent</span>
            <span class="font-medium text-primary dark:text-white">
              {formatCurrency(currentSpending)} / {formatCurrency(
                Number(totalLimit) || 0
              )}
            </span>
          </div>
          <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
            <div
              class="h-2.5 rounded-full transition-all duration-300 {budgetHealth ===
              'healthy'
                ? 'bg-green-500'
                : budgetHealth === 'warning'
                  ? 'bg-yellow-500'
                  : 'bg-red-500'}"
              style="width: {spendingPercentage}%"
            ></div>
          </div>
          <div class="text-right text-xs text-gray-500 dark:text-gray-400">
            {spendingPercentage.toFixed(0)}% used
          </div>
        </div>
      </div>
    {/if}

    <!-- Warning Banners -->
    {#if warnings.length > 0}
      <div class="space-y-2">
        {#each warnings as warning}
          <div
            class="p-4 bg-amber-50 dark:bg-amber-900/20 border-l-4 border-amber-500 rounded-r"
            role="alert"
          >
            <div class="flex items-start space-x-2">
              <span
                class="material-icons-outlined text-amber-600 dark:text-amber-400 mt-0.5"
              >
                warning
              </span>
              <p class="text-sm text-amber-800 dark:text-amber-200">
                {warning}
              </p>
            </div>
          </div>
        {/each}
      </div>
    {/if}

    <!-- Actions -->
    <div
      class="flex justify-end space-x-3 pt-6 border-t-2 border-line-light dark:border-line-dark"
    >
      <!-- Reset to Original (Edit Mode Only) -->
      {#if mode === "edit"}
        <Button
          type="button"
          variant="outline"
          onclick={resetToOriginal}
          disabled={isSaving}
          class="font-medium"
        >
          <span class="flex items-center space-x-2">
            <span class="material-icons-outlined text-sm">restore</span>
            <span>Reset</span>
          </span>
        </Button>
      {/if}

      <Button
        type="button"
        variant="outline"
        onclick={closeModal}
        disabled={isSaving}
        class="font-medium"
      >
        Cancel
      </Button>
      <Button
        type="submit"
        disabled={isSaving || !formValid}
        class="min-w-[120px] font-medium"
      >
        {#if isSaving}
          <span class="flex items-center space-x-2">
            <span class="material-icons-outlined animate-spin text-sm"
              >refresh</span
            >
            <span>Saving...</span>
          </span>
        {:else}
          <span class="flex items-center space-x-2">
            <span class="material-icons-outlined text-sm">
              {mode === "create" ? "add_circle" : "edit"}
            </span>
            <span>{mode === "create" ? "Create Budget" : "Save Changes"}</span>
          </span>
        {/if}
      </Button>
    </div>
  </form>
</CustomModal>
