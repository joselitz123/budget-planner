import { writable, derived, get } from 'svelte/store';
import { budgetStore } from '$lib/db/stores';
import type { Budget } from '$lib/db/schema';
import { getMonthKey, parseMonthKey } from '$lib/utils/format';

/**
 * Budget state management
 */

// All budgets
export const budgets = writable<Budget[]>([]);

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
 * Load budgets from IndexedDB
 */
export async function loadBudgets(): Promise<void> {
	try {
		const allBudgets = await budgetStore.getAll();
		budgets.set(allBudgets);
		console.log(`[Budgets] Loaded ${allBudgets.length} budgets`);
	} catch (error) {
		console.error('[Budgets] Error loading budgets:', error);
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
 */
export async function createBudgetForCurrentMonth(userId: string): Promise<Budget> {
	const month = getMonthKey(new Date());

	// Check if budget already exists
	const existing = get(budgets).find((b) => b.month === month);
	if (existing) {
		currentBudget.set(existing);
		return existing;
	}

	// Create new budget
	const newBudget: Budget = {
		id: crypto.randomUUID(),
		userId,
		month,
		totalLimit: 2000, // Default limit - could be made configurable
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString()
	};

	await budgetStore.create(newBudget);
	budgets.update((b) => [...b, newBudget]);
	currentBudget.set(newBudget);

	console.log('[Budgets] Auto-created budget for month:', month);
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
