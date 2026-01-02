import { env } from "$env/dynamic/private";
import { redirect } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async (event) => {
  // Get auth state from hooks.server.ts
  const { userId, session } = event.locals;

  console.log("[Layout] Load called:", {
    path: event.url.pathname,
    hasUserId: !!userId,
    userId: userId,
    hasSession: !!session,
  });

  // Skip auth check for auth routes (sign-in, sign-up)
  const url = event.url.pathname;
  if (url.startsWith("/sign-in") || url.startsWith("/sign-up")) {
    console.log("[Layout] Skipping auth check for auth route:", url);
    return {
      userId,
      hasSession: !!session,
    };
  }

  // Redirect to sign-in if not authenticated
  if (!userId) {
    console.log("[Layout] No userId found, redirecting to /sign-in");
    throw redirect(302, "/sign-in");
  }

  console.log("[Layout] User authenticated, proceeding to fetch user data");

  // Onboarding flow: Try to fetch user from backend
  let user: any = null;

  try {
    // Try to get existing user from backend
    // Note: We use event.fetch to include auth headers
    const apiUrl = env.API_URL || "http://localhost:8080/api";
    console.log("[Layout] Fetching user from backend:", {
      apiUrl: `${apiUrl}/users/me`,
      hasApiUrl: !!env.API_URL,
    });

    // Get JWT token from Clerk session for backend authentication
    const token = session ? await session.getToken() : null;

    console.log("[Layout] Fetching user from backend:", {
      apiUrl: `${apiUrl}/users/me`,
      hasApiUrl: !!env.API_URL,
      hasToken: !!token,
      tokenLength: token?.length || 0,
    });

    const headers: HeadersInit = {
      "Content-Type": "application/json",
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
      console.log("[Layout] Authorization header added");
    } else {
      console.warn("[Layout] No token available - request may fail with 401");
    }

    const response = await event.fetch(`${apiUrl}/users/me`, {
      headers,
    });

    console.log("[Layout] Backend response:", {
      status: response.status,
      ok: response.ok,
      statusText: response.statusText,
    });

    if (response.ok) {
      user = await response.json();
      console.log("[Layout] User fetched successfully:", user);
    } else if (response.status === 404) {
      // User doesn't exist in backend - call onboarding
      console.log("[Auth] User not found in backend, calling onboarding...");

      // Get user data from Clerk session
      const userData = session?.user || {};
      const firstName = userData.firstName || userData.username || "User";

      // Call onboarding endpoint
      // Get JWT token for onboarding request
      const onboardingToken = session ? await session.getToken() : null;

      const onboardingHeaders: HeadersInit = {
        "Content-Type": "application/json",
      };

      if (onboardingToken) {
        onboardingHeaders["Authorization"] = `Bearer ${onboardingToken}`;
        console.log("[Layout] Authorization header added for onboarding");
      } else {
        console.warn(
          "[Layout] No token available for onboarding - request may fail with 401"
        );
      }

      const onboardingResponse = await event.fetch(
        `${apiUrl}/auth/onboarding`,
        {
          method: "POST",
          headers: onboardingHeaders,
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
    } else {
      console.error("[Layout] Backend returned non-OK response:", {
        status: response.status,
        statusText: response.statusText,
        body: await response.text().catch(() => "Could not read body"),
      });
    }
  } catch (error) {
    console.error("[Auth] Error fetching user data:", error);
    console.error("[Auth] Error details:", {
      message: error instanceof Error ? error.message : "Unknown error",
      name: error instanceof Error ? error.name : "Unknown",
      stack: error instanceof Error ? error.stack : undefined,
    });
    // Continue without user data - app should handle this gracefully
  }

  console.log("[Layout] Returning data:", {
    userId,
    hasSession: !!session,
    hasUser: !!user,
  });

  return {
    userId,
    hasSession: !!session,
    user,
  };
};
