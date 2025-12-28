<script lang="ts">
	import { formatCurrency, formatShortDate } from '$lib/utils/format';
	import { unpaidBills } from '$lib/stores';
	import { getCategoryById } from '$lib/stores/categories';
	import { updateTransaction, transactionsLoading } from '$lib/stores/transactions';
	import { Spinner } from '$lib/components/ui/spinner';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import type { Transaction } from '$lib/db/schema';

	async function markAsPaid(bill: Transaction) {
		try {
			// Create updated transaction with paid = true
			const updated = { ...bill, paid: true };
			await updateTransaction(updated);
		} catch (error) {
			console.error('Failed to mark bill as paid:', error);
			// Could show toast notification here
		}
	}
</script>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="mb-6">
		<h2 class="text-3xl font-display font-bold text-primary dark:text-white mb-2">
			Bill Payment
		</h2>
		<p class="text-sm font-handwriting text-gray-500 dark:text-gray-400 text-xl">
			Track your bills and due dates
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
				<p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Total Due</p>
				<h3 class="text-2xl font-bold font-display text-primary dark:text-white">
					{formatCurrency(
						$unpaidBills.reduce((sum: number, bill: any) => sum + bill.amount, 0)
					)}
				</h3>
				<p class="text-xs text-red-500 mt-1">
					{$unpaidBills.length} unpaid
				</p>
			</div>
			<div
				class="bg-paper-light dark:bg-paper-dark p-4 rounded-xl shadow-paper border border-line-light dark:border-line-dark"
			>
				<p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Paid This Month</p>
				<h3 class="text-2xl font-bold font-display text-green-600 dark:text-green-400">
					{formatCurrency(
						$unpaidBills
							.filter((b: any) => b.paid)
							.reduce((sum: number, bill: any) => sum + bill.amount, 0)
					)}
				</h3>
			</div>
		{/if}
	</div>

	<!-- Bill List -->
	<div
		class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark overflow-hidden relative"
	>
		<!-- Spiral Binding Decoration -->
		<div
			class="absolute left-0 top-0 bottom-0 w-6 bg-gray-100 dark:bg-gray-800 border-r border-line-light dark:border-line-dark flex flex-col items-center py-4 space-y-4 z-10"
		>
			<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
			<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
			<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
		</div>

		<div class="pl-8">
			{#if $transactionsLoading}
				<div class="p-8 text-center">
					<Spinner size="xl" color="primary" />
					<p class="mt-4 text-gray-500 dark:text-gray-400 font-handwriting text-xl">
						Loading bills...
					</p>
				</div>
			{:else}
				<div class="divide-y divide-line-light dark:divide-line-dark">
					{#each $unpaidBills as bill}
				<div
					class="p-4 flex items-center justify-between group {bill.paid
						? 'opacity-60'
						: ''}"
				>
					<div class="flex items-center space-x-3">
						<div
							class="h-10 w-10 rounded-full bg-white dark:bg-gray-700 border border-line-light dark:border-line-dark flex items-center justify-center text-xl shadow-sm"
						>
							{getCategoryById(bill.categoryId)?.icon || '📄'}
						</div>
						<div>
							<p class="font-semibold text-primary dark:text-white leading-tight">
								{bill.description || 'Bill'}
							</p>
							{#if bill.dueDate}
								<p
									class="text-xs {new Date(bill.dueDate) < new Date()
										? 'text-red-500 font-medium'
										: 'text-gray-500 dark:text-gray-400'}"
								>
									Due: {formatShortDate(bill.dueDate)}
								</p>
							{/if}
						</div>
					</div>
					<div class="text-right">
						<p class="font-bold text-primary dark:text-white">
							{formatCurrency(bill.amount)}
						</p>
						{#if !bill.paid}
							<button
								class="mt-1 text-xs text-primary underline decoration-dotted hover:text-blue-600 dark:hover:text-blue-400"
								type="button"
								onclick={() => markAsPaid(bill)}
							>
								Mark Paid
							</button>
						{:else}
							<span
								class="inline-block px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider text-green-600 bg-green-100 dark:bg-green-900/40 dark:text-green-300 rounded-full"
							>
								Paid
							</span>
						{/if}
					</div>
				</div>
			{:else}
				<div class="p-8 text-center text-gray-500 dark:text-gray-400">
					<span class="material-icons-outlined text-6xl text-gray-300 dark:text-gray-600 mb-4">
						receipt_long
					</span>
					<p class="font-handwriting text-xl">No bills yet</p>
					<p class="text-sm mt-2">Your recurring bills will appear here</p>
				</div>
			{/each}
				</div>
			{/if}
		</div>
	</div>
</div>
