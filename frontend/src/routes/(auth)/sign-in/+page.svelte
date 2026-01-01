<script lang="ts">
  import { onMount, tick } from "svelte";

  let clerk: any = null;
  let container: HTMLDivElement;
  let error: string | null = null;
  let loading = true;

  onMount(() => {
    let mounted = true;

    async function load() {
      try {
        // Dynamic import to prevent SSR issues
        const clerkPkg = await import("@clerk/clerk-js");
        const { Clerk } = clerkPkg;

        const publishableKey = import.meta.env.PUBLIC_CLERK_PUBLISHABLE_KEY;

        // Debug logging
        console.log("[Sign In] Environment check:", {
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
          if (!mounted) return;
          error = "Missing Clerk configuration";
          loading = false;
          return;
        }

        // Initialize Clerk
        clerk = new Clerk(publishableKey);
        await clerk.load();

        // Check if already signed in
        if (clerk.user) {
          window.location.href = "/";
          return;
        }

        // Set loading to false to show container
        if (!mounted) return;
        loading = false;

        // Wait for DOM to update
        await tick();

        // Mount Clerk's sign-in component
        if (!mounted) return;
        await clerk.mountSignIn(container, {
          afterSignInUrl: "/",
          signUpUrl: "/sign-up",
        });
      } catch (err) {
        if (!mounted) return;
        console.error("[Sign In] Error:", err);
        error = "Failed to load sign-in form";
        loading = false;
      }
    }

    load();

    // Cleanup on unmount
    return () => {
      mounted = false;
      if (clerk) {
        try {
          clerk.unmountSignIn(container);
        } catch (error) {
          console.error("[Sign In] Error unmounting:", error);
        }
      }
    };
  });
</script>

<div
  class="min-h-screen bg-background-light dark:bg-background-dark flex items-center justify-center px-4"
>
  <div class="max-w-md w-full">
    <div class="text-center mb-8">
      <h1
        class="text-3xl font-display font-bold text-primary dark:text-white mb-2"
      >
        Budget Planner
      </h1>
      <p class="text-gray-600 dark:text-gray-400">
        Sign in to manage your budget
      </p>
    </div>

    {#if loading}
      <div class="text-center py-12">
        <div
          class="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-r-transparent align-[-0.125em] motion-reduce:animate-[spin_1.5s_linear_infinite]"
        ></div>
        <p class="mt-4 text-gray-600 dark:text-gray-400">Loading...</p>
      </div>
    {:else if error}
      <div
        class="bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200 px-4 py-3 rounded-md text-center"
      >
        <p class="font-medium">Error</p>
        <p class="text-sm">{error}</p>
      </div>
    {:else}
      <div
        bind:this={container}
        class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 min-h-[400px]"
      ></div>
    {/if}

    <div class="mt-6 text-center">
      <p class="text-sm text-gray-600 dark:text-gray-400">
        Don't have an account?
        <a
          href="/sign-up"
          class="text-blue-600 dark:text-blue-400 hover:underline">Sign up</a
        >
      </p>
    </div>
  </div>
</div>
