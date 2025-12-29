<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import CustomModal from '$lib/components/ui/CustomModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Select } from '$lib/components/ui/select';
	import { Badge } from '$lib/components/ui/badge';
	import type { PaymentMethod } from '$lib/db/schema';
	import {
		isRequired,
		isValidAmount,
		isValidPaymentMethodName,
		isValidLastFour,
		isValidCreditLimit
	} from '$lib/utils/validation';

	export let open = false;
	export let editMethod: PaymentMethod | null = null;

	const dispatch = createEventDispatcher();

	// Payment method types
	const paymentTypes = [
		{ value: 'credit_card', label: 'Credit Card', icon: '💳' },
		{ value: 'debit_card', label: 'Debit Card', icon: '💳' },
		{ value: 'cash', label: 'Cash', icon: '💵' },
		{ value: 'ewallet', label: 'E-wallet', icon: '📱' }
	];

	// Card brands
	const cardBrands = [
		{ value: 'Visa', label: 'Visa' },
		{ value: 'Mastercard', label: 'Mastercard' },
		{ value: 'Amex', label: 'American Express' },
		{ value: 'Discover', label: 'Discover' },
		{ value: 'Chase', label: 'Chase' },
		{ value: 'Capital One', label: 'Capital One' },
		{ value: 'Citi', label: 'Citi' },
		{ value: 'Bank of America', label: 'Bank of America' }
	];

	// Form state
	let name = '';
	let type: PaymentMethod['type'] = 'credit_card';
	let lastFour = '';
	let brand = '';
	let creditLimit = '';
	let currentBalance = '';
	let isDefault = false;
	let isActive = true;

	// Form errors
	let errors = {
		name: '',
		type: '',
		lastFour: '',
		creditLimit: ''
	};

	$: isEdit = editMethod !== null;

	// Computed: Show card fields
	$: showCardFields = type === 'credit_card' || type === 'debit_card';

	// Computed: Show credit limit
	$: showCreditLimit = type === 'credit_card';

	// Computed: Show brand selector
	$: showBrandSelector = type === 'credit_card' || type === 'debit_card';

	// Computed: Brand label
	$: brandLabel = type === 'ewallet' ? 'Wallet Name' : 'Brand';

	function closeModal() {
		open = false;
		resetForm();
	}

	function resetForm() {
		name = '';
		type = 'credit_card';
		lastFour = '';
		brand = '';
		creditLimit = '';
		currentBalance = '';
		isDefault = false;
		isActive = true;
		errors = { name: '', type: '', lastFour: '', creditLimit: '' };
	}

	function populateForm(method: PaymentMethod) {
		name = method.name;
		type = method.type;
		lastFour = method.lastFour || '';
		brand = method.brand || '';
		creditLimit = method.creditLimit?.toString() || '';
		currentBalance = method.currentBalance?.toString() || '';
		isDefault = method.isDefault;
		isActive = method.isActive;
	}

	// Watch for editMethod changes
	$: if (editMethod && open) {
		populateForm(editMethod);
	}

	function validateForm(): boolean {
		let valid = true;
		errors = { name: '', type: '', lastFour: '', creditLimit: '' };

		// Validate name
		if (!isRequired(name)) {
			errors.name = 'Name is required';
			valid = false;
		} else if (!isValidPaymentMethodName(name)) {
			errors.name = 'Name must be 2-100 characters';
			valid = false;
		}

		// Validate last four (required for cards)
		if (showCardFields && !isRequired(lastFour)) {
			errors.lastFour = 'Last 4 digits are required for cards';
			valid = false;
		} else if (lastFour && !isValidLastFour(lastFour)) {
			errors.lastFour = 'Last 4 digits must be exactly 4 numbers';
			valid = false;
		}

		// Validate credit limit (for credit cards)
		if (showCreditLimit && creditLimit && !isValidCreditLimit(parseFloat(creditLimit))) {
			errors.creditLimit = 'Credit limit must be a positive number';
			valid = false;
		}

		return valid;
	}

	async function handleSubmit() {
		if (!validateForm()) {
			return;
		}

		try {
			const methodData: PaymentMethod = {
				id: editMethod?.id || crypto.randomUUID(),
				userId: '', // Will be filled by auth context
				name,
				type,
				lastFour: showCardFields ? lastFour : null,
				brand: brand || null,
				isDefault,
				isActive,
				creditLimit: showCreditLimit && creditLimit ? parseFloat(creditLimit) : null,
				currentBalance: currentBalance ? parseFloat(currentBalance) : null,
				createdAt: editMethod?.createdAt || new Date().toISOString(),
				updatedAt: new Date().toISOString()
			};

			dispatch('submit', methodData);
			closeModal();
		} catch (error) {
			console.error('[PaymentMethodForm] Error submitting form:', error);
		}
	}
</script>

<CustomModal
	bind:open
	title={isEdit ? 'Edit Payment Method' : 'Add Payment Method'}
	description={isEdit ? 'Update your payment method details' : 'Add a new payment method'}
	className="max-w-2xl"
>
	<form on:submit|preventDefault={handleSubmit} class="space-y-6">
		<!-- Type Selector -->
		<div class="space-y-2">
			<Label>Payment Method Type *</Label>
			<div class="grid grid-cols-2 gap-3">
				{#each paymentTypes as paymentType}
					<label
						class="border rounded-lg p-3 cursor-pointer flex items-center space-x-3 transition-colors {type ===
						paymentType.value
							? 'border-primary bg-primary/10'
							: 'border-line-light dark:border-line-dark hover:bg-gray-50 dark:hover:bg-gray-800'}"
					>
						<input type="radio" bind:group={type} value={paymentType.value} class="sr-only" />
						<span class="text-2xl">{paymentType.icon}</span>
						<span class="font-medium">{paymentType.label}</span>
					</label>
				{/each}
			</div>
		</div>

		<!-- Name -->
		<div class="space-y-2">
			<Label for="name">Name *</Label>
			<Input
				id="name"
				bind:value={name}
				placeholder="My Chase Visa"
				class="font-handwriting text-lg"
				maxlength="100"
			/>
			{#if errors.name}
				<p class="text-sm text-red-500 dark:text-red-400">{errors.name}</p>
			{/if}
		</div>

		<!-- Card-Specific Fields -->
		{#if showCardFields}
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label for="lastFour">Last 4 Digits *</Label>
					<Input
						id="lastFour"
						bind:value={lastFour}
						placeholder="1234"
						maxlength="4"
						class="text-center font-mono text-lg"
					/>
					{#if errors.lastFour}
						<p class="text-sm text-red-500 dark:text-red-400">{errors.lastFour}</p>
					{/if}
				</div>

				{#if showBrandSelector}
					<div class="space-y-2">
						<Label for="brand">{brandLabel}</Label>
						<Select bind:value={brand}>
							<option value="">Select {brandLabel.toLowerCase()}</option>
							{#each cardBrands as cardBrand}
								<option value={cardBrand.value}>{cardBrand.label}</option>
							{/each}
						</Select>
					</div>
				{/if}
			</div>
		{/if}

		<!-- E-wallet Wallet Name -->
		{#if type === 'ewallet'}
			<div class="space-y-2">
				<Label for="brand">Wallet Name</Label>
				<Input
					id="brand"
					bind:value={brand}
					placeholder="PayPal, Venmo, GCash, etc."
					class="font-handwriting"
				/>
			</div>
		{/if}

		<!-- Credit Limit (Credit Card Only) -->
		{#if showCreditLimit}
			<div class="space-y-2">
				<Label for="creditLimit">Credit Limit</Label>
				<Input
					id="creditLimit"
					type="number"
					step="0.01"
					min="0"
					bind:value={creditLimit}
					placeholder="5000.00"
				/>
				{#if errors.creditLimit}
					<p class="text-sm text-red-500 dark:text-red-400">{errors.creditLimit}</p>
				{/if}
			</div>
		{/if}

		<!-- Current Balance (All Types) -->
		<div class="space-y-2">
			<Label for="currentBalance">Current Balance</Label>
			<Input
				id="currentBalance"
				type="number"
				step="0.01"
				bind:value={currentBalance}
				placeholder="0.00"
			/>
			<p class="text-xs text-gray-500 dark:text-gray-400">
				Optional: Track your current balance or debt
			</p>
		</div>

		<!-- Toggles -->
		<div class="space-y-3">
			<label class="flex items-center space-x-2 cursor-pointer">
				<input type="checkbox" bind:checked={isDefault} class="w-4 h-4" />
				<span>Set as default payment method</span>
				{#if isDefault}
					<Badge variant="default">Default</Badge>
				{/if}
			</label>

			<label class="flex items-center space-x-2 cursor-pointer">
				<input type="checkbox" bind:checked={isActive} class="w-4 h-4" />
				<span>Active</span>
			</label>
		</div>

		<!-- Actions -->
		<div class="flex justify-end space-x-3 pt-4 border-t border-line-light dark:border-line-dark">
			<Button type="button" variant="outline" onclick={closeModal}>Cancel</Button>
			<Button type="submit">{isEdit ? 'Update' : 'Add'} Payment Method</Button>
		</div>
	</form>
</CustomModal>
