<script lang="ts">
  import { page } from "$app/stores";
  import { budgetsApi } from "$lib/api/budgets";
  import BudgetFormModal from "$lib/components/budget/BudgetFormModal.svelte";
  import ShareBudgetDialog from "$lib/components/sharing/ShareBudgetDialog.svelte";
  import { Button } from "$lib/components/ui/button";
  import { Skeleton } from "$lib/components/ui/skeleton";
  import {
    createBudgetForCurrentMonth,
    currentMonthBudget,
    totalIncome,
    totalSpent,
  } from "$lib/stores";
  import { budgetsLoading } from "$lib/stores/budgets";
  import { showToast } from "$lib/stores/ui";
  import { formatCurrency, formatMonthYear } from "$lib/utils/format";

  let shareDialogOpen = false;
  let showBudgetFormModal = false;
  let formMode: "create" | "edit" = "create";
  let isSavingBudget = false;

  async function handleCreateBudget(data: {
    name?: string;
    totalLimit: number;
    totalIncome?: number;
  }): Promise<void> {
    const userId = $page.data.userId;
    if (!userId) {
      console.error("[Overview] No userId available for budget creation");
      showToast("Failed to create budget: User not authenticated", "error");
      return;
    }

    isSavingBudget = true;
    try {
      await createBudgetForCurrentMonth(userId, {
        name: data.name,
        totalLimit: data.totalLimit,
        totalIncome: data.totalIncome,
      });

      showToast("Budget created successfully!", "success");
      showBudgetFormModal = false;
    } catch (error) {
      console.error("[Overview] Failed to create budget:", error);
      showToast(
        error instanceof Error ? error.message : "Failed to create budget",
        "error"
      );
    } finally {
      isSavingBudget = false;
    }
  }

  async function handleUpdateBudget(data: {
    name?: string;
    totalLimit: number;
    totalIncome?: number;
  }): Promise<void> {
    if (!$currentMonthBudget) {
      console.error("[Overview] No current budget available for update");
      showToast("Failed to update budget: Budget not found", "error");
      return;
    }

    isSavingBudget = true;
    try {
      await budgetsApi.updateBudget($currentMonthBudget.id, data);
      showToast("Budget updated successfully!", "success");
      showBudgetFormModal = false;
      // Reload budgets to get updated data
      const { loadBudgets } = await import("$lib/stores/budgets");
      await loadBudgets();
    } catch (error) {
      console.error("[Overview] Failed to update budget:", error);
      showToast(
        error instanceof Error ? error.message : "Failed to update budget",
        "error"
      );
    } finally {
      isSavingBudget = false;
    }
  }

  function openCreateModal() {
    formMode = "create";
    showBudgetFormModal = true;
  }

  function openEditModal() {
    formMode = "edit";
    showBudgetFormModal = true;
  }

  function handleFormSubmit(event: CustomEvent) {
    const data = event.detail;
    if (formMode === "create") {
      handleCreateBudget(data);
    } else {
      handleUpdateBudget(data);
    }
  }
</script>

<div class="space-y-6">
  <!-- Page Header -->
  <div class="mb-6 flex justify-between items-start">
    <div>
      <h2
        class="text-3xl font-display font-bold text-primary dark:text-white mb-2"
      >
        Budget Overview
      </h2>
      <p class="text-sm text-gray-500 dark:text-gray-400 italic">
        "Classify and summarize expenditures."
      </p>
    </div>
    {#if $currentMonthBudget}
      <div class="flex space-x-2">
        <Button
          onclick={openEditModal}
          aria-label="Edit this budget"
          class="shrink-0"
          variant="outline"
        >
          <span class="material-icons-outlined text-sm mr-1">edit</span>
          Edit
        </Button>
        <Button
          onclick={() => (shareDialogOpen = true)}
          aria-label="Share this budget"
          class="shrink-0"
        >
          <span class="material-icons-outlined text-sm mr-1">share</span>
          Share
        </Button>
      </div>
    {/if}
  </div>

  {#if $budgetsLoading}
    <!-- Loading State - Skeleton -->
    <section
      class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark overflow-hidden mb-6"
    >
      <div
        class="bg-primary dark:bg-gray-700 text-white p-3 flex justify-between items-center"
      >
        <h3 class="font-display font-semibold tracking-wide uppercase text-sm">
          Budget Review
        </h3>
      </div>
      <div class="p-4 grid grid-cols-2 gap-3 text-center">
        <div
          class="border border-line-light dark:border-line-dark p-2 rounded-lg"
        >
          <Skeleton variant="card" className="h-12 w-full" />
        </div>
        <div
          class="border border-line-light dark:border-line-dark p-2 rounded-lg"
        >
          <Skeleton variant="card" className="h-12 w-full" />
        </div>
        <div
          class="border border-line-light dark:border-line-dark p-2 rounded-lg"
        >
          <Skeleton variant="card" className="h-12 w-full" />
        </div>
        <div class="border-2 p-2 rounded-lg">
          <Skeleton variant="card" className="h-12 w-full" />
        </div>
      </div>
    </section>

    <section
      class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-5"
    >
      <Skeleton variant="card" className="h-64 w-full" />
    </section>
  {:else if $currentMonthBudget}
    <!-- Budget Review Card -->
    <section
      class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark overflow-hidden mb-6"
    >
      <div
        class="bg-primary dark:bg-gray-700 text-white p-3 flex justify-between items-center"
      >
        <h3 class="font-display font-semibold tracking-wide uppercase text-sm">
          Budget Review
        </h3>
        <span class="material-icons text-sm opacity-80">analytics</span>
      </div>
      <div class="p-4 grid grid-cols-2 gap-3 text-center">
        <div
          class="border border-line-light dark:border-line-dark p-2 rounded-lg bg-gray-50 dark:bg-gray-900/50"
        >
          <p
            class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1"
          >
            Budget Limit
          </p>
          <p class="font-bold text-gray-800 dark:text-gray-100 text-lg">
            {formatCurrency($currentMonthBudget.totalLimit)}
          </p>
        </div>
        <div
          class="border border-line-light dark:border-line-dark p-2 rounded-lg bg-gray-50 dark:bg-gray-900/50"
        >
          <p
            class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1"
          >
            Total Income
          </p>
          <p class="font-bold text-green-600 dark:text-green-400 text-lg">
            {formatCurrency($totalIncome)}
          </p>
        </div>
        <div
          class="border border-line-light dark:border-line-dark p-2 rounded-lg bg-gray-50 dark:bg-gray-900/50"
        >
          <p
            class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1"
          >
            Total Expenses
          </p>
          <p class="font-bold text-red-500 dark:text-red-400 text-lg">
            {formatCurrency($totalSpent)}
          </p>
        </div>
        <div
          class="border-2 border-accent-gold p-2 rounded-lg bg-accent-highlight/30 dark:bg-yellow-900/20"
        >
          <p
            class="text-xs text-yellow-700 dark:text-yellow-400 uppercase tracking-wider mb-1 font-semibold"
          >
            Remaining
          </p>
          <p class="font-bold text-yellow-800 dark:text-yellow-300 text-lg">
            {formatCurrency($currentMonthBudget.totalLimit - $totalSpent)}
          </p>
        </div>
      </div>
    </section>

    <!-- Monthly Reflection Section -->
    <section
      class="relative bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-5 overflow-hidden"
    >
      <div
        class="absolute top-0 left-0 bottom-0 w-6 flex flex-col justify-evenly py-4 pl-1 pointer-events-none"
      >
        <div
          class="w-3 h-3 rounded-full bg-gray-300 dark:bg-gray-600 shadow-inner mb-2"
        ></div>
        <div
          class="w-3 h-3 rounded-full bg-gray-300 dark:bg-gray-600 shadow-inner mb-2"
        ></div>
        <div
          class="w-3 h-3 rounded-full bg-gray-300 dark:bg-gray-600 shadow-inner mb-2"
        ></div>
      </div>
      <div class="pl-6">
        <h3
          class="font-display text-lg font-bold text-primary dark:text-gray-100 mb-4 border-b-2 border-primary dark:border-gray-500 inline-block"
        >
          Monthly Reflection
        </h3>
        <div class="mb-5">
          <div
            class="block text-xs font-bold uppercase tracking-wide text-white bg-primary dark:bg-gray-600 py-1 px-2 mb-1 w-full rounded-sm"
          >
            My biggest wins this month
          </div>
          <div
            class="notebook-lines min-h-[4rem] text-sm text-blue-600 dark:text-blue-300 font-handwriting leading-8 pl-1 bg-white dark:bg-gray-800 rounded border border-line-light dark:border-line-dark p-2"
          >
            What went well with your budget?
          </div>
        </div>
        <div class="mb-5">
          <div
            class="block text-xs font-bold uppercase tracking-wide text-white bg-primary dark:bg-gray-600 py-1 px-2 mb-1 w-full rounded-sm"
          >
            Did I meet my budget? If not, why not?
          </div>
          <div
            class="notebook-lines min-h-[4rem] text-sm text-gray-700 dark:text-gray-300 font-handwriting leading-8 pl-1 bg-white dark:bg-gray-800 rounded border border-line-light dark:border-line-dark p-2"
          >
            Reflect on your spending...
          </div>
        </div>
        <div class="mb-2">
          <div
            class="block text-xs font-bold uppercase tracking-wide text-white bg-primary dark:bg-gray-600 py-1 px-2 mb-1 w-full rounded-sm"
          >
            I will do this within 1 month to improve
          </div>
          <div
            class="notebook-lines min-h-[4rem] text-sm text-gray-700 dark:text-gray-300 font-handwriting leading-8 pl-1 bg-white dark:bg-gray-800 rounded border border-line-light dark:border-line-dark p-2"
          >
            What improvements will you make?
          </div>
        </div>
      </div>
    </section>
  {:else}
    <!-- No Budget State -->
    <section
      class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-8 text-center"
    >
      <span
        class="material-icons-outlined text-6xl text-gray-300 dark:text-gray-600 mb-4"
      >
        assignment
      </span>
      <h3
        class="text-xl font-display font-bold text-primary dark:text-white mb-2"
      >
        No Budget Yet
      </h3>
      <p class="text-gray-500 dark:text-gray-400 mb-6">
        Create a budget for this month to start tracking your expenses.
      </p>
      <Button
        disabled={isSavingBudget}
        onclick={openCreateModal}
        aria-label="Create a new budget"
      >
        {isSavingBudget ? "Creating..." : "Create Budget"}
      </Button>
    </section>
  {/if}
</div>

{#if $currentMonthBudget}
  <ShareBudgetDialog
    bind:open={shareDialogOpen}
    budgetId={$currentMonthBudget.id}
    budgetName={formatMonthYear(new Date($currentMonthBudget.month + "-01"))}
  />
{/if}

<BudgetFormModal
  bind:isOpen={showBudgetFormModal}
  mode={formMode}
  budgetData={$currentMonthBudget
    ? {
        id: $currentMonthBudget.id,
        userId: $currentMonthBudget.userId,
        month: $currentMonthBudget.month,
        totalLimit: $currentMonthBudget.totalLimit,
        totalIncome: $currentMonthBudget.totalIncome,
        spent: $totalSpent,
        remaining: $currentMonthBudget.totalLimit - $totalSpent,
        createdAt: $currentMonthBudget.createdAt,
        updatedAt: $currentMonthBudget.updatedAt,
      }
    : undefined}
  bind:isSaving={isSavingBudget}
  currentSpending={$totalSpent}
  on:close={() => (showBudgetFormModal = false)}
  on:submit={handleFormSubmit}
/>
