<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { tv, type VariantProps } from 'tailwind-variants';
	import type { Toast } from '$lib/stores/ui';
	import { removeToast } from '$lib/stores/ui';

	const toastVariants = tv({
		base: 'flex items-center gap-3 p-4 rounded-xl shadow-paper border min-w-[300px] max-w-md pointer-events-auto',
		variants: {
			type: {
				success: 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800',
				error: 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800',
				info: 'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800',
				warning: 'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800'
			}
		},
		defaultVariants: {
			type: 'info'
		}
	});

	type ToastVariants = VariantProps<typeof toastVariants>;

	export let toast: Toast;

	let isExiting = false;
	let remainingTime = toast.duration || 3000;
	let progressInterval: number;

	const iconMap = {
		success: 'check_circle',
		error: 'error',
		info: 'info',
		warning: 'warning'
	};

	$: type = toast.type;
	$: icon = iconMap[type];
	$: allClassName = toastVariants({ type });

	function dismiss() {
		isExiting = true;
		setTimeout(() => {
			removeToast(toast.id);
		}, 300); // Wait for exit animation
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			dismiss();
		}
	}

	onMount(() => {
		// Progress bar animation
		if (remainingTime > 0) {
			progressInterval = window.setInterval(() => {
				remainingTime -= 100;
				if (remainingTime <= 0) {
					dismiss();
				}
			}, 100);
		}

		// Keyboard listener
		document.addEventListener('keydown', handleKeydown);

		return () => {
			if (progressInterval) {
				clearInterval(progressInterval);
			}
			document.removeEventListener('keydown', handleKeydown);
		};
	});

	onDestroy(() => {
		if (progressInterval) {
			clearInterval(progressInterval);
		}
		document.removeEventListener('keydown', handleKeydown);
	});
</script>

<div
	class="toast-item {allClassName}"
	class:toast-exit={isExiting}
	role="alert"
	aria-live={type === 'error' || type === 'warning' ? 'assertive' : 'polite'}
>
	<!-- Icon -->
	<span class="material-icons text-xl flex-shrink-0 {toast.type === 'success'
		? 'text-green-600 dark:text-green-400'
		: toast.type === 'error'
			? 'text-red-600 dark:text-red-400'
			: toast.type === 'info'
				? 'text-blue-600 dark:text-blue-400'
				: 'text-yellow-600 dark:text-yellow-400'}">
		{icon}
	</span>

	<!-- Message -->
	<div class="flex-1 min-w-0">
		<p class="text-sm font-medium text-gray-900 dark:text-gray-100 leading-snug">
			{toast.message}
		</p>
		{#if toast.duration && toast.duration > 0}
			<div class="toast-progress mt-2 h-1 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
				<div
					class="toast-progress-bar h-full bg-current opacity-30"
					style="animation-duration: {toast.duration}ms; animation-name: progress;"
				></div>
			</div>
		{/if}
	</div>

	<!-- Close Button -->
	<button
		type="button"
		on:click={dismiss}
		aria-label="Close notification"
		class="flex-shrink-0 p-1 rounded hover:bg-black/5 dark:hover:bg-white/10 transition-colors">
		<span class="material-icons text-gray-500 dark:text-gray-400">close</span>
	</button>
</div>

<style>
	.toast-item {
		animation: slideInRight 0.3s ease-out forwards;
	}

	.toast-exit {
		animation: slideOutRight 0.3s ease-in forwards;
	}

	@keyframes slideInRight {
		from {
			transform: translateX(100%);
			opacity: 0;
		}
		to {
			transform: translateX(0);
			opacity: 1;
		}
	}

	@keyframes slideOutRight {
		from {
			transform: translateX(0);
			opacity: 1;
		}
		to {
			transform: translateX(100%);
			opacity: 0;
		}
	}

	@keyframes progress {
		from {
			width: 100%;
		}
		to {
			width: 0%;
		}
	}
</style>
