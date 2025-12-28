import { writable, derived, get } from 'svelte/store';
import { transactionStore } from '$lib/db/stores';
import type { Transaction } from '$lib/db/schema';
import { formatCurrency } from '$lib/utils/format';
import { currentBudget } from './budgets';
import { transactionsApi } from '$lib/api/transactions';

/**
 * Transaction state management
 * Supports both API and IndexedDB (offline-first)
 */

// Flag to enable/disable API integration
// Set to false to use only IndexedDB (offline mode)
let USE_API = true;

// All transactions
export const transactions = writable<Transaction[]>([]);

// Transaction loading state
export const transactionsLoading = writable<boolean>(false);

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
 * Load transactions from API (with IndexedDB fallback)
 */
export async function loadTransactions(): Promise<void> {
	transactionsLoading.set(true);

	try {
		// Try API first if enabled
		if (USE_API) {
			try {
				const apiTransactions = await transactionsApi.getTransactions();

				// Update store
				transactions.set(apiTransactions);

				// Update IndexedDB for offline access
				for (const transaction of apiTransactions) {
					await transactionStore.update(transaction);
				}

				console.log(`[Transactions] Loaded ${apiTransactions.length} transactions from API`);
				return;
			} catch (error) {
				console.warn('[Transactions] API load failed, falling back to IndexedDB:', error);
				// Fall through to IndexedDB loading
			}
		}

		// Fallback to IndexedDB
		const allTransactions = await transactionStore.getAll();
		transactions.set(allTransactions);
		console.log(`[Transactions] Loaded ${allTransactions.length} transactions from IndexedDB`);
	} catch (error) {
		console.error('[Transactions] Error loading transactions:', error);
	} finally {
		transactionsLoading.set(false);
	}
}

/**
 * Add transaction
 */
export async function addTransaction(transaction: Transaction): Promise<void> {
	// Try API first if enabled
	if (USE_API) {
		try {
			const created = await transactionsApi.createTransaction(transaction);

			// Update IndexedDB for offline access
			await transactionStore.create(created);

			// Update store
			transactions.update((current) => [...current, created]);

			console.log('[Transactions] Added transaction via API:', created.id);
			return;
		} catch (error) {
			console.warn('[Transactions] API create failed, using local fallback:', error);
			// Fall through to local creation
		}
	}

	// Fallback: Create locally (will sync later)
	try {
		await transactionStore.create(transaction);
		transactions.update((current) => [...current, transaction]);
		console.log('[Transactions] Added transaction locally:', transaction.id);

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
	// Try API first if enabled
	if (USE_API) {
		try {
			const updated = await transactionsApi.updateTransaction(transaction.id, transaction);

			// Update IndexedDB for offline access
			await transactionStore.update(updated);

			// Update store
			transactions.update((current) =>
				current.map((t) => (t.id === transaction.id ? updated : t))
			);

			console.log('[Transactions] Updated transaction via API:', transaction.id);
			return;
		} catch (error) {
			console.warn('[Transactions] API update failed, using local fallback:', error);
			// Fall through to local update
		}
	}

	// Fallback: Update locally (will sync later)
	try {
		await transactionStore.update(transaction);
		transactions.update((current) =>
			current.map((t) => (t.id === transaction.id ? transaction : t))
		);
		console.log('[Transactions] Updated transaction locally:', transaction.id);

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
	// Try API first if enabled
	if (USE_API) {
		try {
			await transactionsApi.deleteTransaction(transactionId);

			// Update IndexedDB for offline access
			await transactionStore.delete(transactionId);

			// Update store
			transactions.update((current) => current.filter((t) => t.id !== transactionId));

			console.log('[Transactions] Deleted transaction via API:', transactionId);
			return;
		} catch (error) {
			console.warn('[Transactions] API delete failed, using local fallback:', error);
			// Fall through to local deletion
		}
	}

	// Fallback: Delete locally (will sync later)
	try {
		await transactionStore.delete(transactionId);
		transactions.update((current) => current.filter((t) => t.id !== transactionId));
		console.log('[Transactions] Deleted transaction locally:', transactionId);

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
