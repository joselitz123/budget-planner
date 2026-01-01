import { redirect } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async (event) => {
  // Get auth state from hooks.server.ts
  const { userId, session } = event.locals;

  // Skip auth check for auth routes (sign-in, sign-up)
  const url = event.url.pathname;
  if (url.startsWith("/sign-in") || url.startsWith("/sign-up")) {
    return {
      userId,
      session,
    };
  }

  // Redirect to sign-in if not authenticated
  if (!userId) {
    throw redirect(302, "/sign-in");
  }

  // Onboarding flow: Try to fetch user from backend
  let user: any = null;

  try {
    // Try to get existing user from backend
    // Note: We use event.fetch to include auth headers
    const response = await event.fetch(`${process.env.API_URL}/users/me`, {
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (response.ok) {
      user = await response.json();
    } else if (response.status === 404) {
      // User doesn't exist in backend - call onboarding
      console.log("[Auth] User not found in backend, calling onboarding...");

      // Get user data from Clerk session
      const userData = session?.user || {};
      const firstName = userData.firstName || userData.username || "User";

      // Call onboarding endpoint
      const onboardingResponse = await event.fetch(
        `${process.env.API_URL}/auth/onboarding`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            clerkUserId: userId,
            name: firstName,
            currency: "PHP",
          }),
        }
      );

      if (onboardingResponse.ok) {
        user = await onboardingResponse.json();
        console.log("[Auth] Onboarding successful for user:", userId);
      } else {
        console.error(
          "[Auth] Onboarding failed:",
          await onboardingResponse.text()
        );
      }
    }
  } catch (error) {
    console.error("[Auth] Error fetching user data:", error);
    // Continue without user data - app should handle this gracefully
  }

  return {
    userId,
    session,
    user,
  };
};
