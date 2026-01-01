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
      console.error("[Clerk] Error getting token:", error);
      return null;
    }
  }
}
