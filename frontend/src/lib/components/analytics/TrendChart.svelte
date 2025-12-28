<script lang="ts">
	import type { MonthlyTrend } from '$lib/stores/analytics';

	export let data: MonthlyTrend[] = [];
	export let height: number = 160; // Default height in pixels

	// Calculate SVG path coordinates
	$: maxAmount = Math.max(...data.map((d) => d.amount), 1);

	$: points = data.map((d, i) => {
		const x = data.length > 1 ? (i / (data.length - 1)) * 100 : 50; // 0-100% width
		const y = 100 - (d.amount / maxAmount) * 100; // Invert Y (100 is top in SVG)
		return { x, y, amount: d.amount, month: d.monthLabel };
	});

	$: pathData =
		points.length > 0
			? `M ${points.map((p) => `${p.x},${p.y}`).join(' L ')}`
			: '';

	// Grid line Y positions
	const gridLines = [25, 50, 75];

	// Empty state
	$: isEmpty = data.length === 0;
</script>

<div class="w-full">
	{#if isEmpty}
		<!-- Empty state -->
		<div
			class="flex flex-col items-center justify-center bg-gray-50 dark:bg-gray-800 rounded-xl border-2 border-dashed border-gray-300 dark:border-gray-600"
			style="height: {height}px;"
		>
			<span class="material-icons-outlined text-4xl text-gray-400 dark:text-gray-500 mb-2"
				>trending_up</span
			>
			<p class="text-sm text-gray-500 dark:text-gray-400">No trend data yet</p>
		</div>
	{:else}
		<!-- Trend chart -->
		<div class="w-full" style="height: {height}px;">
			<svg
				viewBox="0 0 100 100"
				class="w-full h-full"
				preserveAspectRatio="none"
				aria-label="Monthly spending trend chart"
			>
				<!-- Grid lines (horizontal) -->
				{#each gridLines as y}
					<line
						x1="0"
						y1={y}
						x2="100"
						y2={y}
						stroke="currentColor"
						stroke-width="0.5"
						stroke-dasharray="2"
						class="text-gray-200 dark:text-gray-700"
					/>
				{/each}

				<!-- Trend line -->
				<path
					d={pathData}
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					class="text-gray-900 dark:text-gray-100"
				/>

				<!-- Data points -->
				{#each points as point}
					<circle
						cx={point.x}
						cy={point.y}
						r="2.5"
						fill="currentColor"
						class="text-gray-900 dark:text-gray-100"
					>
						<title>{point.month}: ₱{point.amount.toLocaleString('en-PH')}</title>
					</circle>
				{/each}
			</svg>
		</div>

		<!-- X-axis labels (months) -->
		<div class="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-2 px-1">
			{#each points as point}
				<span class="text-center" style="width: {100 / points.length}%;">
					{point.month}
				</span>
			{/each}
		</div>
	{/if}
</div>
