<script lang="ts">
	import { onMount } from 'svelte';
	import ClerkJs from '@clerk/clerk-js';

	let clerk: typeof ClerkJs | null = null;
	let container: HTMLDivElement;
	let error: string | null = null;
	let loading = true;

	onMount(() => {
		let mounted = true;

		async function load() {
			try {
				const publishableKey = import.meta.env.VITE_PUBLIC_CLERK_PUBLISHABLE_KEY;

				if (!publishableKey) {
					if (!mounted) return;
					error = 'Missing Clerk configuration';
					loading = false;
					return;
				}

				// Initialize Clerk
				clerk = new Clerk(publishableKey);
				await clerk.load();

				// Check if already signed in
				if (clerk.user) {
					window.location.href = '/';
					return;
				}

				// Mount Clerk's sign-up component
				if (!mounted) return;
				await clerk.mountSignUp(container, {
					afterSignUpUrl: '/',
					afterSignInUrl: '/',
					signInUrl: '/sign-in'
				});

				if (!mounted) return;
				loading = false;
			} catch (err) {
				if (!mounted) return;
				console.error('[Sign Up] Error:', err);
				error = 'Failed to load sign-up form';
				loading = false;
			}
		}

		load();

		// Cleanup on unmount
		return () => {
			mounted = false;
			if (clerk) {
				try {
					clerk.unmountSignUp(container);
				} catch (error) {
					console.error('[Sign Up] Error unmounting:', error);
				}
			}
		};
	});
</script>

<div class="min-h-screen bg-background-light dark:bg-background-dark flex items-center justify-center px-4">
	<div class="max-w-md w-full">
		<div class="text-center mb-8">
			<h1 class="text-3xl font-display font-bold text-primary dark:text-white mb-2">
				Budget Planner
			</h1>
			<p class="text-gray-600 dark:text-gray-400">Create your account</p>
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
				Already have an account?
				<a href="/sign-in" class="text-blue-600 dark:text-blue-400 hover:underline">Sign in</a>
			</p>
		</div>
	</div>
</div>
