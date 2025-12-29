<script lang="ts">
	import '../app.css';
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { initDB, initBackgroundSync, stopBackgroundSync } from '$lib/db';
	import { loadBudgets, loadTransactions, currentMonth, goToPreviousMonth, goToNextMonth } from '$lib/stores';
	import { formatMonthYear } from '$lib/utils/format';
	import { initTheme, showToast } from '$lib/stores/ui';
	import { isOnline, syncIndicator } from '$lib/stores/offline';
	import { ToastContainer } from '$lib/components/ui/toast';
	import ClerkJs from '@clerk/clerk-js';

	let clerk: typeof ClerkJs | null = null;
	let showUserMenu = false;

	// Close user menu when clicking outside
	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.user-menu-container')) {
			showUserMenu = false;
		}
	}

	// Initialize app
	onMount(async () => {
		// Initialize theme
		initTheme();

		// Initialize Clerk
		try {
			const publishableKey = import.meta.env.VITE_PUBLIC_CLERK_PUBLISHABLE_KEY;
			if (publishableKey) {
				clerk = new Clerk(publishableKey);
				await clerk.load();
			}
		} catch (error) {
			console.error('[App] Error initializing Clerk:', error);
		}

		// Initialize IndexedDB
		try {
			await initDB();
			console.log('[App] IndexedDB initialized');

			// Load initial data
			await loadBudgets();
			await loadTransactions();
			console.log('[App] Initial data loaded');

			// Initialize background sync
			initBackgroundSync();
			console.log('[App] Background sync initialized');
		} catch (error) {
			console.error('[App] Error initializing:', error);
			showToast('Failed to initialize app. Some features may not work.', 'warning');
		}
	});

	async function handleLogout() {
		try {
			if (clerk) {
				await clerk.signOut();
				window.location.href = '/sign-in';
			}
		} catch (error) {
			console.error('[App] Error signing out:', error);
		}
	}

	// Cleanup on destroy
	onDestroy(() => {
		stopBackgroundSync();
		console.log('[App] Background sync stopped');
	});
</script>

<div
	class="min-h-screen bg-background-light dark:bg-background-dark font-body antialiased"
	onclick={handleClickOutside}
	role="application"
>
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
				<!-- Sync Status Indicator -->
				<div class="flex items-center space-x-1 text-xs font-medium" aria-live="polite">
					{#if $syncIndicator.status === 'syncing'}
						<span class="material-icons-outlined text-blue-500 animate-spin" style="font-size: 16px;">sync</span>
						<span class="text-blue-600 dark:text-blue-400">{$syncIndicator.label}</span>
					{:else if $syncIndicator.status === 'error'}
						<span class="material-icons-outlined text-red-500" style="font-size: 16px;">error</span>
						<span class="text-red-600 dark:text-red-400">{$syncIndicator.label}</span>
					{:else if $syncIndicator.status === 'pending'}
						<span class="material-icons-outlined text-orange-500" style="font-size: 16px;">cloud_upload</span>
						<span class="text-orange-600 dark:text-orange-400">{$syncIndicator.label}</span>
					{:else if $syncIndicator.status === 'synced'}
						<span class="material-icons-outlined text-green-500" style="font-size: 16px;">check_circle</span>
						<span class="text-green-600 dark:text-green-400 hidden sm:inline">{$syncIndicator.label}</span>
					{:else}
						<!-- offline -->
						<span class="material-icons-outlined text-gray-500" style="font-size: 16px;">cloud_off</span>
					{/if}
				</div>

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

				<!-- User Menu -->
				<div class="relative user-menu-container">
					<button
						onclick={() => (showUserMenu = !showUserMenu)}
						class="flex items-center space-x-2 px-3 py-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 transition"
						aria-label="User menu"
					>
						<span class="material-icons-outlined text-primary dark:text-white">account_circle</span>
						<span class="hidden sm:inline text-sm font-medium text-gray-700 dark:text-gray-300">
							{$page.data.user?.name || 'User'}
						</span>
					</button>

					{#if showUserMenu}
						<div
							class="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-md shadow-lg py-1 z-50"
						>
							<a
								href="/profile"
								class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
								onclick={() => (showUserMenu = false)}
							>
								<span class="material-icons-outlined text-sm align-middle mr-2">person</span>
								Profile
							</a>
							<hr class="my-1 border-gray-200 dark:border-gray-700" />
							<button
								onclick={() => {
									showUserMenu = false;
									handleLogout();
								}}
								class="w-full text-left px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-700"
							>
								<span class="material-icons-outlined text-sm align-middle mr-2">logout</span>
								Sign Out
							</button>
						</div>
					{/if}
				</div>
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
				href="/shared"
				class="flex flex-col items-center {$page.url.pathname === '/shared'
					? 'text-primary dark:text-white'
					: 'text-gray-400 dark:text-gray-500'}"
			>
				<span class="material-icons-outlined">folder_shared</span>
				<span class="text-[10px] mt-1">Shared</span>
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
