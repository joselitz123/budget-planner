<script lang="ts">
	import { onMount } from 'svelte';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import PaymentMethodList from '$lib/components/payment/PaymentMethodList.svelte';
	import PaymentMethodForm from '$lib/components/payment/PaymentMethodForm.svelte';
	import {
		paymentMethods,
		paymentMethodsLoading,
		loadPaymentMethods,
		addPaymentMethod,
		updatePaymentMethod,
		deletePaymentMethod,
		setDefaultPaymentMethod
	} from '$lib/stores/paymentMethods';
	import type { PaymentMethod } from '$lib/db/schema';
	import { showToast } from '$lib/stores/ui';

	let showForm = false;
	let editingMethod: PaymentMethod | null = null;

	onMount(async () => {
		await loadPaymentMethods();
	});

	function handleOpenForm() {
		editingMethod = null;
		showForm = true;
	}

	function handleEdit(method: PaymentMethod) {
		editingMethod = method;
		showForm = true;
	}

	function handleSetDefault(methodId: string) {
		setDefaultPaymentMethod(methodId);
	}

	function handleDelete(methodId: string) {
		if (confirm('Are you sure you want to delete this payment method?')) {
			deletePaymentMethod(methodId);
		}
	}

	async function handleSubmit(event: CustomEvent<PaymentMethod>) {
		const method = event.detail;

		try {
			if (editingMethod) {
				// Update existing
				await updatePaymentMethod(method);
			} else {
				// Add new
				await addPaymentMethod(method);
			}
			showForm = false;
			editingMethod = null;
		} catch (error) {
			console.error('[PaymentMethodsPage] Error saving payment method:', error);
		}
	}
</script>

<svelte:head>
	<title>Payment Methods - Budget Planner</title>
</svelte:head>

<div class="max-w-4xl mx-auto p-4 md:p-6">
	<!-- Page Header -->
	<div class="mb-6">
		<h1 class="text-3xl font-display font-bold text-primary dark:text-white mb-2">
			Payment Methods
		</h1>
		<p class="text-gray-600 dark:text-gray-400">Manage your payment methods for better expense tracking</p>
	</div>

	<!-- Summary Cards -->
	{#if $paymentMethodsLoading}
		<div class="grid grid-cols-2 gap-4 mb-6">
			<div
				class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-4"
			>
				<Skeleton variant="text" className="h-4 w-20 mb-2" />
				<Skeleton variant="card" className="h-8 w-16" />
			</div>
			<div
				class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-4"
			>
				<Skeleton variant="text" className="h-4 w-20 mb-2" />
				<Skeleton variant="card" className="h-8 w-16" />
			</div>
		</div>
	{:else}
		<div class="grid grid-cols-2 gap-4 mb-6">
			<!-- Total Methods -->
			<div
				class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-4 relative overflow-hidden"
			>
				<div class="absolute top-0 right-0 w-20 h-20 bg-gradient-to-br from-primary/5 to-transparent rounded-bl-full"></div>
				<p class="text-xs text-gray-500 dark:text-gray-400 mb-1 relative z-10">Total Methods</p>
				<p class="text-2xl font-bold font-display text-primary dark:text-white relative z-10">
					{$paymentMethods.length}
				</p>
			</div>

			<!-- Default Method -->
			<div
				class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark p-4 relative overflow-hidden"
			>
				<div class="absolute top-0 right-0 w-20 h-20 bg-gradient-to-br from-primary/5 to-transparent rounded-bl-full"></div>
				<p class="text-xs text-gray-500 dark:text-gray-400 mb-1 relative z-10">Default Method</p>
				{#if $paymentMethods.find((m) => m.isDefault)}
					<p
						class="text-sm font-semibold text-primary dark:text-white relative z-10 truncate"
					>
						{$paymentMethods.find((m) => m.isDefault)?.name}
					</p>
				{:else}
					<p class="text-sm text-gray-400 dark:text-gray-500 relative z-10">None set</p>
				{/if}
			</div>
		</div>
	{/if}

	<!-- Payment Methods List -->
	<div class="relative">
		<!-- Spiral binding decoration -->
		<div
			class="absolute left-0 top-0 bottom-0 w-6 bg-gray-100 dark:bg-gray-800 border-r border-line-light dark:border-line-dark flex flex-col items-center py-4 space-y-4 z-10 rounded-l-xl hidden md:flex"
		>
			<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
			<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
			<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
		</div>

		<!-- List -->
		<div class="pl-0 md:pl-8">
			<PaymentMethodList
				paymentMethods={$paymentMethods}
				loading={$paymentMethodsLoading}
				on:edit={(e) => handleEdit(e.detail)}
				on:setDefault={(e) => handleSetDefault(e.detail)}
				on:delete={(e) => handleDelete(e.detail)}
			/>
		</div>
	</div>

	<!-- Floating Action Button -->
	<button
		onclick={handleOpenForm}
		class="fixed bottom-24 md:bottom-8 right-6 bg-primary text-white p-4 rounded-full shadow-lg hover:shadow-xl hover:bg-gray-700 transition-all transform hover:-translate-y-1 z-30"
		aria-label="Add payment method"
		type="button"
	>
		<span class="material-icons-outlined text-2xl">add</span>
	</button>

	<!-- Form Modal -->
	<PaymentMethodForm
		bind:open={showForm}
		bind:editMethod={editingMethod}
		on:submit={handleSubmit}
	/>
</div>

<style>
	.material-icons-outlined {
		font-family: 'Material Icons Outlined';
		font-weight: normal;
		font-style: normal;
		line-height: 1;
		letter-spacing: normal;
		text-transform: none;
		display: inline-block;
		white-space: nowrap;
		word-wrap: normal;
		direction: ltr;
	}
</style>
