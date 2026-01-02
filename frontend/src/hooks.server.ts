import { env as privateEnv } from "$env/dynamic/private";
import { env as publicEnv } from "$env/dynamic/public";
import { createClerkClient } from "@clerk/backend";
import type { Handle } from "@sveltejs/kit";
import { sequence } from "@sveltejs/kit/hooks";

// Extend App.Locals interface
declare global {
  namespace App {
    interface Locals {
      userId: string | null;
      session: any;
    }
  }
}

// Initialize Clerk backend client
const clerkClient = createClerkClient({
  secretKey: privateEnv.CLERK_SECRET_KEY,
  publishableKey: publicEnv.PUBLIC_CLERK_PUBLISHABLE_KEY,
});

/**
 * Clerk authentication middleware
 * Uses Clerk's backend SDK to verify session and extract authenticated user
 */
const clerkAuth: Handle = async ({ event, resolve }) => {
  const sessionToken = event.cookies.get("__session");

  console.log("[Clerk] Session token check:", {
    hasToken: !!sessionToken,
    tokenLength: sessionToken?.length || 0,
    cookieNames: Object.keys(event.cookies.getAll()),
  });

  let userId: string | null = null;
  let session: any = null;

  // Verify session using Clerk backend SDK with networkless verification
  if (sessionToken) {
    try {
      console.log(
        "[Clerk] Attempting to verify session using networkless verification..."
      );
      console.log(
        "[Clerk] CLERK_SECRET_KEY present:",
        !!privateEnv.CLERK_SECRET_KEY
      );

      // Use Clerk's authenticateRequest for networkless verification
      // This extracts session from request without making network calls
      const requestState = await clerkClient.authenticateRequest(event.request);

      console.log("[Clerk] Request state:", {
        isSignedIn: requestState.isSignedIn,
        status: requestState.status,
        reason: requestState.reason,
      });

      if (requestState.isSignedIn) {
        // Get auth object from request state
        const authObject = requestState.toAuth();

        if (authObject) {
          // Extract session from auth object
          session = authObject;

          // Extract user ID from session
          userId = session.userId || session.claims?.sub || null;

          console.log("[Clerk] Authenticated user:", userId);
          console.log("[Clerk] Session details:", {
            userId: session.userId,
            hasClaims: !!session.claims,
            sessionId: session.id,
          });
        } else {
          console.log(
            "[Clerk] Request state shows signed in but auth object is null"
          );
        }
      } else {
        console.log("[Clerk] User not signed in - request state:", {
          status: requestState.status,
          reason: requestState.reason,
        });
      }
    } catch (error) {
      // Session is invalid, expired, or verification failed
      console.error("[Clerk] Session verification failed:", error);
      console.error("[Clerk] Error details:", {
        message: error instanceof Error ? error.message : "Unknown error",
        name: error instanceof Error ? error.name : "Unknown",
        stack: error instanceof Error ? error.stack : undefined,
      });
      // Treat as unauthenticated
      userId = null;
      session = null;
    }
  } else {
    console.log("[Clerk] No session token found in cookies");
  }

  // Inject auth state into event.locals
  event.locals.userId = userId;
  event.locals.session = session;

  console.log("[Clerk] Auth state injected:", {
    userId,
    hasSession: !!session,
    path: event.url.pathname,
  });

  // Continue to next handler or page
  return resolve(event);
};

/**
 * Main handle function with middleware sequence
 */
export const handle = sequence(clerkAuth);
