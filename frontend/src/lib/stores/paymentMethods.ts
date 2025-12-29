/**
 * Payment Methods Store
 *
 * State management for payment methods with offline-first architecture.
 * Follows the same pattern as budgets.ts (API + IndexedDB hybrid).
 */

import { writable, derived, get } from 'svelte/store';
import { paymentMethodStore } from '$lib/db/stores';
import type { PaymentMethod } from '$lib/db/schema';
import { paymentMethodsApi } from '$lib/api/paymentMethods';
import { showToast } from '$lib/stores/ui';

// Flag to enable/disable API integration
// Set to false to use only IndexedDB (offline mode)
let USE_API = true;

// All payment methods
export const paymentMethods = writable<PaymentMethod[]>([]);

// Loading state
export const paymentMethodsLoading = writable<boolean>(false);

// Derived: default payment method
export const defaultPaymentMethod = derived(
	paymentMethods,
	($paymentMethods) => $paymentMethods.find((m) => m.isDefault) || null
);

// Derived: active payment methods
export const activePaymentMethods = derived(
	paymentMethods,
	($paymentMethods) => $paymentMethods.filter((m) => m.isActive)
);

/**
 * Load payment methods from API (with IndexedDB fallback)
 */
export async function loadPaymentMethods(): Promise<void> {
	paymentMethodsLoading.set(true);

	try {
		// Try API first if enabled
		if (USE_API) {
			try {
				const apiMethods = await paymentMethodsApi.getAllPaymentMethods();

				// Update store
				paymentMethods.set(apiMethods);

				// Update IndexedDB for offline access
				for (const method of apiMethods) {
					await paymentMethodStore.update(method);
				}

				console.log(`[PaymentMethods] Loaded ${apiMethods.length} methods from API`);
				return;
			} catch (error) {
				console.warn('[PaymentMethods] API load failed, falling back to IndexedDB:', error);
				showToast('Using offline data. Connect to internet for latest data.', 'warning');
				// Fall through to IndexedDB loading
			}
		}

		// Fallback to IndexedDB
		const allMethods = await paymentMethodStore.getAll();
		paymentMethods.set(allMethods);
		console.log(`[PaymentMethods] Loaded ${allMethods.length} methods from IndexedDB`);
	} catch (error) {
		console.error('[PaymentMethods] Error loading payment methods:', error);
		showToast(
			error instanceof Error ? error.message : 'Failed to load payment methods',
			'error'
		);
	} finally {
		paymentMethodsLoading.set(false);
	}
}

/**
 * Add payment method
 */
export async function addPaymentMethod(method: PaymentMethod): Promise<void> {
	// Try API first if enabled
	if (USE_API) {
		try {
			const created = await paymentMethodsApi.createPaymentMethod(method);

			// Update IndexedDB for offline access
			await paymentMethodStore.update(created);

			// Update store
			paymentMethods.update((current) => [...current, created]);

			console.log('[PaymentMethods] Added method via API:', created.id);
			showToast('Payment method added successfully!', 'success');
			return;
		} catch (error) {
			console.warn('[PaymentMethods] API create failed, using local fallback:', error);
			showToast('Saved locally. Will sync when connection is restored.', 'warning');
			// Fall through to local creation
		}
	}

	// Fallback: Create locally (will sync later)
	try {
		await paymentMethodStore.create(method);
		paymentMethods.update((current) => [...current, method]);
		console.log('[PaymentMethods] Added method locally:', method.id);

		// Add to sync queue
		const { addToSyncQueue } = await import('$lib/db/sync');
		await addToSyncQueue({
			table: 'paymentMethods',
			recordId: method.id,
			operation: 'CREATE',
			data: method
		});

		showToast('Payment method added locally!', 'success');
	} catch (error) {
		console.error('[PaymentMethods] Error adding payment method:', error);
		showToast(
			error instanceof Error ? error.message : 'Failed to add payment method',
			'error'
		);
		throw error;
	}
}

/**
 * Update payment method
 */
export async function updatePaymentMethod(method: PaymentMethod): Promise<void> {
	// Try API first if enabled
	if (USE_API) {
		try {
			const updated = await paymentMethodsApi.updatePaymentMethod(method.id, method);

			// Update IndexedDB for offline access
			await paymentMethodStore.update(updated);

			// Update store
			paymentMethods.update((current) =>
				current.map((m) => (m.id === method.id ? updated : m))
			);

			console.log('[PaymentMethods] Updated method via API:', method.id);
			showToast('Payment method updated successfully!', 'success');
			return;
		} catch (error) {
			console.warn('[PaymentMethods] API update failed, using local fallback:', error);
			showToast('Updated locally. Will sync when connection is restored.', 'warning');
			// Fall through to local update
		}
	}

	// Fallback: Update locally (will sync later)
	try {
		await paymentMethodStore.update(method);
		paymentMethods.update((current) =>
			current.map((m) => (m.id === method.id ? method : m))
		);
		console.log('[PaymentMethods] Updated method locally:', method.id);

		// Add to sync queue
		const { addToSyncQueue } = await import('$lib/db/sync');
		await addToSyncQueue({
			table: 'paymentMethods',
			recordId: method.id,
			operation: 'UPDATE',
			data: method
		});

		showToast('Payment method updated locally!', 'success');
	} catch (error) {
		console.error('[PaymentMethods] Error updating payment method:', error);
		showToast(
			error instanceof Error ? error.message : 'Failed to update payment method',
			'error'
		);
		throw error;
	}
}

/**
 * Delete payment method (soft delete)
 */
export async function deletePaymentMethod(methodId: string): Promise<void> {
	// Try API first if enabled
	if (USE_API) {
		try {
			await paymentMethodsApi.deletePaymentMethod(methodId);

			// Update IndexedDB for offline access
			await paymentMethodStore.delete(methodId);

			// Update store (remove from list)
			paymentMethods.update((current) => current.filter((m) => m.id !== methodId));

			console.log('[PaymentMethods] Deleted method via API:', methodId);
			showToast('Payment method deleted successfully!', 'success');
			return;
		} catch (error) {
			console.warn('[PaymentMethods] API delete failed, using local fallback:', error);
			showToast('Deleted locally. Will sync when connection is restored.', 'warning');
			// Fall through to local deletion
		}
	}

	// Fallback: Delete locally (will sync later)
	try {
		await paymentMethodStore.delete(methodId);
		paymentMethods.update((current) => current.filter((m) => m.id !== methodId));
		console.log('[PaymentMethods] Deleted method locally:', methodId);

		// Add to sync queue
		const { addToSyncQueue } = await import('$lib/db/sync');
		await addToSyncQueue({
			table: 'paymentMethods',
			recordId: methodId,
			operation: 'DELETE',
			data: { id: methodId }
		});

		showToast('Payment method deleted locally!', 'success');
	} catch (error) {
		console.error('[PaymentMethods] Error deleting payment method:', error);
		showToast(
			error instanceof Error ? error.message : 'Failed to delete payment method',
			'error'
		);
		throw error;
	}
}

/**
 * Get payment method by ID
 */
export function getPaymentMethodById(methodId: string): PaymentMethod | undefined {
	const $paymentMethods = get(paymentMethods);
	return $paymentMethods.find((m) => m.id === methodId);
}

/**
 * Set default payment method
 * Unsets current default and sets new default
 */
export async function setDefaultPaymentMethod(methodId: string): Promise<void> {
	const $paymentMethods = get(paymentMethods);

	// Find current default
	const currentDefault = $paymentMethods.find((m) => m.isDefault);

	// If already default, nothing to do
	if (currentDefault?.id === methodId) {
		return;
	}

	// Unset current default (if exists)
	if (currentDefault) {
		await updatePaymentMethod({ ...currentDefault, isDefault: false });
	}

	// Set new default
	const newDefault = $paymentMethods.find((m) => m.id === methodId);
	if (newDefault) {
		await updatePaymentMethod({ ...newDefault, isDefault: true });
	}
}
