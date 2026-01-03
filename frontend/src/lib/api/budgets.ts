import { apiClient } from "./client";

/**
 * Budget API types matching backend responses
 */

export interface Budget {
  id: string;
  userId: string;
  name?: string;
  month: string; // Format: YYYY-MM-DD
  totalLimit: number;
  totalIncome?: number;
  savings?: number;
  spent: number;
  remaining: number;
  createdAt: string;
  updatedAt: string;
}

export interface BudgetCategory {
  id: string;
  budgetId: string;
  categoryId: string;
  name: string;
  icon?: string;
  color: string;
  limitAmount: number;
  spent: number;
  remaining: number;
}

export interface CreateBudgetRequest {
  name?: string;
  month: string; // Format: YYYY-MM-DD (first day of month)
  totalLimit: number;
  totalIncome?: number;
}

export interface UpdateBudgetRequest {
  name?: string;
  totalLimit?: number;
  totalIncome?: number;
}

export interface AddBudgetCategoryRequest {
  categoryId: string;
  limitAmount: number;
}

export interface UpdateBudgetCategoryRequest {
  limitAmount?: number;
}

/**
 * Budgets API client
 */
export class BudgetsApi {
  /**
   * Get all budgets for the current user
   * GET /api/budgets
   */
  async getAllBudgets(): Promise<Budget[]> {
    return apiClient.get<Budget[]>("/budgets");
  }

  /**
   * Get budgets for a specific month
   * GET /api/budgets/{month}
   * @param month - Format: YYYY-MM (e.g., "2025-01")
   */
  async getBudgetsByMonth(month: string): Promise<Budget[]> {
    return apiClient.get<Budget[]>(`/budgets/${month}`);
  }

  /**
   * Get a single budget by ID
   * GET /api/budgets/id/{id}
   */
  async getBudgetById(id: string): Promise<Budget> {
    return apiClient.get<Budget>(`/budgets/id/${id}`);
  }

  /**
   * Create a new budget
   * POST /api/budgets
   */
  async createBudget(request: CreateBudgetRequest): Promise<Budget> {
    return apiClient.post<Budget>("/budgets", request);
  }

  /**
   * Update a budget
   * PUT /api/budgets/{id}
   */
  async updateBudget(
    id: string,
    request: UpdateBudgetRequest
  ): Promise<Budget> {
    return apiClient.put<Budget>(`/budgets/${id}`, request);
  }

  /**
   * Delete a budget
   * DELETE /api/budgets/{id}
   */
  async deleteBudget(id: string): Promise<void> {
    return apiClient.delete<void>(`/budgets/${id}`);
  }

  /**
   * Get budget categories
   * GET /api/budgets/{id}/categories
   */
  async getBudgetCategories(budgetId: string): Promise<BudgetCategory[]> {
    return apiClient.get<BudgetCategory[]>(`/budgets/${budgetId}/categories`);
  }

  /**
   * Add a category to a budget
   * POST /api/budgets/{id}/categories
   */
  async addBudgetCategory(
    budgetId: string,
    request: AddBudgetCategoryRequest
  ): Promise<BudgetCategory> {
    return apiClient.post<BudgetCategory>(
      `/budgets/${budgetId}/categories`,
      request
    );
  }
}

// Export singleton instance
export const budgetsApi = new BudgetsApi();
