import { derived } from 'svelte/store';
import { currentBudgetTransactions, transactions } from './transactions';
import { allCategories } from './categories';
import { currentBudget } from './budgets';

/**
 * Analytics data types
 */
export interface CategorySpend {
	categoryId: string;
	name: string;
	icon: string | null;
	color: string | null;
	amount: number;
	percentage: number;
	cumulativePercent: number;
}

export interface MonthlyTrend {
	month: string; // "YYYY-MM" format
	amount: number;
	monthLabel: string; // "Jan 2024" format
}

/**
 * Analytics state management
 * All analytics are derived from transaction data (computed client-side)
 */

// Loading state (analytics are computed synchronously from stores)
export const analyticsLoading = derived(
	[currentBudgetTransactions, allCategories],
	() => false
);

// Spending by category (for current budget/month)
export const spendingByCategory = derived(
	[currentBudgetTransactions, allCategories],
	([$transactions, $categories]): CategorySpend[] => {
		// Group expenses by category
		const categoryMap = new Map<string, number>();

		$transactions
			.filter((t) => t.transactionType === 'expense')
			.forEach((t) => {
				const current = categoryMap.get(t.categoryId) || 0;
				categoryMap.set(t.categoryId, current + t.amount);
			});

		// Convert to array with category details
		const spending = Array.from(categoryMap.entries())
			.map(([categoryId, amount]) => {
				const category = $categories.find((c) => c.id === categoryId);
				return {
					categoryId,
					name: category?.name || 'Unknown',
					icon: category?.icon || '❓',
					color: category?.color || '#9CA3AF',
					amount
				};
			})
			.sort((a, b) => b.amount - a.amount);

		// Calculate percentages
		const total = spending.reduce((sum, s) => sum + s.amount, 0);
		let cumulativePercent = 0;

		return spending.map((s) => {
			const percentage = total > 0 ? (s.amount / total) * 100 : 0;
			const result: CategorySpend = {
				...s,
				percentage: Math.round(percentage * 10) / 10, // Round to 1 decimal
				cumulativePercent
			};
			cumulativePercent += percentage;
			return result;
		});
	}
);

// Top 5 spending categories
export const topCategories = derived(
	spendingByCategory,
	($spendingByCategory): CategorySpend[] => $spendingByCategory.slice(0, 5)
);

// Top spending category (for summary card)
export const topCategory = derived(
	spendingByCategory,
	($spendingByCategory): CategorySpend | null => $spendingByCategory[0] || null
);

// Monthly trend (last 6 months)
export const monthlyTrends = derived(
	[transactions, allCategories],
	([$transactions]): MonthlyTrend[] => {
		// Group by month
		const monthMap = new Map<string, number>();

		$transactions
			.filter((t) => t.transactionType === 'expense')
			.forEach((t) => {
				const monthKey = t.transactionDate.substring(0, 7); // YYYY-MM
				const current = monthMap.get(monthKey) || 0;
				monthMap.set(monthKey, current + t.amount);
			});

		// Convert to array and sort
		const trends = Array.from(monthMap.entries())
			.map(([month, amount]) => {
				// Format month label (e.g., "Jan 2024")
				const date = new Date(month + '-01');
				const monthLabel = new Intl.DateTimeFormat('en-US', {
					month: 'short',
					year: 'numeric'
				}).format(date);

				return { month, amount, monthLabel };
			})
			.sort((a, b) => a.month.localeCompare(b.month));

		// Return last 6 months
		return trends.slice(-6);
	}
);

// Average daily spending (current month)
export const avgDailySpending = derived(
	currentBudgetTransactions,
	($transactions): number => {
		const expenses = $transactions.filter((t) => t.transactionType === 'expense');
		const total = expenses.reduce((sum, t) => sum + t.amount, 0);

		if (expenses.length === 0) return 0;

		// Get unique days
		const uniqueDays = new Set(expenses.map((t) => t.transactionDate)).size;
		return uniqueDays > 0 ? total / uniqueDays : 0;
	}
);

// Budget remaining (income - expenses)
export const budgetRemaining = derived(
	[currentBudgetTransactions, currentBudget],
	([$transactions, $currentBudget]): number => {
		if (!$currentBudget) return 0;

		const income = $transactions
			.filter((t) => t.transactionType === 'income')
			.reduce((sum, t) => sum + t.amount, 0);

		const expenses = $transactions
			.filter((t) => t.transactionType === 'expense')
			.reduce((sum, t) => sum + t.amount, 0);

		return income - expenses;
	}
);

// Total spent this month
export const totalSpentThisMonth = derived(
	currentBudgetTransactions,
	($transactions): number => {
		return $transactions
			.filter((t) => t.transactionType === 'expense')
			.reduce((sum, t) => sum + t.amount, 0);
	}
);
