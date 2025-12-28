<script lang="ts">
	// Error page - no imports needed
</script>

<div class="flex flex-col items-center justify-center min-h-[50vh] px-4">
	<div class="text-center">
		<span class="material-icons-outlined text-6xl text-red-500 mb-4">error_outline</span>
		<h2 class="text-2xl font-display font-bold text-gray-900 dark:text-gray-100 mb-2">
			Oops! Something went wrong
		</h2>
		<p class="text-gray-600 dark:text-gray-400 mb-6">
			We couldn't load your analytics. Please try again.
		</p>
		<button
			onclick={() => window.location.reload()}
			class="inline-flex items-center px-4 py-2 bg-primary dark:bg-gray-700 text-white rounded-lg hover:opacity-90 transition-opacity"
		>
			<span class="material-icons text-sm mr-2">refresh</span>
			Reload Page
		</button>
	</div>
</div>

<style>
	button {
		cursor: pointer;
	}
</style>
