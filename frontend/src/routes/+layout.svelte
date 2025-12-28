<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { initDB } from '$lib/db';
	import { loadBudgets, loadTransactions, currentMonth, goToPreviousMonth, goToNextMonth } from '$lib/stores';
	import { formatMonthYear } from '$lib/utils/format';
	import { initTheme, showToast } from '$lib/stores/ui';
	import { isOnline } from '$lib/stores/offline';
	import { ToastContainer } from '$lib/components/ui/toast';

	// Initialize app
	onMount(async () => {
		// Initialize theme
		initTheme();

		// Initialize IndexedDB
		try {
			await initDB();
			console.log('[App] IndexedDB initialized');

			// Load initial data
			await loadBudgets();
			await loadTransactions();
			console.log('[App] Initial data loaded');
		} catch (error) {
			console.error('[App] Error initializing:', error);
			showToast('Failed to initialize app. Some features may not work.', 'warning');
		}
	});
</script>

<div class="min-h-screen bg-background-light dark:bg-background-dark font-body antialiased">
	<!-- Top Navigation -->
	<header
		class="sticky top-0 z-50 bg-paper-light/90 dark:bg-paper-dark/90 backdrop-blur-sm border-b border-line-light dark:border-line-dark px-4 py-3"
	>
		<div class="max-w-6xl mx-auto flex justify-between items-center">
			<div class="flex items-center space-x-2">
				<span class="material-icons-outlined text-primary dark:text-white text-3xl">menu_book</span>
				<h1 class="text-xl font-display font-bold text-primary dark:text-white">Budget Planner</h1>
			</div>

			<div class="flex items-center space-x-3">
				<!-- Month Selector -->
				<div class="flex items-center space-x-2">
					<button
						onclick={() => goToPreviousMonth()}
						class="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-800 transition"
						aria-label="Previous month"
					>
						<span class="material-icons-outlined text-primary dark:text-white">chevron_left</span>
					</button>
					<span class="text-sm font-medium text-gray-600 dark:text-gray-300 font-display">
						{#if $currentMonth}
							{formatMonthYear(new Date($currentMonth + '-01'))}
						{/if}
					</span>
					<button
						onclick={() => goToNextMonth()}
						class="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-800 transition"
						aria-label="Next month"
					>
						<span class="material-icons-outlined text-primary dark:text-white">chevron_right</span>
					</button>
				</div>

				<!-- Theme Toggle -->
				<button
					onclick={() => {
						// Toggle theme (handled in initTheme)
					}}
					class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 transition"
					aria-label="Toggle theme"
				>
					<span class="material-icons-outlined text-primary dark:text-white">dark_mode</span>
				</button>
			</div>
		</div>
	</header>

	<!-- Offline Status Indicator -->
	{#if !$isOnline}
		<div class="bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200 text-center py-2 px-4 text-sm font-medium">
			<span class="material-icons-outlined text-sm align-middle mr-1">wifi_off</span>
			You are offline. Changes will sync when connection is restored.
		</div>
	{/if}

	<!-- Main Content -->
	<main class="max-w-6xl mx-auto px-4 py-6">
		<slot />
	</main>

	<!-- Bottom Navigation (Mobile) -->
	<nav
		class="fixed bottom-0 left-0 right-0 bg-paper-light dark:bg-paper-dark border-t border-line-light dark:border-line-dark py-2 px-4 z-40 md:hidden"
	>
		<div class="flex justify-around items-center">
			<a
				href="/"
				class="flex flex-col items-center {$page.url.pathname === '/'
					? 'text-primary dark:text-white'
					: 'text-gray-400 dark:text-gray-500'}"
			>
				<span class="material-icons-outlined">dashboard</span>
				<span class="text-[10px] mt-1">Overview</span>
			</a>
			<a
				href="/transactions"
				class="flex flex-col items-center {$page.url.pathname === '/transactions'
					? 'text-primary dark:text-white'
					: 'text-gray-400 dark:text-gray-500'}"
			>
				<span class="material-icons-outlined">receipt_long</span>
				<span class="text-[10px] mt-1">Transactions</span>
			</a>
			<a
				href="/analytics"
				class="flex flex-col items-center {$page.url.pathname === '/analytics'
					? 'text-primary dark:text-white'
					: 'text-gray-400 dark:text-gray-500'}"
			>
				<span class="material-icons-outlined">pie_chart</span>
				<span class="text-[10px] mt-1">Analytics</span>
			</a>
			<a
				href="/bills"
				class="flex flex-col items-center {$page.url.pathname === '/bills'
					? 'text-primary dark:text-white'
					: 'text-gray-400 dark:text-gray-500'}"
			>
				<span class="material-icons-outlined">account_balance_wallet</span>
				<span class="text-[10px] mt-1">Bills</span>
			</a>
			<a
				href="/settings"
				class="flex flex-col items-center {$page.url.pathname === '/settings'
					? 'text-primary dark:text-white'
					: 'text-gray-400 dark:text-gray-500'}"
			>
				<span class="material-icons-outlined">settings</span>
				<span class="text-[10px] mt-1">Settings</span>
			</a>
		</div>
	</nav>

	<!-- Spacer for mobile navigation -->
	<div class="h-16 md:hidden"></div>

	<!-- Toast Container -->
	<ToastContainer />
</div>
