import { apiClient } from './client';

/**
 * Sharing API types matching backend responses
 */

export interface ShareInvitation {
	id: string;
	budgetId: string;
	ownerId: string;
	recipientEmail: string;
	permission: 'view' | 'edit';
	status: 'pending' | 'accepted' | 'declined' | 'cancelled' | 'expired';
	expiresAt: string;
	createdAt: string;
	updatedAt: string;
}

export interface ShareAccess {
	id: string;
	budgetId: string;
	ownerId: string;
	sharedWithId: string;
	permission: 'view' | 'edit';
	createdAt: string;
	updatedAt: string;
	// Joined fields from GetShareAccessByBudget:
	sharedWithName?: string;
	sharedWithEmail?: string;
}

export interface SharedBudget {
	// From GetShareAccessForUser with joins
	id: string;
	budgetId: string;
	budgetName?: string;
	budgetMonth?: string;
	ownerId: string;
	ownerName?: string;
	ownerEmail?: string;
	sharedWithId: string;
	permission: 'view' | 'edit';
	createdAt: string;
	updatedAt: string;
}

export interface CreateInvitationRequest {
	budgetId: string;
	recipientEmail: string;
	permission: 'view' | 'edit';
}

export interface RespondInvitationRequest {
	status: 'accepted' | 'declined';
}

/**
 * Shares API client
 *
 * Handles all budget sharing operations:
 * - Create and manage share invitations
 * - Accept/decline invitations
 * - View who has access to budgets
 * - Manage access permissions
 */
export class SharesApi {
	/**
	 * Create a new share invitation
	 * POST /api/sharing/invite
	 */
	async createInvitation(request: CreateInvitationRequest): Promise<ShareInvitation> {
		return apiClient.post<ShareInvitation>('/sharing/invite', request);
	}

	/**
	 * Get pending invitations for current user
	 * GET /api/sharing/invitations
	 */
	async getMyInvitations(): Promise<ShareInvitation[]> {
		return apiClient.get<ShareInvitation[]>('/sharing/invitations');
	}

	/**
	 * Respond to a share invitation (accept/decline)
	 * PUT /api/sharing/invitations/{id}/respond
	 */
	async respondToInvitation(
		id: string,
		request: RespondInvitationRequest
	): Promise<ShareInvitation> {
		return apiClient.put<ShareInvitation>(`/sharing/invitations/${id}/respond`, request);
	}

	/**
	 * Cancel a pending invitation
	 * DELETE /api/sharing/invitations/{id}
	 */
	async cancelInvitation(id: string): Promise<{ message: string }> {
		return apiClient.delete<{ message: string }>(`/sharing/invitations/${id}`);
	}

	/**
	 * Get who has access to a specific budget
	 * GET /api/sharing/budgets/{budgetId}
	 */
	async getBudgetSharing(budgetId: string): Promise<ShareAccess[]> {
		return apiClient.get<ShareAccess[]>(`/sharing/budgets/${budgetId}`);
	}

	/**
	 * Remove someone's access to a budget
	 * DELETE /api/sharing/access/{id}
	 */
	async removeAccess(id: string): Promise<{ message: string }> {
		return apiClient.delete<{ message: string }>(`/sharing/access/${id}`);
	}

	/**
	 * Get budgets shared with current user
	 * GET /api/sharing/shared-with-me
	 */
	async getSharedBudgets(): Promise<SharedBudget[]> {
		return apiClient.get<SharedBudget[]>('/sharing/shared-with-me');
	}
}

// Export singleton instance
export const sharesApi = new SharesApi();
