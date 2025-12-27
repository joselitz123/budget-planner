import { writable, derived, get } from 'svelte/store';
import { transactionStore } from '$lib/db/stores';
import type { Transaction } from '$lib/db/schema';
import { formatCurrency } from '$lib/utils/format';
import { currentBudget } from './budgets';

/**
 * Transaction state management
 */

// All transactions
export const transactions = writable<Transaction[]>([]);

// Filter settings
export const filterCategory = writable<string | null>(null);
export const filterPaid = writable<boolean | null>(null);

// Derived: filtered transactions
export const filteredTransactions = derived(
	[transactions, filterCategory, filterPaid],
	([$transactions, $filterCategory, $filterPaid]) => {
		let filtered = $transactions;

		if ($filterCategory) {
			filtered = filtered.filter((t) => t.categoryId === $filterCategory);
		}

		if ($filterPaid !== null) {
			filtered = filtered.filter((t) => t.paid === $filterPaid);
		}

		return filtered.sort((a, b) => {
			return new Date(b.transactionDate).getTime() - new Date(a.transactionDate).getTime();
		});
	}
);

// Derived: transactions for current budget
export const currentBudgetTransactions = derived(
	[transactions, currentBudget],
	([$transactions, $currentBudget]) => {
		if (!$currentBudget) return [];
		return $transactions.filter((t) => t.budgetId === $currentBudget.id);
	}
);

// Derived: total spent for current budget
export const totalSpent = derived(currentBudgetTransactions, ($transactions) => {
	return $transactions
		.filter((t) => t.transactionType === 'expense')
		.reduce((sum, t) => sum + t.amount, 0);
});

// Derived: total income for current budget
export const totalIncome = derived(currentBudgetTransactions, ($transactions) => {
	return $transactions
		.filter((t) => t.transactionType === 'income')
		.reduce((sum, t) => sum + t.amount, 0);
});

// Derived: unpaid bills
export const unpaidBills = derived(currentBudgetTransactions, ($transactions) => {
	return $transactions.filter(
		(t) => t.isRecurring && !t.paid && t.transactionType === 'expense'
	);
});

/**
 * Load transactions from IndexedDB
 */
export async function loadTransactions(): Promise<void> {
	try {
		const allTransactions = await transactionStore.getAll();
		transactions.set(allTransactions);
		console.log(`[Transactions] Loaded ${allTransactions.length} transactions`);
	} catch (error) {
		console.error('[Transactions] Error loading transactions:', error);
	}
}

/**
 * Add transaction
 */
export async function addTransaction(transaction: Transaction): Promise<void> {
	try {
		await transactionStore.create(transaction);
		transactions.update((current) => [...current, transaction]);
		console.log('[Transactions] Added transaction:', transaction.id);

		// Add to sync queue
		const { addToSyncQueue } = await import('$lib/db/sync');
		await addToSyncQueue({
			table: 'transactions',
			recordId: transaction.id,
			operation: 'CREATE',
			data: transaction
		});
	} catch (error) {
		console.error('[Transactions] Error adding transaction:', error);
		throw error;
	}
}

/**
 * Update transaction
 */
export async function updateTransaction(transaction: Transaction): Promise<void> {
	try {
		await transactionStore.update(transaction);
		transactions.update((current) =>
			current.map((t) => (t.id === transaction.id ? transaction : t))
		);
		console.log('[Transactions] Updated transaction:', transaction.id);

		// Add to sync queue
		const { addToSyncQueue } = await import('$lib/db/sync');
		await addToSyncQueue({
			table: 'transactions',
			recordId: transaction.id,
			operation: 'UPDATE',
			data: transaction
		});
	} catch (error) {
		console.error('[Transactions] Error updating transaction:', error);
		throw error;
	}
}

/**
 * Delete transaction
 */
export async function deleteTransaction(transactionId: string): Promise<void> {
	try {
		await transactionStore.delete(transactionId);
		transactions.update((current) => current.filter((t) => t.id !== transactionId));
		console.log('[Transactions] Deleted transaction:', transactionId);

		// Add to sync queue
		const { addToSyncQueue } = await import('$lib/db/sync');
		await addToSyncQueue({
			table: 'transactions',
			recordId: transactionId,
			operation: 'DELETE',
			data: { id: transactionId }
		});
	} catch (error) {
		console.error('[Transactions] Error deleting transaction:', error);
		throw error;
	}
}
