import { syncQueueStore } from './stores';
import type { SyncOperation } from './schema';
import { get } from 'svelte/store';
import { isOnline, syncStatus, lastSync } from '$lib/stores/offline';

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

	// Try to sync immediately if online
	if (get(isOnline)) {
		await processSyncQueue();
	}
}

/**
 * Process sync queue (send pending operations to server)
 */
export async function processSyncQueue(): Promise<void> {
	const pendingOps = await syncQueueStore.getPending();

	if (pendingOps.length === 0) {
		console.log('[Sync] No pending operations');
		return;
	}

	console.log(`[Sync] Processing ${pendingOps.length} pending operations`);
	syncStatus.set('syncing');

	// Group operations by table for batch processing
	const groupedOps = pendingOps.reduce((acc, op) => {
		if (!acc[op.table]) {
			acc[op.table] = [];
		}
		acc[op.table].push(op);
		return acc;
	}, {} as Record<string, SyncOperation[]>);

	try {
		// Send to server via sync API
		const apiUrl = import.meta.env.VITE_PUBLIC_API_URL || 'http://localhost:8080/api';
		const response = await fetch(`${apiUrl}/sync/push`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ operations: pendingOps })
		});

		if (!response.ok) {
			throw new Error(`Sync failed: ${response.statusText}`);
		}

		const result = await response.json();

		// Remove successful operations from queue
		const successfulOps = result.successful || [];
		const failedOps = result.failed || [];

		for (const opId of successfulOps) {
			await syncQueueStore.delete(opId);
		}

		// Update failed operations
		for (const failedOp of failedOps) {
			const op = pendingOps.find((o) => o.id === failedOp.id);
			if (op) {
				await syncQueueStore.update({
					...op,
					status: 'failed',
					error: failedOp.error,
					retryCount: op.retryCount + 1
				});
			}
		}

		lastSync.set(new Date().toISOString());
		syncStatus.set('idle');
		console.log(`[Sync] Completed: ${successfulOps.length} successful, ${failedOps.length} failed`);
	} catch (error) {
		console.error('[Sync] Error processing queue:', error);

		// Mark operations as failed
		for (const op of pendingOps) {
			await syncQueueStore.update({
				...op,
				status: 'failed',
				error: error instanceof Error ? error.message : 'Unknown error',
				retryCount: op.retryCount + 1
			});
		}

		syncStatus.set('error');
	}
}

/**
 * Pull latest data from server
 */
export async function pullFromServer(lastSyncTimestamp?: string): Promise<void> {
	if (!get(isOnline)) {
		console.log('[Sync] Offline, skipping pull');
		return;
	}

	console.log('[Sync] Pulling from server...');
	syncStatus.set('syncing');

	try {
		const apiUrl = import.meta.env.VITE_PUBLIC_API_URL || 'http://localhost:8080/api';
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
		// TODO: Implement based on API response structure

		lastSync.set(new Date().toISOString());
		syncStatus.set('idle');
		console.log('[Sync] Pull completed successfully');
	} catch (error) {
		console.error('[Sync] Error pulling from server:', error);
		syncStatus.set('error');
	}
}

/**
 * Check for conflicts and resolve (server wins)
 */
export async function resolveConflict(localData: any, serverData: any): Promise<any> {
	// Server wins on conflicts (backend is source of truth)
	const localUpdatedAt = new Date(localData.updated_at).getTime();
	const serverUpdatedAt = new Date(serverData.updated_at).getTime();

	if (serverUpdatedAt > localUpdatedAt) {
		console.log('[Sync] Conflict: Server data is newer, using server version');
		return serverData;
	}

	console.log('[Sync] Conflict: Local data is newer, but server wins');
	return serverData;
}
