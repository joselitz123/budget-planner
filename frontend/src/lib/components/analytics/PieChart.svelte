<script lang="ts">
	import type { CategorySpend } from '$lib/stores/analytics';

	export let data: CategorySpend[] = [];
	export let size: 'sm' | 'md' | 'lg' = 'md';

	// Size variants
	const sizeClasses = {
		sm: 'w-32 h-32',
		md: 'w-40 h-40',
		lg: 'w-56 h-56'
	};

	const donutHoleSizes = {
		sm: 'w-12 h-12',
		md: 'w-16 h-16',
		lg: 'w-20 h-20'
	};

	// Calculate conic-gradient stops
	$: gradientStops = data.reduce((stops, item, index) => {
		const prevPercent = index > 0 ? data[index - 1].cumulativePercent : 0;
		const currentPercent = prevPercent + item.percentage;
		return [...stops, `${item.color} ${prevPercent}% ${currentPercent}%`];
	}, [] as string[]);

	$: gradientStyle =
		data.length > 0
			? `conic-gradient(${gradientStops.join(', ')})`
			: 'conic-gradient(#e5e7eb 0% 100%)';

	// Calculate total amount
	$: totalAmount = data.reduce((sum, item) => sum + item.amount, 0);

	// Empty state
	$: isEmpty = data.length === 0;
</script>

<div class="flex flex-col items-center justify-center space-y-4">
	<!-- Pie/Donut Chart -->
	<div class="relative {sizeClasses[size]} flex-shrink-0">
		{#if isEmpty}
			<!-- Empty state placeholder -->
			<div
				class="w-full h-full rounded-full border-4 border-white dark:border-gray-800 bg-gray-100 dark:bg-gray-700 flex items-center justify-center"
			>
				<div
					class="{donutHoleSizes[size]} bg-white dark:bg-gray-800 rounded-full flex items-center justify-center"
				>
					<span class="material-icons-outlined text-gray-400">pie_chart</span>
				</div>
			</div>
		{:else}
			<!-- Pie chart with conic-gradient -->
			<div
				style="background: {gradientStyle};"
				class="w-full h-full rounded-full border-4 border-white dark:border-gray-800 shadow-lg relative"
			>
				<!-- Donut hole -->
				<div
					class="{donutHoleSizes[size]} absolute inset-0 m-auto bg-white dark:bg-gray-800 rounded-full flex items-center justify-center"
				>
					<div class="text-center">
						<div class="text-xs text-gray-500 dark:text-gray-400">Total</div>
						<div class="text-sm font-bold text-gray-900 dark:text-gray-100">
							₱{totalAmount.toLocaleString('en-PH', { maximumFractionDigits: 0 })}
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>

	<!-- Legend -->
	{#if isEmpty}
		<div class="text-center text-sm text-gray-500 dark:text-gray-400">
			<p>No spending data</p>
		</div>
	{:else}
		<div class="w-full space-y-2">
			{#each data as item}
				<div class="flex items-center justify-between text-sm">
					<div class="flex items-center space-x-2">
						<!-- Color swatch -->
						<div
							class="w-3 h-3 rounded-full flex-shrink-0"
							style="background-color: {item.color};"
						></div>
						<!-- Category icon and name -->
						<div class="flex items-center space-x-1">
							<span>{item.icon}</span>
							<span class="text-gray-700 dark:text-gray-300">{item.name}</span>
						</div>
					</div>
					<!-- Amount and percentage -->
					<div class="text-right">
						<div class="font-medium text-gray-900 dark:text-gray-100">
							₱{item.amount.toLocaleString('en-PH', { maximumFractionDigits: 0 })}
						</div>
						<div class="text-xs text-gray-500 dark:text-gray-400">{item.percentage}%</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
