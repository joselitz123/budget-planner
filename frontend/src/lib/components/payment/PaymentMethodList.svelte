<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import type { PaymentMethod } from '$lib/db/schema';
	import { formatCurrency } from '$lib/utils/format';

	export let paymentMethods: PaymentMethod[] = [];
	export let loading = false;

	const dispatch = createEventDispatcher();

	/**
	 * Get icon for payment method type
	 */
	function getTypeIcon(method: PaymentMethod): string {
		switch (method.type) {
			case 'credit_card':
			case 'debit_card':
				return '💳';
			case 'cash':
				return '💵';
			case 'ewallet':
				return '📱';
			default:
				return '💳';
		}
	}

	/**
	 * Get display name for payment method
	 */
	function getDisplayName(method: PaymentMethod): string {
		if (method.type === 'credit_card' || method.type === 'debit_card') {
			const brand = method.brand || 'Card';
			const lastFour = method.lastFour || '----';
			return `${brand} •••• ${lastFour}`;
		} else if (method.type === 'ewallet') {
			return method.brand || 'E-wallet';
		} else {
			return 'Cash';
		}
	}

	/**
	 * Handle edit button click
	 */
	function handleEdit(method: PaymentMethod) {
		dispatch('edit', method);
	}

	/**
	 * Handle set default button click
	 */
	function handleSetDefault(methodId: string) {
		dispatch('setDefault', methodId);
	}

	/**
	 * Handle delete button click
	 */
	function handleDelete(methodId: string) {
		dispatch('delete', methodId);
	}
</script>

<div class="space-y-3">
	{#if loading}
		<!-- Loading skeletons -->
		{#each Array(3) as _}
			<div
				class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-4"
			>
				<div class="flex items-center space-x-3">
					<Skeleton variant="card" className="h-12 w-12 rounded-full" />
					<div class="flex-1 space-y-2">
						<Skeleton variant="text" className="h-4 w-32" />
						<Skeleton variant="text" className="h-3 w-24" />
					</div>
				</div>
			</div>
		{/each}
	{:else}
		{#if paymentMethods.length === 0}
			<!-- Empty state -->
			<div
				class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-8 text-center"
			>
				<div class="text-4xl mb-2">💳</div>
				<p class="font-display text-lg text-primary dark:text-white font-semibold mb-1">
					No payment methods yet
				</p>
				<p class="text-sm text-gray-500 dark:text-gray-400 font-handwriting">
					Add your first payment method to start tracking expenses
				</p>
			</div>
		{:else}
			<!-- Payment method cards -->
			{#each paymentMethods as method}
				<div
					class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-4 hover:shadow-lg transition-shadow {!method.isActive
						? 'opacity-60'
						: ''}"
				>
					<div class="flex items-center justify-between">
						<!-- Left: Icon + Info -->
						<div class="flex items-center space-x-3">
							<!-- Type Icon -->
							<div
								class="h-12 w-12 rounded-full bg-white dark:bg-gray-700 border border-line-light dark:border-line-dark flex items-center justify-center text-2xl shadow-sm flex-shrink-0"
							>
								{getTypeIcon(method)}
							</div>

							<!-- Info -->
							<div>
								<div class="flex items-center space-x-2">
									<h3 class="font-semibold text-primary dark:text-white leading-tight">
										{method.name}
									</h3>
									{#if method.isDefault}
										<Badge variant="default" class="text-xs">Default</Badge>
									{/if}
								</div>

								{#if method.type === 'credit_card' || method.type === 'debit_card'}
									<p class="text-sm text-gray-500 dark:text-gray-400">
										{getDisplayName(method)}
									</p>
								{:else if method.type === 'ewallet'}
									<p class="text-sm text-gray-500 dark:text-gray-400">{getDisplayName(method)}</p>
								{:else}
									<p class="text-sm text-gray-500 dark:text-gray-400">Cash</p>
								{/if}

								<!-- Balance/Limit Info -->
								{#if method.type === 'credit_card'}
									<div class="flex items-center space-x-3 mt-1">
										{#if method.creditLimit}
											<p class="text-xs text-gray-500 dark:text-gray-400">
												Limit: {formatCurrency(method.creditLimit)}
											</p>
										{/if}
										{#if method.currentBalance !== null}
											<p
												class="text-xs {method.currentBalance > 0
													? 'text-red-500 dark:text-red-400'
													: 'text-green-500 dark:text-green-400'}"
											>
												Balance: {formatCurrency(method.currentBalance)}
											</p>
										{/if}
									</div>
								{:else if method.currentBalance !== null}
									<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
										Balance: {formatCurrency(method.currentBalance)}
									</p>
								{/if}

								<!-- Inactive badge -->
								{#if !method.isActive}
									<Badge variant="outline" class="text-xs mt-1">Inactive</Badge>
								{/if}
							</div>
						</div>

						<!-- Right: Actions -->
						<div class="flex items-center space-x-2">
							{#if !method.isDefault}
								<button
									onclick={() => handleSetDefault(method.id)}
									class="text-xs text-primary hover:text-blue-600 dark:hover:text-blue-400 underline decoration-dotted"
									type="button"
								>
									Set Default
								</button>
							{/if}

							<button
								onclick={() => handleEdit(method)}
								class="text-xs text-primary hover:text-blue-600 dark:hover:text-blue-400"
								type="button"
								aria-label="Edit {method.name}"
							>
								Edit
							</button>

							<button
								onclick={() => handleDelete(method.id)}
								class="text-xs text-red-500 hover:text-red-600 dark:hover:text-red-400"
								type="button"
								aria-label="Delete {method.name}"
							>
								Delete
							</button>
						</div>
					</div>
				</div>
			{/each}
		{/if}
	{/if}
</div>
