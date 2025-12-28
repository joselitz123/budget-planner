import { apiClient } from './client';
import type { Category } from '$lib/db/schema';

/**
 * Category API types (matching backend responses)
 */

export interface CategoryResponse {
	id: string;
	name: string;
	icon?: string;
	color: string;
	isSystem: boolean;
	defaultLimit?: number;
}

export interface CreateCategoryRequest {
	name: string;
	icon?: string;
	color: string;
	defaultLimit?: number;
}

export interface UpdateCategoryRequest {
	name?: string;
	icon?: string;
	color?: string;
	defaultLimit?: number;
}

/**
 * Type adapters to convert between backend and frontend types
 */

/**
 * Convert frontend Category to backend CreateCategoryRequest
 */
export function adaptCategoryToBackend(category: Category): CreateCategoryRequest {
	return {
		name: category.name,
		icon: category.icon || undefined,
		color: category.color || '#000000',
		defaultLimit: category.defaultLimit || undefined
	};
}

/**
 * Convert backend CategoryResponse to frontend Category
 */
export function adaptCategoryToFrontend(response: CategoryResponse, userId: string | null): Category {
	return {
		id: response.id,
		userId, // null for system categories
		name: response.name,
		icon: response.icon || null,
		color: response.color,
		isDefault: response.isSystem,
		defaultLimit: response.defaultLimit || null,
		createdAt: new Date().toISOString(), // Backend doesn't return this
		updatedAt: new Date().toISOString()  // Backend doesn't return this
	};
}

/**
 * Categories API client
 */
export class CategoriesApi {
	/**
	 * Get all categories for the current user
	 * GET /api/categories
	 */
	async getAllCategories(userId: string): Promise<Category[]> {
		const response = await apiClient.get<CategoryResponse[]>('/categories');
		return response.map(cat => adaptCategoryToFrontend(cat, userId));
	}

	/**
	 * Get system categories
	 * GET /api/categories/system
	 */
	async getSystemCategories(): Promise<Category[]> {
		const response = await apiClient.get<CategoryResponse[]>('/categories/system');
		return response.map(cat => adaptCategoryToFrontend(cat, null)); // System categories have null userId
	}

	/**
	 * Get a single category by ID
	 * GET /api/categories/{id}
	 */
	async getCategoryById(id: string, userId: string): Promise<Category> {
		const response = await apiClient.get<CategoryResponse>(`/categories/${id}`);
		return adaptCategoryToFrontend(response, userId);
	}

	/**
	 * Create a new category
	 * POST /api/categories
	 */
	async createCategory(category: Category): Promise<Category> {
		const request = adaptCategoryToBackend(category);
		const response = await apiClient.post<CategoryResponse>('/categories', request);
		return adaptCategoryToFrontend(response, category.userId);
	}

	/**
	 * Update a category
	 * PUT /api/categories/{id}
	 */
	async updateCategory(id: string, category: Category): Promise<Category> {
		const request: UpdateCategoryRequest = {
			name: category.name,
			icon: category.icon || undefined,
			color: category.color || undefined,
			defaultLimit: category.defaultLimit || undefined
		};

		const response = await apiClient.put<CategoryResponse>(`/categories/${id}`, request);
		return adaptCategoryToFrontend(response, category.userId);
	}

	/**
	 * Delete a category
	 * DELETE /api/categories/{id}
	 */
	async deleteCategory(id: string): Promise<void> {
		return apiClient.delete<void>(`/categories/${id}`);
	}
}

// Export singleton instance
export const categoriesApi = new CategoriesApi();
