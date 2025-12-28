<script lang="ts">
	import { page } from '$app/stores';
	import {
		spendingByCategory,
		topCategories,
		monthlyTrends,
		avgDailySpending,
		totalSpentThisMonth,
		topCategory,
		budgetRemaining,
		analyticsLoading
	} from '$lib/stores/analytics';
	import { currentBudget } from '$lib/stores/budgets';
	import { currentMonthBudget } from '$lib/stores';
	import { budgetsLoading } from '$lib/stores/budgets';
	import { formatCurrency } from '$lib/utils/format';
	import PieChart from '$lib/components/analytics/PieChart.svelte';
	import TrendChart from '$lib/components/analytics/TrendChart.svelte';
	import CategoryBreakdown from '$lib/components/analytics/CategoryBreakdown.svelte';
	import { Skeleton } from '$lib/components/ui/skeleton';
</script>

<svelte:head>
	<title>Analytics - Budget Planner</title>
	<meta name="description" content="View your spending analytics and insights" />
</svelte:head>

<div class="space-y-6 pb-20">
	<!-- Page Header -->
	<div class="mb-6">
		<h2 class="text-3xl font-display font-bold text-primary dark:text-white mb-2">
			Analytics & Insights
		</h2>
		<p class="text-sm text-gray-500 dark:text-gray-400 italic">
			"Understand your spending patterns."
		</p>
	</div>

	{#if $analyticsLoading || $budgetsLoading}
		<!-- Loading State - Skeletons -->
		<div class="space-y-6">
			<!-- Summary Cards -->
			<div class="grid grid-cols-2 gap-4">
				{#each [1, 2, 3, 4] as _}
					<div
						class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-4"
					>
						<Skeleton variant="text" className="h-4 w-20 mb-2" />
						<Skeleton variant="card" className="h-8 w-full" />
					</div>
				{/each}
			</div>

			<!-- Charts -->
			<section
				class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-5"
			>
				<Skeleton variant="card" className="h-64 w-full mb-4" />
				<Skeleton variant="card" className="h-48 w-full" />
			</section>
		</div>
	{:else}
		<!-- Summary Cards Row -->
		<div class="grid grid-cols-2 gap-4 mb-6">
			<!-- Total Spent This Month -->
			<div
				class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-4 relative overflow-hidden"
			>
				<div class="absolute top-0 right-0 w-20 h-20 opacity-5">
					<span class="material-icons text-6xl">payments</span>
				</div>
				<p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">
					Total Spent
				</p>
				<p class="font-bold text-2xl text-red-500 dark:text-red-400">
					{formatCurrency($totalSpentThisMonth)}
				</p>
				<p class="text-xs text-gray-400 dark:text-gray-500 mt-1">This month</p>
			</div>

			<!-- Top Spending Category -->
			<div
				class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-4 relative overflow-hidden"
			>
				<div class="absolute top-0 right-0 w-20 h-20 opacity-5">
					<span class="material-icons text-6xl">stars</span>
				</div>
				<p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">
					Top Category
				</p>
				{#if $topCategory}
					<div class="flex items-center space-x-2">
						<span class="text-2xl">{$topCategory.icon}</span>
						<p class="font-bold text-lg text-gray-900 dark:text-gray-100 truncate">
							{$topCategory.name}
						</p>
					</div>
					<p class="text-xs text-gray-400 dark:text-gray-500 mt-1">
						{$topCategory.percentage}% of spending
					</p>
				{:else}
					<p class="font-bold text-lg text-gray-400 dark:text-gray-500">-</p>
					<p class="text-xs text-gray-400 dark:text-gray-500 mt-1">No data yet</p>
				{/if}
			</div>

			<!-- Average Daily Spending -->
			<div
				class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-4 relative overflow-hidden"
			>
				<div class="absolute top-0 right-0 w-20 h-20 opacity-5">
					<span class="material-icons text-6xl">today</span>
				</div>
				<p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">
					Avg Daily
				</p>
				<p class="font-bold text-2xl text-blue-600 dark:text-blue-400">
					{formatCurrency($avgDailySpending)}
				</p>
				<p class="text-xs text-gray-400 dark:text-gray-500 mt-1">Per day</p>
			</div>

			<!-- Budget Remaining -->
			<div
				class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-4 relative overflow-hidden"
			>
				<div class="absolute top-0 right-0 w-20 h-20 opacity-5">
					<span class="material-icons text-6xl">account_balance_wallet</span>
				</div>
				<p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">
					Budget Left
				</p>
				<p
					class="font-bold text-2xl {$budgetRemaining >= 0
						? 'text-green-600 dark:text-green-400'
						: 'text-red-500 dark:text-red-400'}"
				>
					{formatCurrency($budgetRemaining)}
				</p>
				<p class="text-xs text-gray-400 dark:text-gray-500 mt-1">
					{$budgetRemaining >= 0 ? 'On track' : 'Over budget'}
				</p>
			</div>
		</div>

		<!-- Spending by Category (Pie Chart) -->
		<section
			class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark overflow-hidden relative mb-6"
		>
			<!-- Spiral Binding -->
			<div
				class="absolute left-0 top-0 bottom-0 w-6 bg-gray-100 dark:bg-gray-800 border-r border-line-light dark:border-line-dark flex flex-col items-center py-4 space-y-4 z-10"
			>
				<div
					class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 dark:from-gray-500 dark:to-gray-700 rounded-full shadow-sm"
				></div>
				<div
					class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 dark:from-gray-500 dark:to-gray-700 rounded-full shadow-sm"
				></div>
				<div
					class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 dark:from-gray-500 dark:to-gray-700 rounded-full shadow-sm"
				></div>
			</div>

			<div class="pl-8 p-5">
				<h3
					class="font-display font-semibold tracking-wide uppercase text-sm text-gray-700 dark:text-gray-300 mb-4 flex items-center"
				>
					<span class="material-icons mr-2 text-lg">pie_chart</span>
					Spending by Category
				</h3>
				<PieChart data={$spendingByCategory} size="lg" />
			</div>
		</section>

		<!-- Monthly Spending Trend -->
		<section
			class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark overflow-hidden relative mb-6"
		>
			<!-- Spiral Binding -->
			<div
				class="absolute left-0 top-0 bottom-0 w-6 bg-gray-100 dark:bg-gray-800 border-r border-line-light dark:border-line-dark flex flex-col items-center py-4 space-y-4 z-10"
			>
				<div
					class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 dark:from-gray-500 dark:to-gray-700 rounded-full shadow-sm"
				></div>
				<div
					class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 dark:from-gray-500 dark:to-gray-700 rounded-full shadow-sm"
				></div>
				<div
					class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 dark:from-gray-500 dark:to-gray-700 rounded-full shadow-sm"
				></div>
			</div>

			<div class="pl-8 p-5">
				<h3
					class="font-display font-semibold tracking-wide uppercase text-sm text-gray-700 dark:text-gray-300 mb-4 flex items-center"
				>
					<span class="material-icons mr-2 text-lg">show_chart</span>
					Monthly Spending Trend
				</h3>
				<TrendChart data={$monthlyTrends} height={160} />
			</div>
		</section>

		<!-- Top Spending Categories -->
		<section
			class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark overflow-hidden relative"
		>
			<!-- Spiral Binding -->
			<div
				class="absolute left-0 top-0 bottom-0 w-6 bg-gray-100 dark:bg-gray-800 border-r border-line-light dark:border-line-dark flex flex-col items-center py-4 space-y-4 z-10"
			>
				<div
					class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 dark:from-gray-500 dark:to-gray-700 rounded-full shadow-sm"
				></div>
				<div
					class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 dark:from-gray-500 dark:to-gray-700 rounded-full shadow-sm"
				></div>
				<div
					class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 dark:from-gray-500 dark:to-gray-700 rounded-full shadow-sm"
				></div>
			</div>

			<div class="pl-8 p-5">
				<h3
					class="font-display font-semibold tracking-wide uppercase text-sm text-gray-700 dark:text-gray-300 mb-4 flex items-center"
				>
					<span class="material-icons mr-2 text-lg">leaderboard</span>
					Top Spending Categories
				</h3>
				<CategoryBreakdown data={$topCategories} limit={5} />
			</div>
		</section>
	{/if}
</div>
