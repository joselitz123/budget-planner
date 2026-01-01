import { syncQueueStore, budgetStore, transactionStore, categoryStore, reflectionStore } from './stores';
import type { SyncOperation } from './schema';
import { get } from 'svelte/store';
import { isOnline, syncStatus, lastSync, pendingSyncCount } from '$lib/stores/offline';
import { showToast } from '$lib/stores/ui';

// Sync configuration
const MAX_RETRY_ATTEMPTS = 5;
const BASE_RETRY_DELAY = 1000; // 1 second
const MAX_RETRY_DELAY = 60000; // 1 minute
const SYNC_INTERVAL = 30000; // 30 seconds

// Track sync interval for cleanup
let syncIntervalId: ReturnType<typeof setInterval> | null = null;

/**
 * Initialize background sync processor
 * Starts periodic sync queue monitoring
 */
export function initBackgroundSync(): void {
	if (syncIntervalId) {
		console.log('[Sync] Background sync already initialized');
		return;
	}

	console.log('[Sync] Initializing background sync processor');
	syncIntervalId = setInterval(async () => {
		if (get(isOnline)) {
			await processSyncQueue();
		}
	}, SYNC_INTERVAL);

	// Initial sync on init
	if (get(isOnline)) {
		processSyncQueue();
	}
}

/**
 * Stop background sync processor
 */
export function stopBackgroundSync(): void {
	if (syncIntervalId) {
		clearInterval(syncIntervalId);
		syncIntervalId = null;
		console.log('[Sync] Background sync stopped');
	}
}

/**
 * Calculate exponential backoff delay
 */
function calculateRetryDelay(retryCount: number): number {
	const delay = Math.min(BASE_RETRY_DELAY * Math.pow(2, retryCount), MAX_RETRY_DELAY);
	// Add jitter to prevent thundering herd
	return delay + Math.random() * 1000;
}

/**
 * Check if operation should be retried
 */
function shouldRetry(operation: SyncOperation): boolean {
	return operation.retryCount < MAX_RETRY_ATTEMPTS;
}

/**
 * Update pending sync count
 */
async function updatePendingCount(): Promise<void> {
	const pendingOps = await syncQueueStore.getPending();
	pendingSyncCount.set(pendingOps.length);
}

/**
 * Add operation to sync queue
 */
export async function addToSyncQueue(operation: Omit<SyncOperation, 'id' | 'timestamp' | 'status' | 'retryCount'>): Promise<void> {
	const syncOp: SyncOperation = {
		...operation,
		id: crypto.randomUUID(),
		timestamp: new Date().toISOString(),
		status: 'pending',
		retryCount: 0
	};

	await syncQueueStore.create(syncOp);
	console.log('[Sync] Added to queue:', syncOp);
	await updatePendingCount();

	// Try to sync immediately if online
	if (get(isOnline)) {
		await processSyncQueue();
	}
}

/**
 * Process sync queue (send pending operations to server)
 * Implements exponential backoff and max retry logic
 */
export async function processSyncQueue(): Promise<void> {
	const pendingOps = await syncQueueStore.getPending();

	// Filter out operations that are waiting for retry delay
	const now = Date.now();
	const readyToSync = pendingOps.filter((op) => {
		if (op.status === 'failed' && op.retryCount > 0) {
			const delay = calculateRetryDelay(op.retryCount - 1);
			const opTime = new Date(op.timestamp).getTime();
			const nextRetry = opTime + delay;
			return now >= nextRetry;
		}
		return true;
	});

	if (readyToSync.length === 0) {
		console.log('[Sync] No operations ready to sync');
		return;
	}

	console.log(`[Sync] Processing ${readyToSync.length} operations (${pendingOps.length} total pending)`);
	syncStatus.set('syncing');

	try {
		// Send to server via sync API
		const apiUrl = import.meta.env.PUBLIC_API_URL || 'http://localhost:8080/api';
		const response = await fetch(`${apiUrl}/sync/push`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ operations: readyToSync })
		});

		if (!response.ok) {
			throw new Error(`Sync failed: ${response.statusText}`);
		}

		const result = await response.json();

		// Handle successful operations
		const successfulOps = result.successful || [];
		for (const opId of successfulOps) {
			await syncQueueStore.delete(opId);
			console.log('[Sync] Operation synced successfully:', opId);
		}

		// Handle failed operations
		const failedOps = result.failed || [];
		for (const failedOp of failedOps) {
			const op = readyToSync.find((o) => o.id === failedOp.id);
			if (op) {
				if (shouldRetry(op)) {
					// Mark for retry with exponential backoff
					await syncQueueStore.update({
						...op,
						status: 'pending',
						error: failedOp.error,
						retryCount: op.retryCount + 1
					});
					const delay = calculateRetryDelay(op.retryCount);
					console.log(`[Sync] Operation failed, will retry in ${Math.round(delay / 1000)}s:`, op.id);
				} else {
					// Max retries reached, mark as permanently failed
					await syncQueueStore.update({
						...op,
						status: 'failed',
						error: failedOp.error || 'Max retry attempts reached',
						retryCount: op.retryCount + 1
					});
					console.error('[Sync] Operation failed permanently:', op.id);
					showToast(`Sync failed: ${failedOp.error || 'Max retries reached'}`, 'error', 5000);
				}
			}
		}

		lastSync.set(new Date().toISOString());
		syncStatus.set('idle');
		await updatePendingCount();
		console.log(`[Sync] Completed: ${successfulOps.length} successful, ${failedOps.length} failed`);
	} catch (error) {
		console.error('[Sync] Error processing queue:', error);

		// Network error - mark operations for retry
		for (const op of readyToSync) {
			if (shouldRetry(op)) {
				await syncQueueStore.update({
					...op,
					status: 'pending',
					error: error instanceof Error ? error.message : 'Network error',
					retryCount: op.retryCount + 1
				});
			} else {
				await syncQueueStore.update({
					...op,
					status: 'failed',
					error: error instanceof Error ? error.message : 'Network error',
					retryCount: op.retryCount + 1
				});
			}
		}

		syncStatus.set('error');
		showToast('Sync failed. Will retry automatically.', 'warning', 3000);
	}
}

/**
 * Pull latest data from server and update IndexedDB
 */
export async function pullFromServer(lastSyncTimestamp?: string): Promise<void> {
	if (!get(isOnline)) {
		console.log('[Sync] Offline, skipping pull');
		return;
	}

	console.log('[Sync] Pulling from server...');
	syncStatus.set('syncing');

	try {
		const apiUrl = import.meta.env.PUBLIC_API_URL || 'http://localhost:8080/api';
		const response = await fetch(`${apiUrl}/sync/pull`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				lastSync: lastSyncTimestamp || get(lastSync)
			})
		});

		if (!response.ok) {
			throw new Error(`Pull failed: ${response.statusText}`);
		}

		const result = await response.json();

		// Update IndexedDB with server data
		if (result.budgets) {
			for (const budget of result.budgets) {
				const existing = await budgetStore.get(budget.id);
				if (existing) {
					// Resolve conflict using timestamp-based resolution
					const resolved = resolveConflict(existing, budget);
					await budgetStore.update(resolved);
				} else {
					await budgetStore.create(budget);
				}
			}
		}

		if (result.transactions) {
			for (const transaction of result.transactions) {
				const existing = await transactionStore.get(transaction.id);
				if (existing) {
					const resolved = resolveConflict(existing, transaction);
					await transactionStore.update(resolved);
				} else {
					await transactionStore.create(transaction);
				}
			}
		}

		if (result.categories) {
			for (const category of result.categories) {
				const existing = await categoryStore.get(category.id);
				if (existing) {
					const resolved = resolveConflict(existing, category);
					await categoryStore.update(resolved);
				} else {
					await categoryStore.create(category);
				}
			}
		}

		if (result.reflections) {
			for (const reflection of result.reflections) {
				const existing = await reflectionStore.get(reflection.id);
				if (existing) {
					const resolved = resolveConflict(existing, reflection);
					await reflectionStore.update(resolved);
				} else {
					await reflectionStore.create(reflection);
				}
			}
		}

		lastSync.set(new Date().toISOString());
		syncStatus.set('idle');
		console.log('[Sync] Pull completed successfully');
		showToast('Data synced from server', 'success', 2000);
	} catch (error) {
		console.error('[Sync] Error pulling from server:', error);
		syncStatus.set('error');
		showToast('Failed to sync with server. Will retry later.', 'warning', 3000);
	}
}

/**
 * Resolve conflict using timestamp-based resolution
 * Latest update wins (more recent timestamp)
 */
export function resolveConflict(localData: any, serverData: any): any {
	const localUpdatedAt = new Date(localData.updatedAt).getTime();
	const serverUpdatedAt = new Date(serverData.updatedAt).getTime();

	if (serverUpdatedAt > localUpdatedAt) {
		console.log(`[Sync] Conflict: Server data is newer (${serverUpdatedAt} > ${localUpdatedAt}), using server version`);
		return serverData;
	} else if (localUpdatedAt > serverUpdatedAt) {
		console.log(`[Sync] Conflict: Local data is newer (${localUpdatedAt} > ${serverUpdatedAt}), using local version`);
		return localData;
	} else {
		// Timestamps are equal - use server as tiebreaker (source of truth)
		console.log('[Sync] Conflict: Timestamps equal, using server as tiebreaker');
		return serverData;
	}
}

/**
 * Manual sync trigger - pushes pending operations and pulls latest data
 */
export async function manualSync(): Promise<void> {
	if (!get(isOnline)) {
		showToast('Cannot sync while offline', 'warning', 3000);
		return;
	}

	console.log('[Sync] Manual sync triggered');
	showToast('Syncing data...', 'info');

	try {
		// First push pending operations
		await processSyncQueue();

		// Then pull latest data from server
		await pullFromServer();

		showToast('Sync completed successfully', 'success', 3000);
	} catch (error) {
		console.error('[Sync] Manual sync failed:', error);
		showToast('Sync failed. Please try again.', 'error', 3000);
	}
}

/**
 * Get sync status summary for UI display
 */
export async function getSyncStatus(): Promise<{
	pending: number;
	lastSync: string | null;
	status: 'idle' | 'syncing' | 'error';
}> {
	const pendingOps = await syncQueueStore.getPending();
	return {
		pending: pendingOps.length,
		lastSync: get(lastSync),
		status: get(syncStatus)
	};
}
