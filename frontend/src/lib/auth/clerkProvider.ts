import type { AuthTokenProvider } from '$lib/api/client';
import clerkPkg from '@clerk/clerk-js';
const { Clerk } = clerkPkg;

/**
 * Clerk Authentication Token Provider
 * Retrieves JWT tokens from Clerk session for API requests
 */
export class ClerkTokenProvider implements AuthTokenProvider {
	private clerk: any = null;

	constructor() {
		// Initialize Clerk client
		if (typeof window !== 'undefined') {
			try {
				const publishableKey = import.meta.env.VITE_PUBLIC_CLERK_PUBLISHABLE_KEY;

				if (!publishableKey) {
					console.error('[Clerk] Missing PUBLIC_CLERK_PUBLISHABLE_KEY');
					return;
				}

				this.clerk = new Clerk(publishableKey);

				// Load Clerk resources
				this.clerk.load().catch((error: unknown) => {
					console.error('[Clerk] Failed to load:', error);
				});
			} catch (error) {
				console.error('[Clerk] Initialization error:', error);
			}
		}
	}

	/**
	 * Get JWT token from Clerk session
	 */
	async getToken(): Promise<string | null> {
		try {
			if (!this.clerk) {
				return null;
			}

			// Get active session
			const session = this.clerk.session;

			if (!session) {
				return null;
			}

			// Get JWT token
			const token = await session.getToken();

			return token || null;
		} catch (error) {
			console.error('[Clerk] Error getting token:', error);
			return null;
		}
	}
}
