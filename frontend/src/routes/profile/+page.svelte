<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { Clerk } from '@clerk/clerk-js';

	let clerk: Clerk | null = null;
	let loading = true;

	onMount(async () => {
		try {
			const publishableKey = import.meta.env.VITE_PUBLIC_CLERK_PUBLISHABLE_KEY;

			if (publishableKey) {
				clerk = new Clerk(publishableKey);
				await clerk.load();
			}
		} catch (error) {
			console.error('[Profile] Error initializing Clerk:', error);
		}

		loading = false;
	});

	async function handleLogout() {
		try {
			if (clerk) {
				await clerk.signOut();
				window.location.href = '/sign-in';
			}
		} catch (error) {
			console.error('[Profile] Error signing out:', error);
		}
	}
</script>

<div class="max-w-2xl mx-auto">
	<Card>
		<div class="p-6">
			<h1 class="text-2xl font-display font-bold mb-6 text-primary dark:text-white">
				Profile
			</h1>

			{#if loading}
				<div class="text-center py-12">
					<div
						class="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-r-transparent align-[-0.125em] motion-reduce:animate-[spin_1.5s_linear_infinite]"
					></div>
					<p class="mt-4 text-gray-600 dark:text-gray-400">Loading...</p>
				</div>
			{:else}
				{#if $page.data.user}
					<div class="space-y-4">
						<div>
							<label class="text-sm font-medium text-gray-600 dark:text-gray-400"
								>Name</label
							>
							<p class="text-lg text-gray-900 dark:text-white">
								{$page.data.user.name || 'N/A'}
							</p>
						</div>
						<div>
							<label class="text-sm font-medium text-gray-600 dark:text-gray-400"
								>Email</label
							>
							<p class="text-lg text-gray-900 dark:text-white">
								{$page.data.user.email || 'N/A'}
							</p>
						</div>
						<div>
							<label class="text-sm font-medium text-gray-600 dark:text-gray-400"
								>Currency</label
							>
							<p class="text-lg text-gray-900 dark:text-white">
								{$page.data.user.currency || 'PHP'}
							</p>
						</div>
					</div>
				{:else}
					<div
						class="bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200 px-4 py-3 rounded-md"
					>
						<p class="font-medium">User information not available</p>
						<p class="text-sm mt-1">
							Please complete the onboarding process or try refreshing the page.
						</p>
					</div>
				{/if}

				<div class="mt-8 pt-6 border-t border-gray-200 dark:border-gray-700">
					<Button
						onclick={handleLogout}
						variant="destructive"
						class="w-full sm:w-auto"
						disabled={loading}
					>
						Sign Out
					</Button>
				</div>
			{/if}
		</div>
	</Card>
</div>
