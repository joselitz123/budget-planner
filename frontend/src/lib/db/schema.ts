/**
 * Budget Planner - IndexedDB Schema
 *
 * This file defines TypeScript interfaces for the IndexedDB database.
 * These types match the backend PostgreSQL schema.
 */

/**
 * Budget entity - monthly budget with limits
 */
export interface Budget {
	id: string;
	userId: string;
	month: string; // Format: "YYYY-MM"
	totalLimit: number;
	createdAt: string;
	updatedAt: string;
}

/**
 * Category entity - expense categories
 */
export interface Category {
	id: string;
	userId: string | null; // null for system categories
	name: string;
	icon: string | null;
	color: string | null;
	isDefault: boolean;
	defaultLimit: number | null;
	createdAt: string;
	updatedAt: string;
}

/**
 * Transaction entity - individual expense/income
 */
export interface Transaction {
	id: string;
	userId: string;
	budgetId: string;
	categoryId: string;
	amount: number;
	description: string | null;
	transactionDate: string;
	transactionType: 'expense' | 'income';
	paymentMethodId: string | null;
	paid: boolean;
	dueDate: string | null;
	isRecurring: boolean;
	notes: string | null;
	createdAt: string;
	updatedAt: string;
}

/**
 * Reflection entity - monthly budget reflection
 */
export interface Reflection {
	id: string;
	userId: string;
	budgetId: string;
	wins: string | null;
	didMeetBudget: boolean;
	reasons: string | null;
	improvements: string | null;
	createdAt: string;
	updatedAt: string;
}

/**
 * PaymentMethod entity - payment methods (cards, cash, etc.)
 */
export interface PaymentMethod {
	id: string;
	userId: string;
	name: string;
	type: string;
	isDefault: boolean;
	createdAt: string;
	updatedAt: string;
}

/**
 * SyncOperation entity - sync queue for offline changes
 */
export interface SyncOperation {
	id: string;
	table: 'budgets' | 'categories' | 'transactions' | 'reflections' | 'paymentMethods';
	recordId: string;
	operation: 'CREATE' | 'UPDATE' | 'DELETE';
	data: any;
	timestamp: string;
	status: 'pending' | 'syncing' | 'success' | 'failed';
	retryCount: number;
	error?: string;
}

/**
 * IndexedDB database schema definition
 */
export interface BudgetDB {
	budgets: {
		key: string;
		value: Budget;
		indexes: { 'by-month': string; 'by-user': string };
	};
	categories: {
		key: string;
		value: Category;
		indexes: { 'by-user': string };
	};
	transactions: {
		key: string;
		value: Transaction;
		indexes: { 'by-budget': string; 'by-date': string; 'by-category': string };
	};
	reflections: {
		key: string;
		value: Reflection;
		indexes: { 'by-budget': string };
	};
	paymentMethods: {
		key: string;
		value: PaymentMethod;
		indexes: { 'by-user': string };
	};
	syncQueue: {
		key: string;
		value: SyncOperation;
		indexes: { 'by-status': string };
	};
}

/**
 * API response types
 */
export interface ApiResponse<T> {
	success: boolean;
	data: T;
	error?: {
		message: string;
		code: string;
	};
}

/**
 * Dashboard analytics
 */
export interface DashboardAnalytics {
	openingBalance: number;
	totalIncome: number;
	totalExpenses: number;
	totalSavings: number;
	balanceForward: number;
}

/**
 * Spending breakdown by category
 */
export interface SpendingBreakdown {
	categoryName: string;
	amount: number;
	percentage: number;
	color: string;
}
