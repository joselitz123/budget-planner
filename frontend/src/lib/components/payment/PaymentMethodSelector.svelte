<script lang="ts">
	import { Select } from '$lib/components/ui/select';
	import type { PaymentMethod } from '$lib/db/schema';

	/**
	 * Currently selected payment method ID
	 */
	export let selectedId: string | null = null;

	/**
	 * List of available payment methods
	 */
	export let paymentMethods: PaymentMethod[] = [];

	/**
	 * Callback when selection changes
	 */
	export let onChange: (id: string | null) => void;

	/**
	 * Format payment method for display in dropdown
	 */
	function formatPaymentMethod(method: PaymentMethod): string {
		if (method.type === 'credit_card' || method.type === 'debit_card') {
			const brand = method.brand || 'Card';
			const lastFour = method.lastFour || '----';
			const suffix = method.isDefault ? ' (Default)' : '';
			return `${brand} •••• ${lastFour}${suffix}`;
		} else if (method.type === 'ewallet') {
			const name = method.brand || 'E-wallet';
			const suffix = method.isDefault ? ' (Default)' : '';
			return `${name}${suffix}`;
		} else {
			// Cash
			const suffix = method.isDefault ? ' (Default)' : '';
			return `Cash${suffix}`;
		}
	}

	/**
	 * Get icon for payment method type
	 */
	function getTypeIcon(method: PaymentMethod): string {
		switch (method.type) {
			case 'credit_card':
			case 'debit_card':
				return '💳';
			case 'cash':
				return '💵';
			case 'ewallet':
				return '📱';
			default:
				return '💳';
		}
	}

	function handleChange(event: Event) {
		const target = event.target as HTMLSelectElement;
		onChange(target.value || null);
	}
</script>

<select
	bind:value={selectedId}
	onchange={handleChange}
	class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
>
	<option value="">None</option>
	{#each paymentMethods as method}
		<option value={method.id}>
			{getTypeIcon(method)} {formatPaymentMethod(method)}
		</option>
	{/each}
</select>
