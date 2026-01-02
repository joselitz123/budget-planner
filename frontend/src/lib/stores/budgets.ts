import { writable, derived, get } from 'svelte/store';
import { budgetStore } from '$lib/db/stores';
import type { Budget } from '$lib/db/schema';
import { getMonthKey, parseMonthKey } from '$lib/utils/format';
import { budgetsApi } from '$lib/api/budgets';
import { showToast } from '$lib/stores/ui';

/**
 * Budget state management
 * Supports both API and IndexedDB (offline-first)
 */

// Flag to enable/disable API integration
// Set to false to use only IndexedDB (offline mode)
let USE_API = true;

// All budgets
export const budgets = writable<Budget[]>([]);

// Budget loading state
export const budgetsLoading = writable<boolean>(false);

// Current budget (for selected month)
export const currentBudget = writable<Budget | null>(null);

// Current month (default to current month)
export const currentMonth = writable<string>(getMonthKey(new Date()));

// Derived: current month budget
export const currentMonthBudget = derived(
	[budgets, currentMonth],
	([$budgets, $currentMonth]) => {
		return $budgets.find((b) => b.month === $currentMonth) || null;
	}
);

// Derived: all months (sorted)
export const allMonths = derived(budgets, ($budgets) => {
	const months = [...new Set($budgets.map((b) => b.month))];
	return months.sort().reverse();
});

/**
 * Load budgets from API (with IndexedDB fallback)
 */
export async function loadBudgets(): Promise<void> {
	budgetsLoading.set(true);

	try {
		// Try API first if enabled
		if (USE_API) {
			try {
				const apiBudgets = await budgetsApi.getAllBudgets();

				// Convert API response to Budget type
				const budgetsData: Budget[] = apiBudgets.map(b => ({
					id: b.id,
					userId: b.userId,
					month: b.month.substring(0, 7), // Convert YYYY-MM-DD to YYYY-MM
					totalLimit: b.totalLimit,
					createdAt: b.createdAt,
					updatedAt: b.updatedAt
				}));

				// Update store
				budgets.set(budgetsData);

				// Update IndexedDB for offline access
				for (const budget of budgetsData) {
					await budgetStore.update(budget);
				}

				console.log(`[Budgets] Loaded ${budgetsData.length} budgets from API`);
				return;
			} catch (error) {
				console.warn('[Budgets] API load failed, falling back to IndexedDB:', error);
				showToast('Using offline data. Connect to internet for latest data.', 'warning');
				// Fall through to IndexedDB loading
			}
		}

		// Fallback to IndexedDB
		const allBudgets = await budgetStore.getAll();
		budgets.set(allBudgets);
		console.log(`[Budgets] Loaded ${allBudgets.length} budgets from IndexedDB`);
	} catch (error) {
		console.error('[Budgets] Error loading budgets:', error);
		showToast(
			error instanceof Error ? error.message : 'Failed to load budgets',
			'error'
		);
	} finally {
		budgetsLoading.set(false);
	}
}

/**
 * Get or create budget for month
 */
export async function getOrCreateBudgetForMonth(month: string): Promise<Budget | null> {
	const existing = get(budgets).find((b) => b.month === month);
	if (existing) {
		currentBudget.set(existing);
		return existing;
	}

	// TODO: Create new budget via API
	return null;
}

/**
 * Create budget for current month (auto-creation)
 * @param userId - The user ID to create the budget for
 * @param options - Optional parameters for budget creation
 * @param options.name - Optional custom name for the budget (defaults to "${month} Budget")
 * @param options.totalLimit - Optional total limit for the budget (defaults to 2000)
 */
export async function createBudgetForCurrentMonth(
	userId: string,
	options?: { name?: string; totalLimit?: number }
): Promise<Budget> {
	const month = getMonthKey(new Date());
	
	// Extract options with defaults
	const name = options?.name || `${month} Budget`;
	const totalLimit = options?.totalLimit ?? 2000;

	// Check if budget already exists
	const existing = get(budgets).find((b) => b.month === month);
	if (existing) {
		currentBudget.set(existing);
		return existing;
	}

	// Try to create via API if enabled
	if (USE_API) {
		try {
			// Convert YYYY-MM to YYYY-MM-DD (first day of month)
			const monthDate = `${month}-01`;

			const apiBudget = await budgetsApi.createBudget({
				month: monthDate,
				totalLimit
			});

			const newBudget: Budget = {
				id: apiBudget.id,
				userId: apiBudget.userId,
				month: apiBudget.month.substring(0, 7),
				totalLimit: apiBudget.totalLimit,
				createdAt: apiBudget.createdAt,
				updatedAt: apiBudget.updatedAt
			};

			// Update IndexedDB for offline access
			await budgetStore.create(newBudget);

			// Update store
			budgets.update((b) => [...b, newBudget]);
			currentBudget.set(newBudget);

			console.log('[Budgets] Created budget via API for month:', month, 'with name:', name);
			return newBudget;
		} catch (error) {
			console.warn('[Budgets] API create failed, using local fallback:', error);
			showToast('Saved locally. Will sync when connection is restored.', 'warning');
			// Fall through to local creation
		}
	}

	// Fallback: Create locally (will sync later)
	const newBudget: Budget = {
		id: crypto.randomUUID(),
		userId,
		month,
		totalLimit,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString()
	};

	await budgetStore.create(newBudget);
	budgets.update((b) => [...b, newBudget]);
	currentBudget.set(newBudget);

	// Queue for sync
	const { addToSyncQueue } = await import('$lib/db/sync');
	await addToSyncQueue({
		table: 'budgets',
		recordId: newBudget.id,
		operation: 'CREATE',
		data: newBudget
	});

	console.log('[Budgets] Auto-created budget locally for month:', month, 'with name:', name);
	return newBudget;
}

/**
 * Set current month
 */
export function setCurrentMonth(month: string): void {
	currentMonth.set(month);
	getOrCreateBudgetForMonth(month);
}

/**
 * Navigate to previous month
 */
export function goToPreviousMonth(): void {
	const $currentMonth = get(currentMonth);
	const date = parseMonthKey($currentMonth);
	date.setMonth(date.getMonth() - 1);
	const prevMonth = getMonthKey(date);
	setCurrentMonth(prevMonth);
}

/**
 * Navigate to next month
 */
export function goToNextMonth(): void {
	const $currentMonth = get(currentMonth);
	const date = parseMonthKey($currentMonth);
	date.setMonth(date.getMonth() + 1);
	const nextMonth = getMonthKey(date);
	setCurrentMonth(nextMonth);
}
