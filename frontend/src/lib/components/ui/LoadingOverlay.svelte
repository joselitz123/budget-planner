<script lang="ts">
	import { Spinner } from '$lib/components/ui/spinner';

	export let message: string | undefined = undefined;
	export let show = true;
</script>

{#if show}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-paper-light/80 dark:bg-paper-dark/80 backdrop-blur-sm"
		role="status"
		aria-busy="true"
		aria-live="polite"
	>
		<div class="flex flex-col items-center gap-4 p-8 bg-white dark:bg-gray-800 rounded-xl shadow-paper border border-line-light dark:border-line-dark">
			<Spinner size="xl" color="primary" />
			{#if message}
				<p class="text-sm text-gray-600 dark:text-gray-300 font-medium">{message}</p>
			{/if}
		</div>
	</div>
{/if}
