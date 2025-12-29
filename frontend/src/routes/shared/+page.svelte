<script lang="ts">
	import { onMount } from 'svelte';
	import { sharedBudgets, sharedBudgetsLoading, pendingInvitations, invitationsLoading } from '$lib/stores/shares';
	import { loadInvitations, loadSharedBudgets } from '$lib/stores/shares';
	import InvitationList from '$lib/components/sharing/InvitationList.svelte';
	import SharedBudgetCard from '$lib/components/sharing/SharedBudgetCard.svelte';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { showToast } from '$lib/stores/ui';

	onMount(async () => {
		try {
			await Promise.all([loadInvitations(), loadSharedBudgets()]);
		} catch (error) {
			console.error('Failed to load shared budgets:', error);
			showToast('Failed to load shared budgets', 'error');
		}
	});
</script>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="mb-6">
		<h2 class="text-3xl font-display font-bold text-primary dark:text-white mb-2">
			Shared With Me
		</h2>
		<p class="text-sm text-gray-500 dark:text-gray-400 italic">
			"Budgets others have shared with you."
		</p>
	</div>

	<!-- Pending Invitations Section -->
	{#if $pendingInvitations.length > 0}
		<section
			class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-5"
		>
			<InvitationList invitations={$pendingInvitations} loading={$invitationsLoading} />
		</section>
	{/if}

	<!-- Shared Budgets Grid -->
	<section>
		<h3 class="text-lg font-display font-bold text-primary dark:text-white mb-4">
			Shared Budgets
		</h3>

		{#if $sharedBudgetsLoading}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				<Skeleton variant="card" className="h-32 w-full" />
				<Skeleton variant="card" className="h-32 w-full" />
				<Skeleton variant="card" className="h-32 w-full" />
			</div>
		{:else if $sharedBudgets.length === 0}
			<div
				class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-8 text-center"
			>
				<span class="material-icons-outlined text-6xl text-gray-300 dark:text-gray-600 mb-4">
					folder_off
				</span>
				<h3 class="text-xl font-display font-bold text-primary dark:text-white mb-2">
					No Shared Budgets
				</h3>
				<p class="text-gray-500 dark:text-gray-400">
					When someone shares a budget with you, it will appear here.
				</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				{#each $sharedBudgets as budget}
					<SharedBudgetCard sharedBudget={budget} />
				{/each}
			</div>
		{/if}
	</section>
</div>
