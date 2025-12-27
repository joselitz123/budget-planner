import { openDB, type IDBPDatabase } from 'idb';
import type { BudgetDB } from './schema';

/**
 * IndexedDB client for Budget Planner
 *
 * Provides offline-first storage with automatic schema upgrades.
 */

const DB_NAME = 'budget-planner';
const DB_VERSION = 1;

let dbInstance: IDBPDatabase<BudgetDB> | null = null;

/**
 * Initialize and open IndexedDB database
 */
export async function initDB(): Promise<IDBPDatabase<BudgetDB>> {
	if (dbInstance) {
		return dbInstance;
	}

	dbInstance = await openDB<BudgetDB>(DB_NAME, DB_VERSION, {
		upgrade(db) {
			// Create budgets store
			if (!db.objectStoreNames.contains('budgets')) {
				const budgetStore = db.createObjectStore('budgets', { keyPath: 'id' });
				budgetStore.createIndex('by-month', 'month');
				budgetStore.createIndex('by-user', 'userId');
			}

			// Create categories store
			if (!db.objectStoreNames.contains('categories')) {
				const categoryStore = db.createObjectStore('categories', { keyPath: 'id' });
				categoryStore.createIndex('by-user', 'userId');
			}

			// Create transactions store
			if (!db.objectStoreNames.contains('transactions')) {
				const transactionStore = db.createObjectStore('transactions', { keyPath: 'id' });
				transactionStore.createIndex('by-budget', 'budgetId');
				transactionStore.createIndex('by-date', 'transactionDate');
				transactionStore.createIndex('by-category', 'categoryId');
			}

			// Create reflections store
			if (!db.objectStoreNames.contains('reflections')) {
				const reflectionStore = db.createObjectStore('reflections', { keyPath: 'id' });
				reflectionStore.createIndex('by-budget', 'budgetId');
			}

			// Create payment methods store
			if (!db.objectStoreNames.contains('paymentMethods')) {
				const paymentStore = db.createObjectStore('paymentMethods', { keyPath: 'id' });
				paymentStore.createIndex('by-user', 'userId');
			}

			// Create sync queue store
			if (!db.objectStoreNames.contains('syncQueue')) {
				const syncStore = db.createObjectStore('syncQueue', { keyPath: 'id' });
				syncStore.createIndex('by-status', 'status');
			}
		},
		blocked() {
			console.error('[IndexedDB] Database upgrade blocked. Close all tabs and try again.');
		},
		blocking() {
			console.warn('[IndexedDB] Database upgrade blocking. Closing connection...');
			dbInstance?.close();
			dbInstance = null;
		}
	});

	console.log('[IndexedDB] Database initialized successfully');
	return dbInstance;
}

/**
 * Get database instance (initializes if needed)
 */
export async function getDB(): Promise<IDBPDatabase<BudgetDB>> {
	if (!dbInstance) {
		return await initDB();
	}
	return dbInstance;
}

/**
 * Close database connection
 */
export function closeDB(): void {
	if (dbInstance) {
		dbInstance.close();
		dbInstance = null;
		console.log('[IndexedDB] Database connection closed');
	}
}

/**
 * Clear all data (for testing or logout)
 */
export async function clearAllData(): Promise<void> {
	const db = await getDB();
	const tx = db.transaction(['budgets', 'categories', 'transactions', 'reflections', 'paymentMethods', 'syncQueue'], 'readwrite');

	await Promise.all([
		tx.objectStore('budgets').clear(),
		tx.objectStore('categories').clear(),
		tx.objectStore('transactions').clear(),
		tx.objectStore('reflections').clear(),
		tx.objectStore('paymentMethods').clear(),
		tx.objectStore('syncQueue').clear()
	]);

	await tx.done;
	console.log('[IndexedDB] All data cleared');
}
