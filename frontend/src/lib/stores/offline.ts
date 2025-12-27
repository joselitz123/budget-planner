import { writable, derived, get } from 'svelte/store';
import { initDB, pullFromServer, processSyncQueue } from '$lib/db';
import type { Budget, Category, Transaction, Reflection } from '$lib/db/schema';

/**
 * Offline state management
 */

// Online/offline status
export const isOnline = writable<boolean>(
	typeof window !== 'undefined' ? navigator.onLine : true
);

// Sync queue
export const syncQueue = writable<any[]>([]);

// Sync status
export const syncStatus = writable<'idle' | 'syncing' | 'error'>('idle');

// Last sync timestamp
export const lastSync = writable<string | null>(null);

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
