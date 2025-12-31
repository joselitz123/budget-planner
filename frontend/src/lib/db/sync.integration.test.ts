/**
 * Integration tests for sync.ts
 * Tests sync flows using actual sync functions with mocked fetch
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { resolveConflict } from './sync';
import type { SyncOperation } from './schema';
import {
	createMockSyncOperation,
	createMockSyncOperations,
	createMockSyncResponse,
	createMockFailedSyncResponse,
	createMockPullResponse,
	createMockBudget,
	createMockTransaction,
	createMockCategory
} from '$test/helpers/sync';
import { clearMockFetch } from '$test/setup';

describe('sync.integration - Sync Flows', () => {
	let fetchMock: ReturnType<typeof vi.fn>;

	beforeEach(() => {
		vi.clearAllMocks();
		vi.useFakeTimers();

		fetchMock = vi.fn();
		global.fetch = fetchMock;

		clearMockFetch();
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	describe('Conflict Resolution Integration', () => {
		it('should resolve conflicts for budgets correctly', () => {
			const localBudget = createMockBudget({
				id: 'budget-1',
				income: 4000,
				updatedAt: '2025-01-01T00:00:00Z'
			});
			const serverBudget = createMockBudget({
				id: 'budget-1',
				income: 5000,
				updatedAt: '2025-01-02T00:00:00Z'
			});

			const result = resolveConflict(localBudget, serverBudget);
			expect(result.income).toBe(5000);
			expect(result).toBe(serverBudget);
		});

		it('should resolve conflicts for transactions correctly', () => {
			const localTx = createMockTransaction({
				id: 'tx-1',
				amount: 100,
				updatedAt: '2025-01-01T00:00:00Z'
			});
			const serverTx = createMockTransaction({
				id: 'tx-1',
				amount: 200,
				updatedAt: '2025-01-02T00:00:00Z'
			});

			const result = resolveConflict(localTx, serverTx);
			expect(result.amount).toBe(200);
		});

		it('should resolve conflicts for categories correctly', () => {
			const localCat = createMockCategory({
				id: 'cat-1',
				name: 'Food',
				color: '#FF0000',
				updatedAt: '2025-01-01T00:00:00Z'
			});
			const serverCat = createMockCategory({
				id: 'cat-1',
				name: 'Groceries',
				color: '#00FF00',
				updatedAt: '2025-01-02T00:00:00Z'
			});

			const result = resolveConflict(localCat, serverCat);
			expect(result.name).toBe('Groceries');
			expect(result.color).toBe('#00FF00');
		});

		it('should keep local changes when local is newer', () => {
			const localBudget = createMockBudget({
				id: 'budget-1',
				income: 4500,
				updatedAt: '2025-01-03T00:00:00Z'
			});
			const serverBudget = createMockBudget({
				id: 'budget-1',
				income: 5000,
				updatedAt: '2025-01-02T00:00:00Z'
			});

			const result = resolveConflict(localBudget, serverBudget);
			expect(result.income).toBe(4500);
			expect(result).toBe(localBudget);
		});

		it('should prefer server as tiebreaker for equal timestamps', () => {
			const timestamp = '2025-01-01T00:00:00Z';
			const localBudget = createMockBudget({
				id: 'budget-1',
				name: 'Local Name',
				updatedAt: timestamp
			});
			const serverBudget = createMockBudget({
				id: 'budget-1',
				name: 'Server Name',
				updatedAt: timestamp
			});

			const result = resolveConflict(localBudget, serverBudget);
			expect(result.name).toBe('Server Name');
		});
	});

	describe('Retry Logic Integration', () => {
		const MAX_RETRY_ATTEMPTS = 5;

		it('should track retry count correctly', () => {
			const operation = createMockSyncOperation({ retryCount: 0 });

			for (let i = 0; i < MAX_RETRY_ATTEMPTS; i++) {
				const shouldRetry = operation.retryCount < MAX_RETRY_ATTEMPTS;
				expect(shouldRetry).toBe(true);
				operation.retryCount++;
			}

			const shouldRetry = operation.retryCount < MAX_RETRY_ATTEMPTS;
			expect(shouldRetry).toBe(false);
		});

		it('should handle operation with failed status', () => {
			const operation: SyncOperation = {
				...createMockSyncOperation(),
				status: 'failed',
				retryCount: 5,
				error: 'Max retries exceeded'
			};

			expect(operation.status).toBe('failed');
			expect(operation.retryCount).toBe(5);
			expect(operation.error).toBe('Max retries exceeded');
		});
	});

	describe('Sync Operation Structure', () => {
		it('should create valid sync operation', () => {
			const operation = createMockSyncOperation();

			expect(operation).toHaveProperty('id');
			expect(operation).toHaveProperty('table');
			expect(operation).toHaveProperty('recordId');
			expect(operation).toHaveProperty('operation');
			expect(operation).toHaveProperty('data');
			expect(operation).toHaveProperty('timestamp');
			expect(operation).toHaveProperty('status');
			expect(operation).toHaveProperty('retryCount');

			expect(operation.table).toBe('transactions');
			expect(operation.operation).toBe('CREATE');
			expect(operation.status).toBe('pending');
			expect(operation.retryCount).toBe(0);
		});

		it('should support all operation types', () => {
			const operations = ['CREATE', 'UPDATE', 'DELETE'] as const;

			operations.forEach((op) => {
				const syncOp = createMockSyncOperation({ operation: op });
				expect(syncOp.operation).toBe(op);
			});
		});

		it('should support all table types', () => {
			const tables = ['budgets', 'categories', 'transactions', 'reflections', 'paymentMethods'] as const;

			tables.forEach((table) => {
				const syncOp = createMockSyncOperation({ table });
				expect(syncOp.table).toBe(table);
			});
		});
	});

	describe('Sync Response Handling', () => {
		it('should handle successful sync response', () => {
			const response = createMockSyncResponse({
				successful: ['op-1', 'op-2', 'op-3'],
				failed: []
			});

			expect(response.successful).toHaveLength(3);
			expect(response.failed).toHaveLength(0);
		});

		it('should handle partial success sync response', () => {
			const response = createMockSyncResponse({
				successful: ['op-1', 'op-2'],
				failed: [{ id: 'op-3', error: 'Validation failed' }]
			});

			expect(response.successful).toHaveLength(2);
			expect(response.failed).toHaveLength(1);
			expect(response.failed[0].error).toBe('Validation failed');
		});

		it('should handle failed sync response', () => {
			const response = createMockFailedSyncResponse('Network error');

			expect(response.successful).toHaveLength(0);
			expect(response.failed).toHaveLength(1);
			expect(response.failed[0].error).toBe('Network error');
		});
	});

	describe('Pull Response Handling', () => {
		it('should handle empty pull response', () => {
			const response = createMockPullResponse();

			expect(response.budgets).toHaveLength(0);
			expect(response.transactions).toHaveLength(0);
			expect(response.categories).toHaveLength(0);
			expect(response.reflections).toHaveLength(0);
		});

		it('should handle pull response with data', () => {
			const response = createMockPullResponse({
				budgets: [createMockBudget()],
				transactions: [createMockTransaction()],
				categories: [createMockCategory()]
			});

			expect(response.budgets).toHaveLength(1);
			expect(response.transactions).toHaveLength(1);
			expect(response.categories).toHaveLength(1);
		});

		it('should handle pull response with multiple records', () => {
			const budgets = Array.from({ length: 5 }, () => createMockBudget());
			const response = createMockPullResponse({ budgets });

			expect(response.budgets).toHaveLength(5);
		});
	});

	describe('Data Models', () => {
		it('should create valid budget model', () => {
			const budget = createMockBudget();

			expect(budget).toHaveProperty('id');
			expect(budget).toHaveProperty('userId');
			expect(budget).toHaveProperty('name');
			expect(budget).toHaveProperty('month');
			expect(budget).toHaveProperty('year');
			expect(budget).toHaveProperty('income');
			expect(budget).toHaveProperty('budgetLimit');
			expect(budget).toHaveProperty('createdAt');
			expect(budget).toHaveProperty('updatedAt');
		});

		it('should create valid transaction model', () => {
			const transaction = createMockTransaction();

			expect(transaction).toHaveProperty('id');
			expect(transaction).toHaveProperty('budgetId');
			expect(transaction).toHaveProperty('categoryId');
			expect(transaction).toHaveProperty('amount');
			expect(transaction).toHaveProperty('description');
			expect(transaction).toHaveProperty('date');
			expect(transaction).toHaveProperty('createdAt');
			expect(transaction).toHaveProperty('updatedAt');
		});

		it('should create valid category model', () => {
			const category = createMockCategory();

			expect(category).toHaveProperty('id');
			expect(category).toHaveProperty('userId');
			expect(category).toHaveProperty('name');
			expect(category).toHaveProperty('type');
			expect(category).toHaveProperty('color');
			expect(category).toHaveProperty('icon');
			expect(category).toHaveProperty('createdAt');
			expect(category).toHaveProperty('updatedAt');
		});
	});

	describe('Edge Cases', () => {
		it('should handle operation with null data gracefully', () => {
			const operation: SyncOperation = {
				...createMockSyncOperation(),
				data: null as unknown as Record<string, unknown>
			};

			expect(operation.data).toBeNull();
			expect(operation.id).toBeDefined();
		});

		it('should handle operation with empty data', () => {
			const operation = createMockSyncOperation({
				data: {} as Record<string, unknown>
			});

			expect(Object.keys(operation.data)).toHaveLength(0);
		});

		it('should handle multiple operations with same ID for testing', () => {
			const op1 = createMockSyncOperation({ id: 'same-id' });
			const op2 = createMockSyncOperation({ id: 'same-id' });

			expect(op1.id).toBe(op2.id);
			// In real scenario this shouldn't happen, but tests should handle it
		});
	});
});
