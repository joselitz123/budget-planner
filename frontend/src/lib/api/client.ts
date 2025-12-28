import type { ApiResponse } from '$lib/db/schema';
import { showToast } from '$lib/stores/ui';

/**
 * Authentication token provider interface
 * Allows for flexible auth integration (Clerk, custom, etc.)
 */
export interface AuthTokenProvider {
	getToken: () => Promise<string | null>;
}

/**
 * In-memory token storage for development/testing
 * In production, this should be replaced with Clerk or similar
 */
class DevTokenProvider implements AuthTokenProvider {
	private token: string | null = null;

	setToken(token: string) {
		this.token = token;
	}

	clearToken() {
		this.token = null;
	}

	async getToken(): Promise<string | null> {
		return this.token;
	}
}

// Global dev token provider (for development/testing)
export const devTokenProvider = new DevTokenProvider();

/**
 * Base API client with error handling and JWT authentication
 */
export class ApiClient {
	private baseUrl: string;
	private authProvider: AuthTokenProvider;

	constructor(baseUrl?: string, authProvider?: AuthTokenProvider) {
		this.baseUrl = baseUrl || import.meta.env.VITE_PUBLIC_API_URL || 'http://localhost:8080/api';
		this.authProvider = authProvider || devTokenProvider;
	}

	/**
	 * Get auth headers including JWT token
	 */
	private async getAuthHeaders(): Promise<HeadersInit> {
		const token = await this.authProvider.getToken();
		const headers: HeadersInit = {
			'Content-Type': 'application/json'
		};

		if (token) {
			headers['Authorization'] = `Bearer ${token}`;
		}

		return headers;
	}

	/**
	 * Handle API errors with user-friendly messages
	 */
	private handleApiError(response: Response, data?: any): never {
		let message = '';

		// Handle 401 Unauthorized - auth token expired or invalid
		if (response.status === 401) {
			message = 'Session expired. Please log in again.';
			console.error('[API] Unauthorized - please log in again');
			showToast(message, 'error', 5000);
			throw new Error(message);
		}

		// Handle 403 Forbidden
		if (response.status === 403) {
			message = 'You do not have permission to perform this action.';
			console.error('[API] Forbidden - insufficient permissions');
			showToast(message, 'error', 5000);
			throw new Error(message);
		}

		// Handle 404 Not Found
		if (response.status === 404) {
			message = 'The requested resource was not found.';
			console.error('[API] Resource not found');
			showToast(message, 'error', 3000);
			throw new Error(message);
		}

		// Handle 500 Server Error
		if (response.status >= 500) {
			message = 'Server error. Please try again later.';
			console.error('[API] Server error:', response.statusText);
			showToast(message, 'error', 5000);
			throw new Error(message);
		}

		// Handle other errors
		message = data?.error?.message || response.statusText || 'API request failed';
		console.error('[API] Error:', message);
		showToast(message, 'error', 3000);
		throw new Error(message);
	}

	/**
	 * Make a GET request
	 */
	async get<T>(path: string, options?: RequestInit): Promise<T> {
		const headers = await this.getAuthHeaders();

		const response = await fetch(`${this.baseUrl}${path}`, {
			...options,
			method: 'GET',
			headers: {
				...headers,
				...options?.headers
			}
		});

		const data = await response.json().catch(() => ({}));

		if (!response.ok) {
			this.handleApiError(response, data);
		}

		const result: ApiResponse<T> = data;

		if (!result.success) {
			throw new Error(result.error?.message || 'API request failed');
		}

		return result.data;
	}

	/**
	 * Make a POST request
	 */
	async post<T>(path: string, data?: any, options?: RequestInit): Promise<T> {
		const headers = await this.getAuthHeaders();

		const response = await fetch(`${this.baseUrl}${path}`, {
			...options,
			method: 'POST',
			headers: {
				...headers,
				...options?.headers
			},
			body: JSON.stringify(data)
		});

		const responseData = await response.json().catch(() => ({}));

		if (!response.ok) {
			this.handleApiError(response, responseData);
		}

		const result: ApiResponse<T> = responseData;

		if (!result.success) {
			throw new Error(result.error?.message || 'API request failed');
		}

		return result.data;
	}

	/**
	 * Make a PUT request
	 */
	async put<T>(path: string, data?: any, options?: RequestInit): Promise<T> {
		const headers = await this.getAuthHeaders();

		const response = await fetch(`${this.baseUrl}${path}`, {
			...options,
			method: 'PUT',
			headers: {
				...headers,
				...options?.headers
			},
			body: JSON.stringify(data)
		});

		const responseData = await response.json().catch(() => ({}));

		if (!response.ok) {
			this.handleApiError(response, responseData);
		}

		const result: ApiResponse<T> = responseData;

		if (!result.success) {
			throw new Error(result.error?.message || 'API request failed');
		}

		return result.data;
	}

	/**
	 * Make a DELETE request
	 */
	async delete<T>(path: string, options?: RequestInit): Promise<T> {
		const headers = await this.getAuthHeaders();

		const response = await fetch(`${this.baseUrl}${path}`, {
			...options,
			method: 'DELETE',
			headers: {
				...headers,
				...options?.headers
			}
		});

		const data = await response.json().catch(() => ({}));

		if (!response.ok) {
			this.handleApiError(response, data);
		}

		const result: ApiResponse<T> = data;

		if (!result.success) {
			throw new Error(result.error?.message || 'API request failed');
		}

		return result.data;
	}

	/**
	 * Set auth provider (e.g., for Clerk integration)
	 */
	setAuthProvider(provider: AuthTokenProvider) {
		this.authProvider = provider;
	}
}

// Export singleton instance
export const apiClient = new ApiClient();

/**
 * Helper function to set dev JWT token for testing
 * Call this with a valid JWT token to authenticate API requests
 */
export function setDevAuthToken(token: string) {
	devTokenProvider.setToken(token);
	console.log('[API] Dev auth token set');
}

/**
 * Clear dev auth token
 */
export function clearDevAuthToken() {
	devTokenProvider.clearToken();
	console.log('[API] Dev auth token cleared');
}
