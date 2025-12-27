import type { ApiResponse } from '$lib/db/schema';

/**
 * Base API client with error handling
 */
export class ApiClient {
	private baseUrl: string;

	constructor(baseUrl?: string) {
		this.baseUrl = baseUrl || import.meta.env.VITE_PUBLIC_API_URL || 'http://localhost:8080/api';
	}

	/**
	 * Make a GET request
	 */
	async get<T>(path: string, options?: RequestInit): Promise<T> {
		const response = await fetch(`${this.baseUrl}${path}`, {
			...options,
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
				...options?.headers
			}
		});

		if (!response.ok) {
			throw new Error(`API Error: ${response.statusText}`);
		}

		const result: ApiResponse<T> = await response.json();

		if (!result.success) {
			throw new Error(result.error?.message || 'API request failed');
		}

		return result.data;
	}

	/**
	 * Make a POST request
	 */
	async post<T>(path: string, data?: any, options?: RequestInit): Promise<T> {
		const response = await fetch(`${this.baseUrl}${path}`, {
			...options,
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				...options?.headers
			},
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			throw new Error(`API Error: ${response.statusText}`);
		}

		const result: ApiResponse<T> = await response.json();

		if (!result.success) {
			throw new Error(result.error?.message || 'API request failed');
		}

		return result.data;
	}

	/**
	 * Make a PUT request
	 */
	async put<T>(path: string, data?: any, options?: RequestInit): Promise<T> {
		const response = await fetch(`${this.baseUrl}${path}`, {
			...options,
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json',
				...options?.headers
			},
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			throw new Error(`API Error: ${response.statusText}`);
		}

		const result: ApiResponse<T> = await response.json();

		if (!result.success) {
			throw new Error(result.error?.message || 'API request failed');
		}

		return result.data;
	}

	/**
	 * Make a DELETE request
	 */
	async delete<T>(path: string, options?: RequestInit): Promise<T> {
		const response = await fetch(`${this.baseUrl}${path}`, {
			...options,
			method: 'DELETE',
			headers: {
				'Content-Type': 'application/json',
				...options?.headers
			}
		});

		if (!response.ok) {
			throw new Error(`API Error: ${response.statusText}`);
		}

		const result: ApiResponse<T> = await response.json();

		if (!result.success) {
			throw new Error(result.error?.message || 'API request failed');
		}

		return result.data;
	}
}

// Export singleton instance
export const apiClient = new ApiClient();
