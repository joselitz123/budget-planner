/**
 * Test helpers for sync testing
 * Provides factories and utilities for creating test data
 */

import type { SyncOperation } from '$lib/db/schema';

/**
 * Create a mock sync operation
 */
export function createMockSyncOperation(overrides: Partial<SyncOperation> = {}): SyncOperation {
	return {
		id: 'test-op-1',
		table: 'transactions',
		recordId: 'test-record-1',
		operation: 'CREATE',
		data: { id: 'test-record-1', amount: 100, description: 'Test' },
		timestamp: new Date().toISOString(),
		status: 'pending',
		retryCount: 0,
		...overrides
	};
}

/**
 * Create multiple mock sync operations
 */
export function createMockSyncOperations(count: number): SyncOperation[] {
	return Array.from({ length: count }, (_, i) =>
		createMockSyncOperation({
			id: `test-op-${i}`,
			recordId: `test-record-${i}`
		})
	);
}

/**
 * Create a mock budget
 */
export function createMockBudget(overrides: Record<string, unknown> = {}) {
	return {
		id: 'budget-1',
		userId: 'user-1',
		name: 'Test Budget',
		month: 1,
		year: 2025,
		income: 5000,
		budgetLimit: 4000,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString(),
		...overrides
	};
}

/**
 * Create a mock transaction
 */
export function createMockTransaction(overrides: Record<string, unknown> = {}) {
	return {
		id: 'transaction-1',
		budgetId: 'budget-1',
		categoryId: 'category-1',
		amount: 100,
		description: 'Test Transaction',
		date: new Date().toISOString(),
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString(),
		...overrides
	};
}

/**
 * Create a mock category
 */
export function createMockCategory(overrides: Record<string, unknown> = {}) {
	return {
		id: 'category-1',
		userId: 'user-1',
		name: 'Test Category',
		type: 'expense',
		color: '#FF0000',
		icon: 'shopping',
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString(),
		...overrides
	};
}

/**
 * Mock sync API response
 */
export function createMockSyncResponse(overrides: Record<string, unknown> = {}) {
	return {
		successful: ['test-op-1'],
		failed: [],
		...overrides
	};
}

/**
 * Mock failed sync response
 */
export function createMockFailedSyncResponse(error: string) {
	return {
		successful: [],
		failed: [{ id: 'test-op-1', error }]
	};
}

/**
 * Mock pull response from server
 */
export function createMockPullResponse(overrides: Record<string, unknown> = {}) {
	return {
		budgets: [],
		transactions: [],
		categories: [],
		reflections: [],
		...overrides
	};
}

/**
 * Wait for async operations
 */
export function waitFor(ms: number): Promise<void> {
	return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * Flush all pending promises
 */
export async function flushPromises(): Promise<void> {
	await waitFor(0);
}

/**
 * Create a mock store for testing
 */
export function createMockStore() {
	const data = new Map<string, unknown>();

	return {
		get: async (id: string) => data.get(id),
		create: async (item: { id: string }) => {
			data.set(item.id, item);
			return item;
		},
		update: async (item: { id: string }) => {
			data.set(item.id, item);
			return item;
		},
		delete: async (id: string) => {
			data.delete(id);
		},
		getAll: async () => Array.from(data.values()),
		clear: async () => data.clear()
	};
}

/**
 * Create a mock sync queue store
 */
export function createMockSyncQueueStore() {
	const operations: Map<string, SyncOperation> = new Map();

	return {
		get: async (id: string) => operations.get(id),
		create: async (op: SyncOperation) => {
			operations.set(op.id, op);
			return op;
		},
		update: async (op: SyncOperation) => {
			operations.set(op.id, op);
			return op;
		},
		delete: async (id: string) => {
			operations.delete(id);
		},
		getPending: async () => {
			return Array.from(operations.values()).filter((op) => op.status === 'pending');
		},
		getAll: async () => Array.from(operations.values()),
		clear: async () => operations.clear()
	};
}
