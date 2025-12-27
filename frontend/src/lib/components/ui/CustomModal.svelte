<script lang="ts">
	import { cn } from '$lib/utils/cn';
	import { createEventDispatcher } from 'svelte';

	export let open = false;
	export let title = '';
	export let description = '';
	export let className = '';

	const dispatch = createEventDispatcher();

	function close() {
		open = false;
		dispatch('close');
	}

	// Close on escape key
	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') close();
	}

	// Close on backdrop click
	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) close();
	}

	$: if (open) {
		if (typeof document !== 'undefined') {
			document.addEventListener('keydown', handleKeydown);
			// Prevent body scroll when modal is open
			document.body.style.overflow = 'hidden';
		}
	} else {
		if (typeof document !== 'undefined') {
			document.removeEventListener('keydown', handleKeydown);
			document.body.style.overflow = '';
		}
	}
</script>

{#if open}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		onclick={handleBackdropClick}
		role="dialog"
		aria-modal="true"
	>
		<div
			class={cn(
				'relative w-full max-w-lg bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-6 m-4 max-h-[90vh] overflow-y-auto',
				className
			)}
		>
			<div class="flex items-start justify-between mb-4">
				<div class="flex-1">
					<h2 class="text-2xl font-display font-bold text-primary dark:text-white">
						{title}
					</h2>
					{#if description}
						<p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{description}</p>
					{/if}
				</div>
				<button
					onclick={close}
					class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 p-1 rounded-md hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
					aria-label="Close"
				>
					<span class="material-icons-outlined">close</span>
				</button>
			</div>

			<div class="mt-4">
				<slot />
			</div>
		</div>
	</div>
{/if}
