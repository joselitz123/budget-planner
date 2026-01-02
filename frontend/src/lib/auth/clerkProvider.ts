import type { AuthTokenProvider } from "$lib/api/client";

/**
 * Clerk Authentication Token Provider
 * Retrieves JWT tokens from Clerk session for API requests
 */
export class ClerkTokenProvider implements AuthTokenProvider {
  private clerk: any = null;

  constructor() {
    // Initialize Clerk client only in browser
    if (typeof window !== "undefined") {
      this.initializeClerk();
    }
  }

  private async initializeClerk() {
    try {
      // Dynamic import to prevent SSR issues
      const clerkPkg = await import("@clerk/clerk-js");
      const { Clerk } = clerkPkg;

      const publishableKey = import.meta.env.PUBLIC_CLERK_PUBLISHABLE_KEY;

      console.log("[Clerk] Environment check:", {
        hasKey: !!publishableKey,
        keyPrefix: publishableKey
          ? publishableKey.substring(0, 10) + "..."
          : "N/A",
        envMode: import.meta.env.MODE,
        allEnvVars: Object.keys(import.meta.env).filter((k) =>
          k.startsWith("PUBLIC_")
        ),
      });

      if (!publishableKey) {
        console.error("[Clerk] Missing PUBLIC_CLERK_PUBLISHABLE_KEY");
        return;
      }

      this.clerk = new Clerk(publishableKey);

      // Load Clerk resources
      await this.clerk.load();
      console.log("[Clerk] Successfully initialized");
    } catch (error) {
      console.error("[Clerk] Initialization error:", error);
    }
  }

  /**
   * Get JWT token from Clerk session
   */
  async getToken(): Promise<string | null> {
    try {
      console.log("[Clerk Provider] getToken called");

      if (!this.clerk) {
        console.warn("[Clerk Provider] Clerk client not initialized");
        console.warn(
          "[Clerk Provider] This may mean Clerk is still loading or initialization failed"
        );
        return null;
      }

      console.log("[Clerk Provider] Clerk client is initialized");

      // Get active session
      const session = this.clerk.session;

      if (!session) {
        console.warn("[Clerk Provider] No active session found");
        console.warn(
          "[Clerk Provider] User may not be signed in or session has expired"
        );
        console.log("[Clerk Provider] Clerk state:", {
          hasUser: !!this.clerk.user,
          userId: this.clerk.user?.id || null,
        });
        return null;
      }

      console.log("[Clerk Provider] Active session found, retrieving token...");

      // Get JWT token
      const token = await session.getToken();

      if (!token) {
        console.warn("[Clerk Provider] getToken returned null or empty");
        console.warn(
          "[Clerk Provider] Session exists but token retrieval failed"
        );
        return null;
      }

      console.log("[Clerk Provider] Token retrieved successfully:", {
        tokenLength: token.length,
        tokenPrefix: token.substring(0, 20) + "...",
      });

      return token;
    } catch (error) {
      console.error("[Clerk Provider] Error getting token:", error);
      console.error("[Clerk Provider] Error details:", {
        message: error instanceof Error ? error.message : "Unknown error",
        name: error instanceof Error ? error.name : "Unknown",
        stack: error instanceof Error ? error.stack : undefined,
      });
      return null;
    }
  }
}
