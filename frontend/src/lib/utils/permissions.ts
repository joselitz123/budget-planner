import type { SharedBudget } from '$lib/api/shares';

/**
 * Check if user can edit a budget
 * @param sharedBudget - The shared budget object
 * @returns true if user has edit permission
 */
export function canEditBudget(sharedBudget: SharedBudget | null): boolean {
	if (!sharedBudget) return false;
	return sharedBudget.permission === 'edit';
}

/**
 * Check if user can only view a budget
 * @param sharedBudget - The shared budget object
 * @returns true if user has view-only permission
 */
export function canOnlyViewBudget(sharedBudget: SharedBudget | null): boolean {
	if (!sharedBudget) return false;
	return sharedBudget.permission === 'view';
}

/**
 * Get permission badge variant
 * @param permission - The permission level
 * @returns Badge variant name
 */
export function getPermissionBadgeVariant(permission: 'view' | 'edit'): string {
	return permission === 'edit' ? 'default' : 'secondary';
}

/**
 * Get permission display text
 * @param permission - The permission level
 * @returns Human-readable permission text
 */
export function getPermissionText(permission: 'view' | 'edit'): string {
	return permission === 'edit' ? 'View & Edit' : 'View Only';
}

/**
 * Check if user is the owner of a budget
 * @param userId - Current user ID
 * @param budgetOwnerId - Budget owner ID
 * @returns true if user is the owner
 */
export function isOwner(userId: string, budgetOwnerId: string): boolean {
	return userId === budgetOwnerId;
}
