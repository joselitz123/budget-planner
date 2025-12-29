<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import type { ShareInvitation } from '$lib/api/shares';
	import { respondToInvitation } from '$lib/stores/shares';

	export let invitations: ShareInvitation[] = [];
	export let loading = false;

	const dispatch = createEventDispatcher();
	let responding = new Set<string>();

	async function handleAccept(invitation: ShareInvitation) {
		responding.add(invitation.id);
		try {
			await respondToInvitation(invitation.id, true);
			dispatch('respond', { id: invitation.id, accepted: true });
		} catch (err) {
			// Error already handled by respondToInvitation
			console.error('Failed to accept invitation:', err);
		} finally {
			responding.delete(invitation.id);
		}
	}

	async function handleDecline(invitation: ShareInvitation) {
		responding.add(invitation.id);
		try {
			await respondToInvitation(invitation.id, false);
			dispatch('respond', { id: invitation.id, accepted: false });
		} catch (err) {
			// Error already handled by respondToInvitation
			console.error('Failed to decline invitation:', err);
		} finally {
			responding.delete(invitation.id);
		}
	}

	function formatDate(dateString: string): string {
		return new Intl.DateTimeFormat('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		}).format(new Date(dateString));
	}

	function getPermissionVariant(permission: string) {
		return permission === 'edit' ? 'default' : 'secondary';
	}
</script>

<section class="space-y-4">
	<h3 class="text-lg font-display font-bold text-primary dark:text-white">Pending Invitations</h3>

	{#if loading}
		<div class="space-y-3">
			<Skeleton variant="card" className="h-24 w-full" />
			<Skeleton variant="card" className="h-24 w-full" />
		</div>
	{:else if invitations.length === 0}
		<p
			class="text-gray-500 dark:text-gray-400 text-center py-8 italic bg-paper-light dark:bg-paper-dark rounded-lg border border-line-light dark:border-line-dark"
		>
			No pending invitations
		</p>
	{:else}
		<div class="space-y-3">
			{#each invitations as invitation}
				<div
					class="bg-paper-light dark:bg-paper-dark border border-line-light dark:border-line-dark rounded-lg p-4 shadow-sm"
				>
					<div class="flex justify-between items-start gap-4">
						<div class="flex-1 min-w-0">
							<!-- Email -->
							<p class="font-medium text-primary dark:text-white flex items-center gap-2">
								<span class="material-icons-outlined text-sm text-gray-500">email</span>
								<span class="truncate">{invitation.recipientEmail}</span>
							</p>

							<!-- Permission & Expiry -->
							<div class="flex flex-wrap items-center gap-3 mt-2 text-sm">
								<div class="flex items-center gap-1">
									<span class="text-gray-500 dark:text-gray-400">Permission:</span>
									<Badge variant={getPermissionVariant(invitation.permission)}>
										{invitation.permission}
									</Badge>
								</div>
								<div class="flex items-center gap-1 text-gray-500 dark:text-gray-400">
									<span class="material-icons-outlined text-sm">schedule</span>
									<span>Expires: {formatDate(invitation.expiresAt)}</span>
								</div>
							</div>
						</div>

						<!-- Actions -->
						<div class="flex flex-col sm:flex-row gap-2 shrink-0">
							<Button
								size="sm"
								variant="destructive"
								onclick={() => handleDecline(invitation)}
								disabled={responding.has(invitation.id)}
								aria-label="Decline invitation from {invitation.recipientEmail}"
							>
								{#if responding.has(invitation.id)}
									<span class="material-icons-outlined text-sm animate-spin">sync</span>
								{:else}
									Decline
								{/if}
							</Button>
							<Button
								size="sm"
								onclick={() => handleAccept(invitation)}
								disabled={responding.has(invitation.id)}
								aria-label="Accept invitation from {invitation.recipientEmail}"
							>
								{#if responding.has(invitation.id)}
									<span class="material-icons-outlined text-sm animate-spin">sync</span>
								{:else}
									Accept
								{/if}
							</Button>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</section>
