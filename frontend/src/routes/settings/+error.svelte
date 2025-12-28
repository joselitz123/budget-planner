<script lang="ts">
	import { browser } from '$app/environment';

	export let data;
	export let status;
	export let error;

	const isDev = browser && import.meta.env.DEV;
</script>

<div class="min-h-screen bg-background-light dark:bg-background-dark font-body antialiased flex items-center justify-center p-4">
	<div
		class="max-w-md w-full bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-8">
		<div class="text-center">
			<!-- Error Icon -->
			<span class="material-icons-outlined text-6xl text-gray-300 dark:text-gray-600 mb-4">
				settings
			</span>

			<!-- Error Heading -->
			<h1 class="text-3xl font-display font-bold text-primary dark:text-white mb-2">
				Settings Error
			</h1>

			<!-- Error Message -->
			<p class="text-gray-600 dark:text-gray-400 mb-6 font-handwriting text-xl">
				Something went wrong while loading settings. This doesn't affect your data.
			</p>

			<!-- Error Details (Dev Only) -->
			{#if isDev && error}
				<details
					class="mb-6 p-4 bg-red-50 dark:bg-red-900/20 rounded-lg border border-red-200 dark:border-red-800 text-left">
					<summary class="cursor-pointer text-sm font-medium text-red-600 dark:text-red-400">
						Error Details
					</summary>
					<div class="mt-3 text-xs">
						<pre class="overflow-auto text-gray-600 dark:text-gray-400">{error?.message || 'Unknown error'}</pre>
					</div>
				</details>
			{/if}

			<!-- Action Buttons -->
			<div class="flex flex-col sm:flex-row gap-3 justify-center">
				<a
					href="/settings"
					class="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2">
					<span class="material-icons-outlined text-sm mr-2">refresh</span>
					Back to Settings
				</a>
				<a
					href="/"
					class="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 border border-input bg-background hover:bg-accent hover:text-accent-foreground h-10 px-4 py-2">
					Go to Dashboard
				</a>
			</div>
		</div>
	</div>
</div>

<style>
	:global(pre) {
		white-space: pre-wrap;
		word-wrap: break-word;
	}
</style>
