import { apiClient } from './client';
import type { Transaction } from '$lib/db/schema';

/**
 * Transaction API types (matching backend responses)
 */

export interface TransactionResponse {
	id: string;
	budgetId?: string;
	categoryId?: string;
	paymentMethodId?: string;
	amount: number;
	type: string; // 'expense' | 'income' in backend
	isTransfer: boolean;
	transferToAccountId?: string;
	description?: string;
	transactionDate: string;
	isRecurring: boolean;
	recurrencePattern?: any;
	createdAt: string;
	updatedAt: string;
}

export interface CreateTransactionRequest {
	budgetId?: string;
	categoryId?: string;
	paymentMethodId?: string;
	amount: number;
	type: string; // 'expense' | 'income'
	isTransfer?: boolean;
	transferToAccountId?: string;
	description?: string;
	transactionDate: string;
	isRecurring?: boolean;
	recurrencePattern?: any;
}

export interface UpdateTransactionRequest {
	budgetId?: string;
	categoryId?: string;
	paymentMethodId?: string;
	amount?: number;
	type?: string;
	isTransfer?: boolean;
	transferToAccountId?: string;
	description?: string;
	transactionDate?: string;
	isRecurring?: boolean;
	recurrencePattern?: any;
}

export interface TransactionFilters {
	budgetId?: string;
	categoryId?: string;
	type?: string;
	startDate?: string;
	endDate?: string;
	isRecurring?: boolean;
}

/**
 * Type adapters to convert between backend and frontend types
 */

/**
 * Convert backend TransactionResponse to frontend Transaction
 */
export function adaptTransactionToBackend(transaction: Transaction): CreateTransactionRequest {
	return {
		budgetId: transaction.budgetId || undefined,
		categoryId: transaction.categoryId || undefined,
		paymentMethodId: transaction.paymentMethodId || undefined,
		amount: transaction.amount,
		type: transaction.transactionType, // Map transactionType -> type
		isTransfer: false,
		description: transaction.description || undefined,
		transactionDate: transaction.transactionDate,
		isRecurring: transaction.isRecurring,
		recurrencePattern: undefined
	};
}

/**
 * Convert frontend Transaction to backend TransactionResponse
 */
export function adaptTransactionToFrontend(response: TransactionResponse): Transaction {
	return {
		id: response.id,
		userId: '', // Will be filled by auth context
		budgetId: response.budgetId || '',
		categoryId: response.categoryId || '',
		amount: response.amount,
		description: response.description || null,
		transactionDate: response.transactionDate,
		transactionType: response.type as 'expense' | 'income', // Map type -> transactionType
		paymentMethodId: response.paymentMethodId || null,
		paid: false, // Backend doesn't track this, defaults to false
		dueDate: null, // Backend doesn't have this field
		isRecurring: response.isRecurring,
		notes: null, // Backend doesn't have this field
		createdAt: response.createdAt,
		updatedAt: response.updatedAt
	};
}

/**
 * Transactions API client
 */
export class TransactionsApi {
	/**
	 * Get all transactions with optional filters
	 * GET /api/transactions
	 */
	async getTransactions(filters?: TransactionFilters): Promise<Transaction[]> {
		// Build query string from filters
		const params = new URLSearchParams();
		if (filters?.budgetId) params.append('budgetId', filters.budgetId);
		if (filters?.categoryId) params.append('categoryId', filters.categoryId);
		if (filters?.type) params.append('type', filters.type);
		if (filters?.startDate) params.append('startDate', filters.startDate);
		if (filters?.endDate) params.append('endDate', filters.endDate);
		if (filters?.isRecurring !== undefined) params.append('isRecurring', String(filters.isRecurring));

		const queryString = params.toString();
		const path = `/transactions${queryString ? `?${queryString}` : ''}`;

		const response = await apiClient.get<TransactionResponse[]>(path);
		return response.map(adaptTransactionToFrontend);
	}

	/**
	 * Get a single transaction by ID
	 * GET /api/transactions/{id}
	 */
	async getTransactionById(id: string): Promise<Transaction> {
		const response = await apiClient.get<TransactionResponse>(`/transactions/${id}`);
		return adaptTransactionToFrontend(response);
	}

	/**
	 * Create a new transaction
	 * POST /api/transactions
	 */
	async createTransaction(transaction: Transaction): Promise<Transaction> {
		const request = adaptTransactionToBackend(transaction);
		const response = await apiClient.post<TransactionResponse>('/transactions', request);
		return adaptTransactionToFrontend(response);
	}

	/**
	 * Update a transaction
	 * PUT /api/transactions/{id}
	 */
	async updateTransaction(id: string, transaction: Transaction): Promise<Transaction> {
		const request: UpdateTransactionRequest = {
			budgetId: transaction.budgetId || undefined,
			categoryId: transaction.categoryId || undefined,
			paymentMethodId: transaction.paymentMethodId || undefined,
			amount: transaction.amount,
			type: transaction.transactionType,
			isTransfer: false,
			description: transaction.description || undefined,
			transactionDate: transaction.transactionDate,
			isRecurring: transaction.isRecurring
		};

		const response = await apiClient.put<TransactionResponse>(`/transactions/${id}`, request);
		return adaptTransactionToFrontend(response);
	}

	/**
	 * Delete a transaction
	 * DELETE /api/transactions/{id}
	 */
	async deleteTransaction(id: string): Promise<void> {
		return apiClient.delete<void>(`/transactions/${id}`);
	}

	/**
	 * Get transactions for a specific budget
	 * GET /api/budgets/{budgetId}/transactions
	 */
	async getBudgetTransactions(budgetId: string): Promise<Transaction[]> {
		const response = await apiClient.get<TransactionResponse[]>(`/budgets/${budgetId}/transactions`);
		return response.map(adaptTransactionToFrontend);
	}
}

// Export singleton instance
export const transactionsApi = new TransactionsApi();
