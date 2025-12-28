<script lang="ts">
	import { onMount } from 'svelte';
	import { formatCurrency, formatShortDate, getCategoryColor } from '$lib/utils/format';
	import { filteredTransactions, totalSpent } from '$lib/stores';
	import { getCategoryById } from '$lib/stores/categories';
	import { addTransaction, transactionsLoading } from '$lib/stores/transactions';
	import { showToast } from '$lib/stores/ui';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { Spinner } from '$lib/components/ui/spinner';
	import AddExpenseModal from './AddExpenseModal.svelte';
	import type { Transaction } from '$lib/db/schema';

	let modalOpen = false;
	let userId = 'temp-user'; // TODO: Get from Clerk auth

	async function handleAddTransaction(transaction: Transaction) {
		try {
			await addTransaction(transaction);
			showToast('Expense added successfully!', 'success');
			modalOpen = false;
		} catch (error) {
			console.error('[Transactions] Error adding transaction:', error);
			showToast('Failed to add expense', 'error');
		}
	}
</script>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="mb-6">
		<h2 class="text-3xl font-display font-bold text-primary dark:text-white mb-2">
			Expense Tracker
		</h2>
		<p class="text-sm font-handwriting text-gray-500 dark:text-gray-400 text-xl">
			Track your spending
		</p>
	</div>

	<!-- Summary Cards -->
	<div class="grid grid-cols-2 gap-4 mb-6">
		{#if $transactionsLoading}
			<div
				class="bg-paper-light dark:bg-paper-dark p-4 rounded-xl shadow-paper border border-line-light dark:border-line-dark"
			>
				<Skeleton variant="text" className="h-4 w-20 mb-2" />
				<Skeleton variant="card" className="h-8 w-full" />
			</div>
			<div
				class="bg-paper-light dark:bg-paper-dark p-4 rounded-xl shadow-paper border border-line-light dark:border-line-dark"
			>
				<Skeleton variant="text" className="h-4 w-20 mb-2" />
				<Skeleton variant="card" className="h-8 w-full" />
			</div>
		{:else}
			<div
				class="bg-paper-light dark:bg-paper-dark p-4 rounded-xl shadow-paper border border-line-light dark:border-line-dark"
			>
				<p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Total Spent</p>
				<p class="text-2xl font-bold font-display text-primary dark:text-white">
					{formatCurrency($totalSpent)}
				</p>
			</div>
			<div
				class="bg-paper-light dark:bg-paper-dark p-4 rounded-xl shadow-paper border border-line-light dark:border-line-dark"
			>
				<p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Transactions</p>
				<p class="text-2xl font-bold font-display text-primary dark:text-white">
					{$filteredTransactions.length}
				</p>
			</div>
		{/if}
	</div>

	<!-- Transaction List -->
	<div
		class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper overflow-hidden relative border border-line-light dark:border-line-dark"
	>
		<!-- Spiral Binding Decoration -->
		<div
			class="absolute left-0 top-0 bottom-0 w-6 bg-gray-100 dark:bg-gray-800 border-r border-line-light dark:border-line-dark flex flex-col items-center py-4 space-y-4 z-10"
		>
			<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
			<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
			<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
			<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
			<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
		</div>

		<!-- Table -->
		<div class="pl-8 overflow-x-auto">
			{#if $transactionsLoading}
				<div class="p-8 text-center">
					<Spinner size="xl" color="primary" />
					<p class="mt-4 text-gray-500 dark:text-gray-400 font-handwriting text-lg">
						Loading transactions...
					</p>
				</div>
			{:else}
				<table class="w-full text-left">
					<thead>
						<tr class="border-b-2 border-primary">
							<th class="py-3 pl-4 pr-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
								Date
							</th>
							<th class="py-3 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
								Description
							</th>
							<th class="py-3 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
								Category
							</th>
							<th class="py-3 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 text-right">
								Amount
							</th>
							<th class="py-3 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 text-center">
								Paid
							</th>
						</tr>
					</thead>
					<tbody class="text-sm font-hand text-lg">
						{#each $filteredTransactions as transaction}
							<tr
								class="border-b border-line-light dark:border-line-dark hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors"
							>
								<td class="pl-4 py-2 text-gray-600 dark:text-gray-300 whitespace-nowrap">
									{formatShortDate(transaction.transactionDate)}
								</td>
								<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">
									{transaction.description || 'No description'}
								</td>
								<td class="px-2 py-2">
									{#if getCategoryById(transaction.categoryId)}
										{#if getCategoryById(transaction.categoryId)}
											<span class="text-lg" title={getCategoryById(transaction.categoryId)?.name}>
												{getCategoryById(transaction.categoryId)?.icon}
											</span>
										{/if}
									{/if}
								</td>
								<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">
									{formatCurrency(transaction.amount)}
								</td>
								<td class="px-2 py-2 text-center">
									{#if transaction.paid}
										<span class="material-icons-outlined text-green-500 text-sm">check_circle</span>
									{:else}
										<span class="material-icons-outlined text-gray-300 dark:text-gray-600 text-sm"
											>radio_button_unchecked</span
										>
									{/if}
								</td>
							</tr>
						{:else}
							<tr>
								<td colspan="5" class="px-4 py-8 text-center text-gray-500 dark:text-gray-400">
									<p class="font-handwriting text-xl">No transactions yet</p>
									<p class="text-sm mt-2">Add your first expense to get started</p>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			{/if}
		</div>
	</div>

	<!-- Add Expense FAB -->
	<button
		onclick={() => (modalOpen = true)}
		class="fixed bottom-24 right-6 md:bottom-8 bg-primary text-white p-4 rounded-full shadow-lg hover:shadow-xl hover:bg-gray-700 transition-all transform hover:-translate-y-1 z-30"
		aria-label="Add expense"
	>
		<span class="material-icons-outlined text-2xl">add</span>
	</button>

	<!-- Add Expense Modal -->
	<AddExpenseModal
		bind:open={modalOpen}
		{userId}
		on:submit={(e) => handleAddTransaction(e.detail)}
	/>
</div>
