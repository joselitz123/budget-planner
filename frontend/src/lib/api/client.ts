import { ClerkTokenProvider } from "$lib/auth/clerkProvider";
import type { ApiResponse } from "$lib/db/schema";
import { showToast } from "$lib/stores/ui";

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
    this.baseUrl =
      baseUrl || import.meta.env.PUBLIC_API_URL || "http://localhost:8080/api";
    this.authProvider = authProvider || devTokenProvider;
  }

  /**
   * Get auth headers including JWT token
   */
  private async getAuthHeaders(): Promise<Record<string, string>> {
    console.log("[API Client] getAuthHeaders called");
    console.log(
      "[API Client] Auth provider type:",
      this.authProvider.constructor.name
    );

    const token = await this.authProvider.getToken();
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    console.log("[API Client] Token retrieval result:", {
      hasToken: !!token,
      tokenLength: token?.length || 0,
      tokenPrefix: token ? token.substring(0, 20) + "..." : "N/A",
      tokenType: typeof token,
    });

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
      console.log("[API Client] Authorization header set successfully:", {
        headerLength: headers["Authorization"].length,
        headerPrefix: headers["Authorization"].substring(0, 30) + "...",
      });
    } else {
      console.warn(
        "[API Client] No token available - request may fail with 401"
      );
      console.warn(
        "[API Client] This usually means the user is not authenticated or the session has expired"
      );
    }

    console.log("[API Client] Final headers being returned:", {
      keys: Object.keys(headers),
      hasAuth: !!headers["Authorization"],
      hasContentType: !!headers["Content-Type"],
    });

    return headers;
  }

  /**
   * Handle API errors with user-friendly messages
   */
  private handleApiError(response: Response, data?: any): never {
    let message = "";

    // Handle 401 Unauthorized - auth token expired or invalid
    if (response.status === 401) {
      message = "Session expired. Please log in again.";
      console.error("[API] Unauthorized - please log in again");
      showToast(message, "error", 5000);
      throw new Error(message);
    }

    // Handle 403 Forbidden
    if (response.status === 403) {
      message = "You do not have permission to perform this action.";
      console.error("[API] Forbidden - insufficient permissions");
      showToast(message, "error", 5000);
      throw new Error(message);
    }

    // Handle 404 Not Found
    if (response.status === 404) {
      message = "The requested resource was not found.";
      console.error("[API] Resource not found");
      showToast(message, "error", 3000);
      throw new Error(message);
    }

    // Handle 500 Server Error
    if (response.status >= 500) {
      message = "Server error. Please try again later.";
      console.error("[API] Server error:", response.statusText);
      showToast(message, "error", 5000);
      throw new Error(message);
    }

    // Handle other errors
    message =
      data?.error?.message || response.statusText || "API request failed";
    console.error("[API] Error:", message);
    showToast(message, "error", 3000);
    throw new Error(message);
  }

  /**
   * Make a GET request
   */
  async get<T>(path: string, options?: RequestInit): Promise<T> {
    const headers = await this.getAuthHeaders();

    console.log("[API Client] Making GET request:", {
      url: `${this.baseUrl}${path}`,
      hasAuthHeader: !!headers["Authorization"],
      authHeaderLength: headers["Authorization"]?.length || 0,
      authHeaderPrefix:
        headers["Authorization"]?.substring(0, 30) + "..." || "N/A",
      allHeaderKeys: Object.keys(headers),
    });

    // Log the complete request for debugging
    console.log("[API Client] Complete request details:", {
      url: `${this.baseUrl}${path}`,
      method: "GET",
      headers: headers,
      hasAuthHeader: !!headers["Authorization"],
      authHeaderValue: headers["Authorization"]
        ? `${headers["Authorization"].substring(0, 20)}...`
        : "none",
    });

    const response = await fetch(`${this.baseUrl}${path}`, {
      ...options,
      method: "GET",
      headers: {
        ...headers,
        ...options?.headers,
      },
    });

    console.log("[API Client] GET response:", {
      status: response.status,
      ok: response.ok,
      statusText: response.statusText,
      url: `${this.baseUrl}${path}`,
    });

    const data = await response.json().catch(() => ({}));

    if (!response.ok) {
      this.handleApiError(response, data);
    }

    const result: ApiResponse<T> = data;

    if (!result.success) {
      throw new Error(result.error?.message || "API request failed");
    }

    return result.data;
  }

  /**
   * Make a POST request
   */
  async post<T>(path: string, data?: any, options?: RequestInit): Promise<T> {
    const headers = await this.getAuthHeaders();

    console.log("[API Client] Making POST request:", {
      url: `${this.baseUrl}${path}`,
      hasAuthHeader: !!headers["Authorization"],
      authHeaderLength: headers["Authorization"]?.length || 0,
      authHeaderPrefix:
        headers["Authorization"]?.substring(0, 30) + "..." || "N/A",
      hasBody: !!data,
    });

    const response = await fetch(`${this.baseUrl}${path}`, {
      ...options,
      method: "POST",
      headers: {
        ...headers,
        ...options?.headers,
      },
      body: JSON.stringify(data),
    });

    console.log("[API Client] POST response:", {
      status: response.status,
      ok: response.ok,
      statusText: response.statusText,
      url: `${this.baseUrl}${path}`,
    });

    const responseData = await response.json().catch(() => ({}));

    if (!response.ok) {
      this.handleApiError(response, responseData);
    }

    const result: ApiResponse<T> = responseData;

    if (!result.success) {
      throw new Error(result.error?.message || "API request failed");
    }

    return result.data;
  }

  /**
   * Make a PUT request
   */
  async put<T>(path: string, data?: any, options?: RequestInit): Promise<T> {
    const headers = await this.getAuthHeaders();

    console.log("[API Client] Making PUT request:", {
      url: `${this.baseUrl}${path}`,
      hasAuthHeader: !!headers["Authorization"],
      authHeaderLength: headers["Authorization"]?.length || 0,
      authHeaderPrefix:
        headers["Authorization"]?.substring(0, 30) + "..." || "N/A",
      hasBody: !!data,
    });

    const response = await fetch(`${this.baseUrl}${path}`, {
      ...options,
      method: "PUT",
      headers: {
        ...headers,
        ...options?.headers,
      },
      body: JSON.stringify(data),
    });

    console.log("[API Client] PUT response:", {
      status: response.status,
      ok: response.ok,
      statusText: response.statusText,
      url: `${this.baseUrl}${path}`,
    });

    const responseData = await response.json().catch(() => ({}));

    if (!response.ok) {
      this.handleApiError(response, responseData);
    }

    const result: ApiResponse<T> = responseData;

    if (!result.success) {
      throw new Error(result.error?.message || "API request failed");
    }

    return result.data;
  }

  /**
   * Make a DELETE request
   */
  async delete<T>(path: string, options?: RequestInit): Promise<T> {
    const headers = await this.getAuthHeaders();

    console.log("[API Client] Making DELETE request:", {
      url: `${this.baseUrl}${path}`,
      hasAuthHeader: !!headers["Authorization"],
      authHeaderLength: headers["Authorization"]?.length || 0,
      authHeaderPrefix:
        headers["Authorization"]?.substring(0, 30) + "..." || "N/A",
    });

    const response = await fetch(`${this.baseUrl}${path}`, {
      ...options,
      method: "DELETE",
      headers: {
        ...headers,
        ...options?.headers,
      },
    });

    console.log("[API Client] DELETE response:", {
      status: response.status,
      ok: response.ok,
      statusText: response.statusText,
      url: `${this.baseUrl}${path}`,
    });

    const data = await response.json().catch(() => ({}));

    if (!response.ok) {
      this.handleApiError(response, data);
    }

    const result: ApiResponse<T> = data;

    if (!result.success) {
      throw new Error(result.error?.message || "API request failed");
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
// Development switch: use dev provider when VITE_DEV_AUTH_PROVIDER=dev
// Use ClerkTokenProvider when set to 'clerk' or unset
const authProvider =
  import.meta.env.PUBLIC_DEV_AUTH_PROVIDER === "dev"
    ? devTokenProvider
    : new ClerkTokenProvider();

export const apiClient = new ApiClient(undefined, authProvider);

/**
 * Helper function to set dev JWT token for testing
 * Call this with a valid JWT token to authenticate API requests
 */
export function setDevAuthToken(token: string) {
  devTokenProvider.setToken(token);
  console.log("[API] Dev auth token set");
}

/**
 * Clear dev auth token
 */
export function clearDevAuthToken() {
  devTokenProvider.clearToken();
  console.log("[API] Dev auth token cleared");
}
