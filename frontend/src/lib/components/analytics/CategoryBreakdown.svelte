<script lang="ts">
	import type { CategorySpend } from '$lib/stores/analytics';

	export let data: CategorySpend[] = [];
	export let limit: number = 5;

	// Limit data and calculate total for percentages
	$: limitedData = data.slice(0, limit);

	$: totalAmount = data.reduce((sum, item) => sum + item.amount, 0);

	// Empty state
	$: isEmpty = data.length === 0;
</script>

<div class="w-full">
	{#if isEmpty}
		<!-- Empty state -->
		<div class="flex flex-col items-center justify-center py-8 text-center">
			<span class="material-icons-outlined text-5xl text-gray-300 dark:text-gray-600 mb-3"
				>category</span
			>
			<p class="text-gray-500 dark:text-gray-400 font-medium">No spending categories yet</p>
			<p class="text-sm text-gray-400 dark:text-gray-500 mt-1">
				Add expenses to see your top categories
			</p>
		</div>
	{:else}
		<!-- Category breakdown list -->
		<div class="space-y-4">
			{#each limitedData as item, index}
				<div class="flex items-center space-x-3">
					<!-- Rank badge -->
					<div
						class="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold flex-shrink-0
							{index === 0
							? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
							: index === 1
							? 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
							: index === 2
							? 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200'
							: 'bg-gray-50 text-gray-600 dark:bg-gray-800 dark:text-gray-400'}"
					>
						{index + 1}
					</div>

					<!-- Category icon -->
					<div
						class="w-10 h-10 rounded-full flex items-center justify-center text-xl flex-shrink-0
							bg-gradient-to-br from-gray-100 to-gray-200 dark:from-gray-700 dark:to-gray-600"
						style="background-color: {item.color}20;"
					>
						<span>{item.icon}</span>
					</div>

					<!-- Category details -->
					<div class="flex-1 min-w-0">
						<!-- Name and amount -->
						<div class="flex justify-between items-center mb-1">
							<span class="font-medium text-sm text-gray-900 dark:text-gray-100 truncate">
								{item.name}
							</span>
							<span class="text-sm font-bold text-gray-900 dark:text-gray-100 flex-shrink-0 ml-2">
								₱{item.amount.toLocaleString('en-PH', { maximumFractionDigits: 0 })}
							</span>
						</div>

						<!-- Progress bar -->
						<div class="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
							<div
								class="h-full rounded-full transition-all duration-300 ease-out"
								style="width: {item.percentage}%; background-color: {item.color};"
							></div>
						</div>

						<!-- Percentage of total -->
						<div class="flex justify-between mt-1">
							<span class="text-xs text-gray-500 dark:text-gray-400">
								{item.percentage}% of total spending
							</span>
						</div>
					</div>
				</div>
			{/each}
		</div>

		<!-- Show more indicator if data exceeds limit -->
		{#if data.length > limit}
			<div class="text-center mt-4">
				<span class="text-sm text-gray-500 dark:text-gray-400">
					Showing top {limit} of {data.length} categories
				</span>
			</div>
		{/if}
	{/if}
</div>
