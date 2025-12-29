import { writable, derived } from 'svelte/store';
import { sharesApi } from '$lib/api/shares';
import type { ShareInvitation, SharedBudget, ShareAccess } from '$lib/api/shares';
import { showToast } from '$lib/stores/ui';

/**
 * Sharing state management
 *
 * Manages invitations, shared budgets, and access lists
 */

// Writable stores
export const invitations = writable<ShareInvitation[]>([]);
export const invitationsLoading = writable<boolean>(false);
export const sharedBudgets = writable<SharedBudget[]>([]);
export const sharedBudgetsLoading = writable<boolean>(false);
export const budgetAccessList = writable<Map<string, ShareAccess[]>>(new Map());

// Derived store for pending invitations only
export const pendingInvitations = derived(invitations, ($invitations) =>
	$invitations.filter((i) => i.status === 'pending')
);

/**
 * Load pending invitations for current user
 */
export async function loadInvitations(): Promise<void> {
	invitationsLoading.set(true);
	try {
		const data = await sharesApi.getMyInvitations();
		invitations.set(data);
	} catch (error) {
		console.error('[Shares] Failed to load invitations:', error);
		showToast('Failed to load invitations', 'error');
	} finally {
		invitationsLoading.set(false);
	}
}

/**
 * Load budgets shared with current user
 */
export async function loadSharedBudgets(): Promise<void> {
	sharedBudgetsLoading.set(true);
	try {
		const data = await sharesApi.getSharedBudgets();
		sharedBudgets.set(data);
	} catch (error) {
		console.error('[Shares] Failed to load shared budgets:', error);
		showToast('Failed to load shared budgets', 'error');
	} finally {
		sharedBudgetsLoading.set(false);
	}
}

/**
 * Load who has access to a specific budget
 */
export async function loadBudgetAccess(budgetId: string): Promise<void> {
	try {
		const data = await sharesApi.getBudgetSharing(budgetId);
		budgetAccessList.update((list) => {
			const newList = new Map(list);
			newList.set(budgetId, data);
			return newList;
		});
	} catch (error) {
		console.error('[Shares] Failed to load budget access:', error);
		showToast('Failed to load sharing info', 'error');
	}
}

/**
 * Create a new share invitation
 */
export async function createInvitation(
	budgetId: string,
	email: string,
	permission: 'view' | 'edit'
): Promise<void> {
	try {
		await sharesApi.createInvitation({
			budgetId,
			recipientEmail: email,
			permission
		});
		showToast(`Invitation sent to ${email}`, 'success');
		// Reload invitations to update the list
		await loadInvitations();
	} catch (error: any) {
		console.error('[Shares] Failed to create invitation:', error);
		const errorMessage = error?.message || 'Failed to send invitation';
		showToast(errorMessage, 'error');
		throw error;
	}
}

/**
 * Respond to a share invitation (accept or decline)
 */
export async function respondToInvitation(invitationId: string, accept: boolean): Promise<void> {
	try {
		await sharesApi.respondToInvitation(invitationId, {
			status: accept ? 'accepted' : 'declined'
		});
		// Remove from pending invitations
		invitations.update((list) => list.filter((i) => i.id !== invitationId));
		// Reload shared budgets
		await loadSharedBudgets();
	} catch (error: any) {
		console.error('[Shares] Failed to respond to invitation:', error);
		const errorMessage = error?.message || 'Failed to respond to invitation';
		showToast(errorMessage, 'error');
		throw error;
	}
}

/**
 * Cancel a pending invitation
 */
export async function cancelInvitation(invitationId: string): Promise<void> {
	try {
		await sharesApi.cancelInvitation(invitationId);
		// Remove from invitations list
		invitations.update((list) => list.filter((i) => i.id !== invitationId));
		showToast('Invitation cancelled', 'success');
	} catch (error: any) {
		console.error('[Shares] Failed to cancel invitation:', error);
		const errorMessage = error?.message || 'Failed to cancel invitation';
		showToast(errorMessage, 'error');
		throw error;
	}
}

/**
 * Remove someone's access to a budget
 */
export async function removeAccess(accessId: string, budgetId: string): Promise<void> {
	try {
		await sharesApi.removeAccess(accessId);
		// Update budget access list
		budgetAccessList.update((list) => {
			const newList = new Map(list);
			const currentList = newList.get(budgetId) || [];
			newList.set(
				budgetId,
				currentList.filter((access) => access.id !== accessId)
			);
			return newList;
		});
		showToast('Access removed', 'success');
	} catch (error: any) {
		console.error('[Shares] Failed to remove access:', error);
		const errorMessage = error?.message || 'Failed to remove access';
		showToast(errorMessage, 'error');
		throw error;
	}
}

/**
 * Get access list for a specific budget from store
 */
export function getBudgetAccess(budgetId: string): ShareAccess[] {
	let accessList: ShareAccess[] = [];
	budgetAccessList.subscribe((list) => {
		accessList = list.get(budgetId) || [];
	})();
	return accessList;
}
