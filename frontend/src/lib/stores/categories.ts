import { writable, derived, get } from 'svelte/store';
import { categoryStore } from '$lib/db/stores';
import type { Category } from '$lib/db/schema';
import { showToast } from '$lib/stores/ui';

/**
 * Category state management
 */

// User categories
export const userCategories = writable<Category[]>([]);

// System categories (read-only, hardcoded for now)
export const systemCategories = writable<Category[]>([
	{
		id: 'sys-housing',
		userId: null,
		name: 'Housing',
		icon: '🏠',
		color: '#A78BFA',
		isDefault: true,
		defaultLimit: null,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString()
	},
	{
		id: 'sys-food',
		userId: null,
		name: 'Food',
		icon: '🍔',
		color: '#FBBF24',
		isDefault: true,
		defaultLimit: null,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString()
	},
	{
		id: 'sys-transport',
		userId: null,
		name: 'Transportation',
		icon: '🚗',
		color: '#60A5FA',
		isDefault: true,
		defaultLimit: null,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString()
	},
	{
		id: 'sys-healthcare',
		userId: null,
		name: 'Health Care',
		icon: '💊',
		color: '#34D399',
		isDefault: true,
		defaultLimit: null,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString()
	},
	{
		id: 'sys-bills',
		userId: null,
		name: 'Bills',
		icon: '💡',
		color: '#F87171',
		isDefault: true,
		defaultLimit: null,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString()
	},
	{
		id: 'sys-entertainment',
		userId: null,
		name: 'Entertainment',
		icon: '🎬',
		color: '#9CA3AF',
		isDefault: true,
		defaultLimit: null,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString()
	},
	{
		id: 'sys-social',
		userId: null,
		name: 'Social',
		icon: '🎉',
		color: '#F472B6',
		isDefault: true,
		defaultLimit: null,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString()
	}
]);

// Derived: all categories combined
export const allCategories = derived(
	[userCategories, systemCategories],
	([$userCategories, $systemCategories]) => {
		return [...$systemCategories, ...$userCategories];
	}
);

/**
 * Load user categories from IndexedDB
 */
export async function loadCategories(userId?: string): Promise<void> {
	try {
		if (userId) {
			const categories = await categoryStore.getByUser(userId);
			userCategories.set(categories);
			console.log(`[Categories] Loaded ${categories.length} user categories`);
		}
	} catch (error) {
		console.error('[Categories] Error loading categories:', error);
		showToast(
			error instanceof Error ? error.message : 'Failed to load categories',
			'error'
		);
	}
}

/**
 * Get category by ID
 */
export function getCategoryById(categoryId: string): Category | undefined {
	const $userCategories = get(userCategories);
	const $systemCategories = get(systemCategories);

	return [...$systemCategories, ...$userCategories].find((c) => c.id === categoryId);
}

/**
 * Add user category
 */
export async function addCategory(category: Category): Promise<void> {
	try {
		await categoryStore.create(category);
		userCategories.update((current) => [...current, category]);
		console.log('[Categories] Added category:', category.id);
	} catch (error) {
		console.error('[Categories] Error adding category:', error);
		showToast(
			error instanceof Error ? error.message : 'Failed to add category',
			'error'
		);
		throw error;
	}
}
