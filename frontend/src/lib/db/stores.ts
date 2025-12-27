import type { Budget, Category, Transaction, Reflection, PaymentMethod, SyncOperation } from './schema';
import { getDB } from './client';
import type { IDBPDatabase } from 'idb';
import type { BudgetDB } from './schema';

/**
 * Budget CRUD operations
 */
export const budgetStore = {
	async getAll(): Promise<Budget[]> {
		const db = await getDB();
		return db.getAll('budgets');
	},

	async get(id: string): Promise<Budget | undefined> {
		const db = await getDB();
		return db.get('budgets', id);
	},

	async getByMonth(month: string): Promise<Budget[]> {
		const db = await getDB();
		return db.getAllFromIndex('budgets', 'by-month', month);
	},

	async create(budget: Budget): Promise<string> {
		const db = await getDB();
		await db.put('budgets', budget);
		return budget.id;
	},

	async update(budget: Budget): Promise<void> {
		const db = await getDB();
		await db.put('budgets', budget);
	},

	async delete(id: string): Promise<void> {
		const db = await getDB();
		await db.delete('budgets', id);
	}
};

/**
 * Category CRUD operations
 */
export const categoryStore = {
	async getAll(): Promise<Category[]> {
		const db = await getDB();
		return db.getAll('categories');
	},

	async get(id: string): Promise<Category | undefined> {
		const db = await getDB();
		return db.get('categories', id);
	},

	async getByUser(userId: string): Promise<Category[]> {
		const db = await getDB();
		return db.getAllFromIndex('categories', 'by-user', userId);
	},

	async create(category: Category): Promise<string> {
		const db = await getDB();
		await db.put('categories', category);
		return category.id;
	},

	async update(category: Category): Promise<void> {
		const db = await getDB();
		await db.put('categories', category);
	},

	async delete(id: string): Promise<void> {
		const db = await getDB();
		await db.delete('categories', id);
	}
};

/**
 * Transaction CRUD operations
 */
export const transactionStore = {
	async getAll(): Promise<Transaction[]> {
		const db = await getDB();
		return db.getAll('transactions');
	},

	async get(id: string): Promise<Transaction | undefined> {
		const db = await getDB();
		return db.get('transactions', id);
	},

	async getByBudget(budgetId: string): Promise<Transaction[]> {
		const db = await getDB();
		return db.getAllFromIndex('transactions', 'by-budget', budgetId);
	},

	async getByDateRange(startDate: string, endDate: string): Promise<Transaction[]> {
		const db = await getDB();
		const tx = db.transaction('transactions', 'readonly');
		const index = tx.store.index('by-date');

		const transactions: Transaction[] = [];
		let cursor = await index.openCursor(IDBKeyRange.lowerBound(startDate));

		while (cursor) {
			if (cursor.value.transactionDate > endDate) break;
			transactions.push(cursor.value);
			cursor = await cursor.continue();
		}

		return transactions;
	},

	async create(transaction: Transaction): Promise<string> {
		const db = await getDB();
		await db.put('transactions', transaction);
		return transaction.id;
	},

	async update(transaction: Transaction): Promise<void> {
		const db = await getDB();
		await db.put('transactions', transaction);
	},

	async delete(id: string): Promise<void> {
		const db = await getDB();
		await db.delete('transactions', id);
	}
};

/**
 * Reflection CRUD operations
 */
export const reflectionStore = {
	async getAll(): Promise<Reflection[]> {
		const db = await getDB();
		return db.getAll('reflections');
	},

	async get(id: string): Promise<Reflection | undefined> {
		const db = await getDB();
		return db.get('reflections', id);
	},

	async getByBudget(budgetId: string): Promise<Reflection[]> {
		const db = await getDB();
		return db.getAllFromIndex('reflections', 'by-budget', budgetId);
	},

	async create(reflection: Reflection): Promise<string> {
		const db = await getDB();
		await db.put('reflections', reflection);
		return reflection.id;
	},

	async update(reflection: Reflection): Promise<void> {
		const db = await getDB();
		await db.put('reflections', reflection);
	},

	async delete(id: string): Promise<void> {
		const db = await getDB();
		await db.delete('reflections', id);
	}
};

/**
 * PaymentMethod CRUD operations
 */
export const paymentMethodStore = {
	async getAll(): Promise<PaymentMethod[]> {
		const db = await getDB();
		return db.getAll('paymentMethods');
	},

	async get(id: string): Promise<PaymentMethod | undefined> {
		const db = await getDB();
		return db.get('paymentMethods', id);
	},

	async create(method: PaymentMethod): Promise<string> {
		const db = await getDB();
		await db.put('paymentMethods', method);
		return method.id;
	},

	async update(method: PaymentMethod): Promise<void> {
		const db = await getDB();
		await db.put('paymentMethods', method);
	},

	async delete(id: string): Promise<void> {
		const db = await getDB();
		await db.delete('paymentMethods', id);
	}
};

/**
 * SyncQueue operations
 */
export const syncQueueStore = {
	async getAll(): Promise<SyncOperation[]> {
		const db = await getDB();
		return db.getAll('syncQueue');
	},

	async getPending(): Promise<SyncOperation[]> {
		const db = await getDB();
		return db.getAllFromIndex('syncQueue', 'by-status', 'pending');
	},

	async create(operation: SyncOperation): Promise<string> {
		const db = await getDB();
		await db.put('syncQueue', operation);
		return operation.id;
	},

	async update(operation: SyncOperation): Promise<void> {
		const db = await getDB();
		await db.put('syncQueue', operation);
	},

	async delete(id: string): Promise<void> {
		const db = await getDB();
		await db.delete('syncQueue', id);
	},

	async clear(): Promise<void> {
		const db = await getDB();
		await db.clear('syncQueue');
	}
};
