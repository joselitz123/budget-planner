<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import CustomModal from '$lib/components/ui/CustomModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Select } from '$lib/components/ui/select';
	import { Label } from '$lib/components/ui/label';
	import { isValidEmail } from '$lib/utils/validation';
	import { createInvitation } from '$lib/stores/shares';
	import { cn } from '$lib/utils/cn';

	export let open = false;
	export let budgetId = '';
	export let budgetName = '';

	const dispatch = createEventDispatcher();

	let email = '';
	let permission: 'view' | 'edit' = 'view';
	let submitting = false;
	let error = '';

	function handleClose() {
		open = false;
		error = '';
		email = '';
		permission = 'view';
		dispatch('close');
	}

	async function handleSubmit() {
		// Validate email
		if (!email || !isValidEmail(email)) {
			error = 'Please enter a valid email address';
			return;
		}

		submitting = true;
		error = '';

		try {
			await createInvitation(budgetId, email, permission);
			dispatch('invitationSent', { email, permission });
			handleClose();
		} catch (err) {
			// Error is already handled by createInvitation with toast
			error = 'Failed to send invitation. Please try again.';
		} finally {
			submitting = false;
		}
	}

	function isFormValid(): boolean {
		return isValidEmail(email) && !submitting;
	}
</script>

<CustomModal
	bind:open
	title="Share Budget"
	description="Invite someone to collaborate on this budget"
	on:close={handleClose}
>
	<form onsubmit={handleSubmit} class="space-y-6">
		<!-- Budget Info -->
		<div class="mb-4 p-3 bg-gray-50 dark:bg-gray-900/50 rounded-lg border border-line-light dark:border-line-dark">
			<p class="text-sm text-gray-600 dark:text-gray-400">
				Sharing: <span class="font-semibold text-primary dark:text-white">{budgetName}</span>
			</p>
		</div>

		<!-- Email Input -->
		<div class="space-y-2">
			<Label for="email">Email Address *</Label>
			<Input
				id="email"
				bind:value={email}
				type="email"
				placeholder="friend@example.com"
				required
				disabled={submitting}
				aria-invalid={!!error}
				aria-describedby={error ? 'email-error' : undefined}
				class={cn(error && 'border-red-500 focus:ring-red-500')}
			/>
			{#if error}
				<p id="email-error" class="text-sm text-red-500 dark:text-red-400" role="alert">
					{error}
				</p>
			{/if}
		</div>

		<!-- Permission Dropdown -->
		<div class="space-y-2">
			<Label for="permission">Permission</Label>
			<Select bind:value={permission} disabled={submitting}>
				<option value="view">View Only</option>
				<option value="edit">View & Edit</option>
			</Select>
			<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
				{#if permission === 'view'}
					View Only: Can see budget but not make changes
				{:else}
					View & Edit: Can add, edit, and delete transactions and categories
				{/if}
			</p>
		</div>

		<!-- Actions -->
		<div class="flex justify-end space-x-3 pt-4">
			<Button
				type="button"
				variant="outline"
				onclick={handleClose}
				disabled={submitting}
				aria-label="Cancel sharing"
			>
				Cancel
			</Button>
			<Button
				type="submit"
				disabled={!isFormValid()}
				aria-label={submitting ? 'Sending invitation...' : 'Send invitation'}
			>
				{#if submitting}
					<span class="material-icons-outlined text-sm animate-spin mr-1">sync</span>
					Sending...
				{:else}
					<span class="material-icons-outlined text-sm mr-1">send</span>
					Send Invitation
				{/if}
			</Button>
		</div>
	</form>
</CustomModal>
