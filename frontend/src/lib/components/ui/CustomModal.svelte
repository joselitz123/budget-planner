<script lang="ts">
	import { cn } from '$lib/utils/cn';
	import { createEventDispatcher, onMount } from 'svelte';

	export let open = false;
	export let title = '';
	export let description = '';
	export let className = '';

	const dispatch = createEventDispatcher();

	// Generate unique IDs for accessibility
	let titleId = `modal-title-${Math.random().toString(36).substr(2, 9)}`;
	let descrId = `modal-desc-${Math.random().toString(36).substr(2, 9)}`;

	// Focus management
	let modalElement: HTMLDivElement;
	let closeButton: HTMLButtonElement;
	let previousActiveElement: Element | null = null;

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

	// Focus trap - keep Tab/Shift+Tab within modal
	function handleTabKey(e: KeyboardEvent) {
		if (e.key !== 'Tab') return;

		const focusableElements = modalElement?.querySelectorAll(
			'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
		);

		if (!focusableElements || focusableElements.length === 0) return;

		const firstElement = focusableElements[0] as HTMLElement;
		const lastElement = focusableElements[focusableElements.length - 1] as HTMLElement;

		if (e.shiftKey) {
			// Shift+Tab: if on first element, move to last
			if (document.activeElement === firstElement) {
				e.preventDefault();
				lastElement.focus();
			}
		} else {
			// Tab: if on last element, move to first
			if (document.activeElement === lastElement) {
				e.preventDefault();
				firstElement.focus();
			}
		}
	}

	// Focus management when modal opens
	function focusModal() {
		// Store the currently focused element to restore later
		previousActiveElement = document.activeElement;

		// Focus the close button (first interactive element)
		setTimeout(() => {
			closeButton?.focus();
		}, 0);
	}

	// Restore focus when modal closes
	function restoreFocus() {
		if (previousActiveElement instanceof HTMLElement) {
			previousActiveElement.focus();
		}
	}

	$: if (open) {
		if (typeof document !== 'undefined') {
			document.addEventListener('keydown', handleKeydown);
			document.addEventListener('keydown', handleTabKey);
			// Prevent body scroll when modal is open
			document.body.style.overflow = 'hidden';
			// Focus the modal
			focusModal();
		}
	} else {
		if (typeof document !== 'undefined') {
			document.removeEventListener('keydown', handleKeydown);
			document.removeEventListener('keydown', handleTabKey);
			document.body.style.overflow = '';
			// Restore focus to previous element
			restoreFocus();
		}
	}
</script>

{#if open}
	<div
		bind:this={modalElement}
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		onclick={handleBackdropClick}
		role="dialog"
		aria-modal="true"
		aria-labelledby={titleId}
		aria-describedby={description ? descrId : undefined}
		tabindex="-1"
	>
		<div
			class={cn(
				'relative w-full max-w-lg bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-6 m-4 max-h-[90vh] overflow-y-auto',
				className
			)}
		>
			<div class="flex items-start justify-between mb-4">
				<div class="flex-1">
					<h2
						id={titleId}
						class="text-2xl font-display font-bold text-primary dark:text-white"
					>
						{title}
					</h2>
					{#if description}
						<p id={descrId} class="text-sm text-gray-500 dark:text-gray-400 mt-1">
							{description}
						</p>
					{/if}
				</div>
				<button
					bind:this={closeButton}
					onclick={close}
					class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 p-1 rounded-md hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors focus:outline-none focus:ring-2 focus:ring-primary"
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
