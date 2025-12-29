/**
 * Auth Route Group Server Layout
 *
 * This layout overrides the root layout's authentication requirement
 * for sign-in and sign-up pages.
 *
 * Routes in this group:
 * - /sign-in
 * - /sign-up
 *
 * These routes MUST be accessible without authentication.
 */

import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async (event) => {
	// Get auth state from hooks.server.ts (already verified JWT)
	const { userId, session } = event.locals;

	// DO NOT redirect to sign-in - these pages ARE the sign-in/sign-up pages
	// They must be accessible to unauthenticated users

	return {
		userId,
		session
	};
};
