import { apiClient } from './client';

/**
 * Sync API types (matching backend responses)
 */

export interface SyncOperation {
	table: string;
	recordId: string;
	operation: 'create' | 'update' | 'delete';
	localData: Record<string, any>;
	serverData?: Record<string, any>;
}

export interface PushRequest {
	operations: SyncOperation[];
}

export interface PushResponse {
	results: Array<{
		success: boolean;
		recordId: string;
		error?: string;
	}>;
	syncedAt: string;
}

export interface PullRequest {
	lastSyncTime: string; // ISO 8601 timestamp
}

export interface PullResponse {
	hasMore: boolean;
	lastSyncTime: string;
	changes: Record<string, any[]>; // table name -> records
}

export interface SyncStatus {
	lastSyncTime: string;
	pendingOperations: number;
}

/**
 * Sync API client
 */
export class SyncApi {
	/**
	 * Push local changes to the server
	 * POST /api/sync/push
	 */
	async pushSync(operations: SyncOperation[]): Promise<PushResponse> {
		const request: PushRequest = { operations };
		return apiClient.post<PushResponse>('/sync/push', request);
	}

	/**
	 * Pull server changes down to the client
	 * POST /api/sync/pull
	 */
	async pullSync(lastSyncTime?: string): Promise<PullResponse> {
		const request: PullRequest = {
			lastSyncTime: lastSyncTime || ''
		};
		return apiClient.post<PullResponse>('/sync/pull', request);
	}

	/**
	 * Get sync status
	 * GET /api/sync/status
	 */
	async getSyncStatus(): Promise<SyncStatus> {
		return apiClient.get<SyncStatus>('/sync/status');
	}

	/**
	 * Resolve a sync conflict
	 * POST /api/sync/resolve-conflict
	 */
	async resolveConflict(conflictId: string, resolution: 'local' | 'server'): Promise<void> {
		return apiClient.post<void>('/sync/resolve-conflict', {
			conflictId,
			resolution
		});
	}

	/**
	 * Get conflict history
	 * GET /api/sync/conflict-history
	 */
	async getConflictHistory(): Promise<any[]> {
		return apiClient.get<any[]>('/sync/conflict-history');
	}
}

// Export singleton instance
export const syncApi = new SyncApi();
