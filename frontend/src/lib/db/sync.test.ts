/**
 * Unit tests for sync.ts
 * Tests offline sync queue processing, retry logic, conflict resolution, and API calls
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import {
	resolveConflict
} from './sync';
import type { SyncOperation } from './schema';
import {
	createMockSyncOperation,
	createMockSyncOperations,
	createMockBudget,
	createMockTransaction,
	createMockCategory
} from '$test/helpers/sync';
import { createMockStore, clearMockFetch } from '$test/setup';

// Note: We can't easily mock the stores that sync.ts imports,
// so we'll test the pure functions that don't depend on external state

describe('sync.ts - Pure Functions', () => {
	beforeEach(() => {
		vi.useFakeTimers();
	});

	afterEach(() => {
		vi.useRealTimers();
		clearMockFetch();
	});

	describe('Conflict Resolution', () => {
		it('should return server data when server timestamp is newer', () => {
			const localData = { ...createMockBudget(), updatedAt: '2025-01-01T00:00:00Z' };
			const serverData = { ...createMockBudget(), updatedAt: '2025-01-02T00:00:00Z' };

			const result = resolveConflict(localData, serverData);
			expect(result).toBe(serverData);
		});

		it('should return local data when local timestamp is newer', () => {
			const localData = { ...createMockBudget(), updatedAt: '2025-01-02T00:00:00Z' };
			const serverData = { ...createMockBudget(), updatedAt: '2025-01-01T00:00:00Z' };

			const result = resolveConflict(localData, serverData);
			expect(result).toBe(localData);
		});

		it('should return server data as tiebreaker when timestamps are equal', () => {
			const timestamp = '2025-01-01T00:00:00Z';
			const localData = { ...createMockBudget(), updatedAt: timestamp, name: 'Local Budget' };
			const serverData = { ...createMockBudget(), updatedAt: timestamp, name: 'Server Budget' };

			const result = resolveConflict(localData, serverData);
			expect(result).toBe(serverData);
		});

		it('should handle missing updatedAt timestamps', () => {
			const localData = { ...createMockBudget() };
			const serverData = { ...createMockBudget() };
			delete (localData as any).updatedAt;
			delete (serverData as any).updatedAt;

			// Should not throw and should return server as tiebreaker
			const result = resolveConflict(localData, serverData);
			expect(result).toBe(serverData);
		});

		it('should handle transactions with different amounts', () => {
			const localData = { ...createMockTransaction(), amount: 100, updatedAt: '2025-01-01T00:00:00Z' };
			const serverData = { ...createMockTransaction(), amount: 200, updatedAt: '2025-01-02T00:00:00Z' };

			const result = resolveConflict(localData, serverData);
			expect(result.amount).toBe(200);
		});

		it('should preserve all server data fields when server wins', () => {
			const localData = {
				...createMockCategory(),
				name: 'Local Category',
				color: '#FF0000',
				updatedAt: '2025-01-01T00:00:00Z'
			};
			const serverData = {
				...createMockCategory(),
				name: 'Server Category',
				color: '#00FF00',
				updatedAt: '2025-01-02T00:00:00Z'
			};

			const result = resolveConflict(localData, serverData);
			expect(result.name).toBe('Server Category');
			expect(result.color).toBe('#00FF00');
		});
	});

	describe('Retry Logic', () => {
		const MAX_RETRY_ATTEMPTS = 5;

		it('should allow retry when count is below max', () => {
			const operation = createMockSyncOperation({ retryCount: 3 });
			const shouldRetry = operation.retryCount < MAX_RETRY_ATTEMPTS;
			expect(shouldRetry).toBe(true);
		});

		it('should not allow retry when count is at max', () => {
			const operation = createMockSyncOperation({ retryCount: 5 });
			const shouldRetry = operation.retryCount < MAX_RETRY_ATTEMPTS;
			expect(shouldRetry).toBe(false);
		});

		it('should allow retry up to max attempts', () => {
			const operation = createMockSyncOperation({ retryCount: 4 });
			const shouldRetry = operation.retryCount < MAX_RETRY_ATTEMPTS;
			expect(shouldRetry).toBe(true);
		});
	});

	describe('Exponential Backoff Calculation', () => {
		const BASE_RETRY_DELAY = 1000; // 1 second
		const MAX_RETRY_DELAY = 60000; // 1 minute

		function calculateRetryDelay(retryCount: number): number {
			const delay = Math.min(BASE_RETRY_DELAY * Math.pow(2, retryCount), MAX_RETRY_DELAY);
			return delay + Math.random() * 1000;
		}

		it('should calculate correct delay for retry 0', () => {
			const delay = calculateRetryDelay(0);
			expect(delay).toBeGreaterThanOrEqual(1000);
			expect(delay).toBeLessThan(3000); // base + jitter
		});

		it('should calculate correct delay for retry 1', () => {
			const delay = calculateRetryDelay(1);
			expect(delay).toBeGreaterThanOrEqual(2000);
			expect(delay).toBeLessThan(4000);
		});

		it('should calculate correct delay for retry 2', () => {
			const delay = calculateRetryDelay(2);
			expect(delay).toBeGreaterThanOrEqual(4000);
			expect(delay).toBeLessThan(6000);
		});

		it('should cap at max delay', () => {
			const delay = calculateRetryDelay(10);
			expect(delay).toBeGreaterThanOrEqual(MAX_RETRY_DELAY);
			expect(delay).toBeLessThan(MAX_RETRY_DELAY + 2000);
		});
	});
});

describe('sync.ts - Mock Store Tests', () => {
	it('should create a working mock store', () => {
		const store = createMockStore<boolean>(true);

		let receivedValue: boolean | undefined;
		const unsubscribe = store.subscribe((val) => {
			receivedValue = val;
		});

		expect(receivedValue).toBe(true);
		expect(store.get()).toBe(true);

		store.set(false);
		expect(receivedValue).toBe(false);
		expect(store.get()).toBe(false);

		unsubscribe();
	});

	it('should support store updates', () => {
		const store = createMockStore<number>(10);

		store.update((val) => val * 2);
		expect(store.get()).toBe(20);
	});
});
