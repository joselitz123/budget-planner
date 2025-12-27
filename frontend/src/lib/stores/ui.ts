import { writable, derived, get } from 'svelte/store';

/**
 * UI state management
 */

// Theme (light/dark mode)
export const theme = writable<'light' | 'dark'>('light');

// Mobile navigation open/closed
export const mobileNavOpen = writable<boolean>(false);

// Current route (derived from SvelteKit)
export const currentRoute = writable<string>('/');

// Loading state
export const isLoading = writable<boolean>(false);

// Toast notifications
export interface Toast {
	id: string;
	message: string;
	type: 'success' | 'error' | 'info' | 'warning';
	duration?: number;
}

export const toasts = writable<Toast[]>([]);

/**
 * Initialize theme from localStorage or system preference
 */
export function initTheme(): void {
	if (typeof window === 'undefined') return;

	const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null;
	if (savedTheme) {
		theme.set(savedTheme);
	} else {
		const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		theme.set(prefersDark ? 'dark' : 'light');
	}

	// Apply theme to document
	theme.subscribe((t) => {
		if (typeof document !== 'undefined') {
			document.documentElement.classList.toggle('dark', t === 'dark');
			localStorage.setItem('theme', t);
		}
	});
}

/**
 * Toggle theme
 */
export function toggleTheme(): void {
	theme.update((t) => (t === 'light' ? 'dark' : 'light'));
}

/**
 * Show toast notification
 */
export function showToast(
	message: string,
	type: 'success' | 'error' | 'info' | 'warning' = 'info',
	duration = 3000
): void {
	const id = crypto.randomUUID();
	const toast: Toast = { id, message, type, duration };

	toasts.update((current) => [...current, toast]);

	// Auto-remove after duration
	if (duration > 0) {
		setTimeout(() => {
			removeToast(id);
		}, duration);
	}
}

/**
 * Remove toast notification
 */
export function removeToast(id: string): void {
	toasts.update((current) => current.filter((t) => t.id !== id));
}

/**
 * Clear all toasts
 */
export function clearToasts(): void {
	toasts.set([]);
}
