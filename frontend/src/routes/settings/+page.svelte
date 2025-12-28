<script lang="ts">
	import { get } from 'svelte/store';
	import { theme, toggleTheme, showToast } from '$lib/stores/ui';
	import { currency, formatCurrencyWithCode } from '$lib/stores/settings';
	import { budgetStore, transactionStore, categoryStore } from '$lib/db/stores';
	import { loadBudgets, loadTransactions } from '$lib/stores';

	/**
	 * Export all data as JSON
	 */
	async function exportData() {
		try {
			// Get all data from IndexedDB
			const budgets = await budgetStore.getAll();
			const transactions = await transactionStore.getAll();
			const categories = await categoryStore.getAll();

			const data = {
				version: 1,
				exportedAt: new Date().toISOString(),
				currency: get(currency),
				budgets,
				transactions,
				categories
			};

			// Create downloadable JSON file
			const blob = new Blob([JSON.stringify(data, null, 2)], {
				type: 'application/json'
			});
			const url = URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = `budget-planner-backup-${new Date().toISOString().split('T')[0]}.json`;
			a.click();
			URL.revokeObjectURL(url);

			showToast('Data exported successfully!', 'success');
		} catch (error) {
			console.error('Export failed:', error);
			showToast('Failed to export data', 'error');
		}
	}

	/**
	 * Import data from JSON file
	 */
	async function importData(event: Event) {
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];

		if (!file) return;

		try {
			const text = await file.text();
			const data = JSON.parse(text);

			// Validate structure
			if (!data.budgets || !Array.isArray(data.budgets)) {
				throw new Error('Invalid backup file: missing budgets');
			}
			if (!data.transactions || !Array.isArray(data.transactions)) {
				throw new Error('Invalid backup file: missing transactions');
			}

			// Confirm import
			if (
				!confirm(
					`This will replace all your current data with ${data.transactions.length} transactions and ${data.budgets.length} budgets. Continue?`
				)
			) {
				input.value = '';
				return;
			}

			// Clear existing data
			const allBudgets = await budgetStore.getAll();
			const allTransactions = await transactionStore.getAll();
			const allCategories = await categoryStore.getAll();

			for (const budget of allBudgets) {
				await budgetStore.delete(budget.id);
			}
			for (const transaction of allTransactions) {
				await transactionStore.delete(transaction.id);
			}
			if (data.categories) {
				for (const category of allCategories) {
					await categoryStore.delete(category.id);
				}
			}

			// Import data
			for (const budget of data.budgets) {
				await budgetStore.create(budget);
			}
			for (const transaction of data.transactions) {
				await transactionStore.create(transaction);
			}
			if (data.categories) {
				for (const category of data.categories) {
					await categoryStore.create(category);
				}
			}

			// Reload stores
			await loadBudgets();
			await loadTransactions();

			// Import currency if present
			if (data.currency && ['PHP', 'USD', 'EUR', 'GBP', 'JPY'].includes(data.currency)) {
				currency.set(data.currency);
			}

			showToast('Data imported successfully!', 'success');
		} catch (error) {
			console.error('Import failed:', error);
			showToast(error instanceof Error ? error.message : 'Failed to import data', 'error');
		}

		// Reset input
		input.value = '';
	}
</script>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="mb-6">
		<h2 class="text-3xl font-display font-bold text-primary dark:text-white mb-2">
			Settings
		</h2>
		<p class="text-sm font-handwriting text-gray-500 dark:text-gray-400 text-xl">
			Customize your budgeting experience
		</p>
	</div>

	<!-- Appearance Section -->
	<div
		class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark overflow-hidden mb-6"
	>
		<div class="bg-primary dark:bg-gray-700 text-white p-3 flex justify-between items-center">
			<h3 class="font-display font-semibold tracking-wide uppercase text-sm">Appearance</h3>
			<span class="material-icons text-sm opacity-80">palette</span>
		</div>
		<div class="p-4">
			<div class="flex items-center justify-between">
				<div>
					<p class="font-semibold text-gray-800 dark:text-gray-100">Dark Mode</p>
					<p class="text-sm text-gray-500 dark:text-gray-400">
						{$theme === 'dark' ? 'Currently using dark theme' : 'Currently using light theme'}
					</p>
				</div>
				<button
					onclick={toggleTheme}
					class="px-4 py-2 rounded-lg border border-line-light dark:border-line-dark hover:bg-gray-100 dark:hover:bg-gray-800 transition flex items-center space-x-2"
					aria-label="Toggle theme"
				>
					<span class="text-2xl">{$theme === 'dark' ? '🌙' : '☀️'}</span>
					<span class="text-sm font-medium text-gray-700 dark:text-gray-300">
						{$theme === 'dark' ? 'Dark' : 'Light'}
					</span>
				</button>
			</div>
		</div>
	</div>

	<!-- Preferences Section -->
	<div
		class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark overflow-hidden mb-6"
	>
		<div class="bg-primary dark:bg-gray-700 text-white p-3 flex justify-between items-center">
			<h3 class="font-display font-semibold tracking-wide uppercase text-sm">Preferences</h3>
			<span class="material-icons text-sm opacity-80">tune</span>
		</div>
		<div class="p-4">
			<div class="flex items-center justify-between">
				<div>
					<p class="font-semibold text-gray-800 dark:text-gray-100">Currency</p>
					<p class="text-sm text-gray-500 dark:text-gray-400">
						Select your preferred currency
					</p>
				</div>
				<select
					bind:value={$currency}
					class="px-4 py-2 rounded-lg border border-line-light dark:border-line-dark bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-100 focus:ring-2 focus:ring-primary focus:outline-none"
				>
					<option value="PHP">🇵🇭 Philippine Peso (₱)</option>
					<option value="USD">🇺🇸 US Dollar ($)</option>
					<option value="EUR">🇪🇺 Euro (€)</option>
					<option value="GBP">🇬🇧 British Pound (£)</option>
					<option value="JPY">🇯🇵 Japanese Yen (¥)</option>
				</select>
			</div>
			<p class="mt-3 text-xs text-gray-400 dark:text-gray-500">
				Example: {formatCurrencyWithCode(1234.56, $currency)}
			</p>
		</div>
	</div>

	<!-- Data Management Section -->
	<div
		class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark overflow-hidden"
	>
		<div class="bg-primary dark:bg-gray-700 text-white p-3 flex justify-between items-center">
			<h3 class="font-display font-semibold tracking-wide uppercase text-sm">Data Management</h3>
			<span class="material-icons text-sm opacity-80">storage</span>
		</div>
		<div class="p-4">
			<p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
				Export your data as a backup or import from a previous backup file.
			</p>
			<div class="flex flex-col sm:flex-row gap-3">
				<button
					onclick={exportData}
					class="flex-1 flex items-center justify-center space-x-2 px-4 py-3 rounded-lg border border-line-light dark:border-line-dark hover:bg-gray-100 dark:hover:bg-gray-800 transition"
				>
					<span class="material-icons text-primary dark:text-white">download</span>
					<span class="font-medium text-gray-700 dark:text-gray-300">Export Data</span>
				</button>
				<label class="flex-1 flex items-center justify-center space-x-2 px-4 py-3 rounded-lg border border-line-light dark:border-line-dark hover:bg-gray-100 dark:hover:bg-gray-800 transition cursor-pointer">
					<input type="file" accept=".json" onchange={importData} hidden />
					<span class="material-icons text-primary dark:text-white">upload</span>
					<span class="font-medium text-gray-700 dark:text-gray-300">Import Data</span>
				</label>
			</div>
			<p class="mt-4 text-xs text-gray-400 dark:text-gray-500">
				⚠️ Importing will replace all existing data. Make sure to export a backup first!
			</p>
		</div>
	</div>
</div>
