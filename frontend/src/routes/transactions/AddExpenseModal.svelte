<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import CustomModal from '$lib/components/ui/CustomModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Select } from '$lib/components/ui/select';
	import { allCategories } from '$lib/stores/categories';
	import { currentBudget } from '$lib/stores/budgets';
	import { createBudgetForCurrentMonth } from '$lib/stores/budgets';
	import { activePaymentMethods } from '$lib/stores/paymentMethods';
	import PaymentMethodSelector from '$lib/components/payment/PaymentMethodSelector.svelte';
	import type { Transaction } from '$lib/db/schema';
	import { isValidAmount, isRequired } from '$lib/utils/validation';

	export let open = false;
	export let userId = 'temp-user'; // TODO: Get from Clerk auth

	const dispatch = createEventDispatcher();

	// Form state
	let amount = '';
	let date = new Date().toISOString().split('T')[0];
	let description = '';
	let categoryId = $allCategories[0]?.id || '';
	let isRecurring = false;
	let dueDate = '';
	let notes = '';
	let paymentMethodId: string | null = null;

	// Form errors
	let errors = {
		amount: '',
		date: '',
		description: '',
		category: ''
	};

	function closeModal() {
		open = false;
		resetForm();
		dispatch('close');
	}

	function resetForm() {
		amount = '';
		date = new Date().toISOString().split('T')[0];
		description = '';
		categoryId = $allCategories[0]?.id || '';
		isRecurring = false;
		dueDate = '';
		notes = '';
		paymentMethodId = null;
		errors = { amount: '', date: '', description: '', category: '' };
	}

	function validateForm(): boolean {
		let valid = true;
		errors = { amount: '', date: '', description: '', category: '' };

		// Validate amount
		if (!isRequired(amount)) {
			errors.amount = 'Amount is required';
			valid = false;
		} else if (!isValidAmount(amount)) {
			errors.amount = 'Amount must be a positive number';
			valid = false;
		}

		// Validate date
		if (!isRequired(date)) {
			errors.date = 'Date is required';
			valid = false;
		}

		// Validate description
		if (!isRequired(description)) {
			errors.description = 'Description is required';
			valid = false;
		} else if (description.length > 255) {
			errors.description = 'Description must be 255 characters or less';
			valid = false;
		}

		// Validate category
		if (!isRequired(categoryId)) {
			errors.category = 'Category is required';
			valid = false;
		}

		return valid;
	}

	async function handleSubmit() {
		if (!validateForm()) {
			return;
		}

		try {
			// Ensure budget exists for current month
			let budget = $currentBudget;
			if (!budget) {
				budget = await createBudgetForCurrentMonth(userId);
			}

			// Create transaction object
			const transaction: Transaction = {
				id: crypto.randomUUID(),
				userId,
				budgetId: budget.id,
				categoryId,
				amount: parseFloat(amount),
				description,
				transactionDate: date,
				transactionType: 'expense',
				paymentMethodId: paymentMethodId,
				paid: false,
				dueDate: isRecurring && dueDate ? dueDate : null,
				isRecurring,
				notes: notes || null,
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString()
			};

			// Dispatch submit event with transaction data
			dispatch('submit', transaction);
			closeModal();
		} catch (error) {
			console.error('[AddExpenseModal] Error creating transaction:', error);
		}
	}
</script>

<CustomModal
	bind:open
	title="Add Expense"
	description="Enter the details of your expense"
	className="max-w-2xl"
>
	<form on:submit|preventDefault={handleSubmit} class="space-y-6">
		<!-- Amount -->
		<div class="space-y-2">
			<Label for="amount">Amount *</Label>
			<Input
				id="amount"
				type="number"
				step="0.01"
				min="0"
				bind:value={amount}
				placeholder="0.00"
				class="text-2xl font-bold"
				required
			/>
			{#if errors.amount}
				<p class="text-sm text-red-500 dark:text-red-400">{errors.amount}</p>
			{/if}
		</div>

		<!-- Date -->
		<div class="space-y-2">
			<Label for="date">Date *</Label>
			<Input id="date" type="date" bind:value={date} required />
			{#if errors.date}
				<p class="text-sm text-red-500 dark:text-red-400">{errors.date}</p>
			{/if}
		</div>

		<!-- Description -->
		<div class="space-y-2">
			<Label for="description">Description *</Label>
			<Input
				id="description"
				bind:value={description}
				placeholder="What did you spend on?"
				class="font-handwriting text-lg"
				maxlength="255"
				required
			/>
			{#if errors.description}
				<p class="text-sm text-red-500 dark:text-red-400">{errors.description}</p>
			{/if}
		</div>

		<!-- Category -->
		<div class="space-y-2">
			<Label for="category">Category *</Label>
			<Select bind:value={categoryId} required>
				<option value="">Select a category</option>
				{#each $allCategories as category}
					<option value={category.id}>
						{category.icon} {category.name}
					</option>
				{/each}
			</Select>
			{#if errors.category}
				<p class="text-sm text-red-500 dark:text-red-400">{errors.category}</p>
			{/if}
		</div>

		<!-- Payment Method -->
		<div class="space-y-2">
			<Label for="paymentMethod">Payment Method</Label>
			<PaymentMethodSelector
				bind:selectedId={paymentMethodId}
				paymentMethods={$activePaymentMethods}
				onChange={(id) => (paymentMethodId = id)}
			/>
		</div>

		<!-- Is Recurring -->
		<div class="flex items-center space-x-2">
			<input type="checkbox" id="recurring" bind:checked={isRecurring} class="w-4 h-4" />
			<Label for="recurring" class="cursor-pointer">This is a recurring expense</Label>
		</div>

		<!-- Due Date (if recurring) -->
		{#if isRecurring}
			<div class="space-y-2">
				<Label for="dueDate">Due Date</Label>
				<Input id="dueDate" type="date" bind:value={dueDate} />
			</div>
		{/if}

		<!-- Notes -->
		<div class="space-y-2">
			<Label for="notes">Notes</Label>
			<Textarea
				id="notes"
				bind:value={notes}
				placeholder="Any additional details..."
				class="font-handwriting"
				rows="3"
			/>
		</div>

		<!-- Actions -->
		<div class="flex justify-end space-x-3 pt-4">
			<Button type="button" variant="outline" onclick={closeModal}>Cancel</Button>
			<Button type="submit">Add Expense</Button>
		</div>
	</form>
</CustomModal>
