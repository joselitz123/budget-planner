/**
 * Payment Methods API Client
 *
 * Handles all payment method CRUD operations with the backend API.
 */

import { apiClient } from './client';
import type { PaymentMethod } from '$lib/db/schema';

/**
 * Backend response type for payment methods
 */
export interface PaymentMethodResponse {
	id: string;
	name: string;
	type: string;
	lastFour?: string;
	brand?: string;
	isDefault: boolean;
	isActive: boolean;
	creditLimit?: number;
	currentBalance?: number;
	createdAt: string;
	updatedAt: string;
}

/**
 * Request type for creating payment methods
 */
export interface CreatePaymentMethodRequest {
	name: string;
	type: string;
	lastFour?: string;
	brand?: string;
	isDefault?: boolean;
	creditLimit?: number;
	currentBalance?: number;
}

/**
 * Request type for updating payment methods (all fields optional)
 */
export interface UpdatePaymentMethodRequest {
	name?: string;
	type?: string;
	lastFour?: string;
	brand?: string;
	isDefault?: boolean;
	isActive?: boolean;
	creditLimit?: number;
	currentBalance?: number;
}

/**
 * Convert backend response to frontend PaymentMethod type
 */
function adaptResponseToFrontend(response: PaymentMethodResponse): PaymentMethod {
	return {
		id: response.id,
		userId: '', // Will be filled by auth context
		name: response.name,
		type: response.type as PaymentMethod['type'],
		lastFour: response.lastFour || null,
		brand: response.brand || null,
		isDefault: response.isDefault,
		isActive: response.isActive,
		creditLimit: response.creditLimit || null,
		currentBalance: response.currentBalance || null,
		createdAt: response.createdAt,
		updatedAt: response.updatedAt
	};
}

/**
 * Convert frontend PaymentMethod to create request
 */
function adaptFrontendToCreateRequest(method: PaymentMethod): CreatePaymentMethodRequest {
	return {
		name: method.name,
		type: method.type,
		lastFour: method.lastFour || undefined,
		brand: method.brand || undefined,
		isDefault: method.isDefault,
		creditLimit: method.creditLimit || undefined,
		currentBalance: method.currentBalance || undefined
	};
}

/**
 * Convert frontend PaymentMethod to update request
 */
function adaptFrontendToUpdateRequest(method: PaymentMethod): UpdatePaymentMethodRequest {
	return {
		name: method.name,
		type: method.type,
		lastFour: method.lastFour || undefined,
		brand: method.brand || undefined,
		isDefault: method.isDefault,
		isActive: method.isActive,
		creditLimit: method.creditLimit || undefined,
		currentBalance: method.currentBalance || undefined
	};
}

/**
 * Payment Methods API Client
 */
export class PaymentMethodsApi {
	/**
	 * Get all payment methods for the current user
	 */
	async getAllPaymentMethods(): Promise<PaymentMethod[]> {
		const response = await apiClient.get<PaymentMethodResponse[]>('/payment-methods');
		return response.map(adaptResponseToFrontend);
	}

	/**
	 * Get a specific payment method by ID
	 */
	async getPaymentMethodById(id: string): Promise<PaymentMethod> {
		const response = await apiClient.get<PaymentMethodResponse>(`/payment-methods/${id}`);
		return adaptResponseToFrontend(response);
	}

	/**
	 * Create a new payment method
	 */
	async createPaymentMethod(method: PaymentMethod): Promise<PaymentMethod> {
		const request = adaptFrontendToCreateRequest(method);
		const response = await apiClient.post<PaymentMethodResponse>('/payment-methods', request);
		return adaptResponseToFrontend(response);
	}

	/**
	 * Update an existing payment method
	 */
	async updatePaymentMethod(id: string, method: PaymentMethod): Promise<PaymentMethod> {
		const request = adaptFrontendToUpdateRequest(method);
		const response = await apiClient.put<PaymentMethodResponse>(
			`/payment-methods/${id}`,
			request
		);
		return adaptResponseToFrontend(response);
	}

	/**
	 * Delete (soft delete) a payment method
	 */
	async deletePaymentMethod(id: string): Promise<void> {
		return apiClient.delete<void>(`/payment-methods/${id}`);
	}
}

// Export singleton instance
export const paymentMethodsApi = new PaymentMethodsApi();
