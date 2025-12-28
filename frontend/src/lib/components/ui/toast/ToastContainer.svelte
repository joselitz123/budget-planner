<script lang="ts">
	import { toasts } from '$lib/stores/ui';
	import Toast from './Toast.svelte';
</script>

<div class="toast-container">
	{#each $toasts as toast (toast.id)}
		<Toast {toast} />
	{/each}
</div>

<style>
	.toast-container {
		position: fixed;
		bottom: 1rem;
		right: 1rem;
		z-index: 9999;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		pointer-events: none;
	}

	.toast-container :global(*) {
		pointer-events: auto;
	}

	/* Mobile: Position above bottom navigation */
	@media (max-width: 768px) {
		.toast-container {
			bottom: 5rem;
			right: 0.5rem;
			left: 0.5rem;
		}
	}
</style>
