import { sequence } from '@sveltejs/kit/hooks';
import type { Handle } from '@sveltejs/kit';
import { verifyToken } from '@clerk/backend';

// Extend App.Locals interface
declare global {
	namespace App {
		interface Locals {
			userId: string | null;
			session: any;
		}
	}
}

/**
 * Clerk authentication middleware
 * Verifies JWT signatures and extracts authenticated user
 */
const clerkAuth: Handle = async ({ event, resolve }) => {
	// Get session token from cookie (Clerk uses __session cookie by default)
	const sessionToken = event.cookies.get('__session');

	let userId: string | null = null;
	let session: any = null;

	// Verify JWT signature with Clerk backend SDK
	if (sessionToken) {
		try {
			// Verify token signature and expiration
			// This throws if token is invalid, expired, or signature is wrong
			const verifiedToken = await verifyToken(sessionToken, {
				secretKey: process.env.CLERK_SECRET_KEY
			});

			// Extract user ID from verified token
			userId = verifiedToken?.sub || null;
			session = verifiedToken;

			console.log('[Clerk] Authenticated user:', userId);
		} catch (error) {
			// Token is invalid, expired, or signature verification failed
			console.error('[Clerk] Token verification failed:', error);
			// Treat as unauthenticated
			userId = null;
			session = null;
		}
	}

	// Inject auth state into event.locals
	event.locals.userId = userId;
	event.locals.session = session;

	// Continue to next handler or page
	return resolve(event);
};

/**
 * Main handle function with middleware sequence
 */
export const handle = sequence(clerkAuth);
