import { writable, derived, get } from 'svelte/store';
import { processSyncQueue } from '$lib/db/sync';
import type { Budget, Category, Transaction, Reflection } from '$lib/db/schema';

/**
 * Offline state management
 */

// Online/offline status
export const isOnline = writable<boolean>(
	typeof window !== 'undefined' ? navigator.onLine : true
);

// Sync status
export const syncStatus = writable<'idle' | 'syncing' | 'error'>('idle');

// Last sync timestamp
export const lastSync = writable<string | null>(null);

// Pending sync operations count
export const pendingSyncCount = writable<number>(0);

// Initialize offline event listeners
if (typeof window !== 'undefined') {
	window.addEventListener('online', () => {
		console.log('[Offline] Connection restored');
		isOnline.set(true);
		// Trigger sync when back online
		processSyncQueue();
	});

	window.addEventListener('offline', () => {
		console.log('[Offline] Connection lost');
		isOnline.set(false);
	});
}

// Derived store for offline indicator
export const showOfflineIndicator = derived(isOnline, ($isOnline) => !$isOnline);

// Derived store for sync status indicator (combines status with pending count)
export const syncIndicator = derived(
	[syncStatus, pendingSyncCount, isOnline],
	([$syncStatus, $pendingSyncCount, $isOnline]) => {
		if (!$isOnline) {
			return { status: 'offline' as const, count: 0, label: 'Offline' };
		}
		if ($syncStatus === 'syncing') {
			return { status: 'syncing' as const, count: $pendingSyncCount, label: 'Syncing...' };
		}
		if ($syncStatus === 'error') {
			return { status: 'error' as const, count: $pendingSyncCount, label: 'Sync error' };
		}
		if ($pendingSyncCount > 0) {
			return { status: 'pending' as const, count: $pendingSyncCount, label: `${$pendingSyncCount} pending` };
		}
		return { status: 'synced' as const, count: 0, label: 'Synced' };
	}
);
