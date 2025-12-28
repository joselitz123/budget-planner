import { writable, get } from 'svelte/store';

/**
 * Settings state management
 */

export type Currency = 'PHP' | 'USD' | 'EUR' | 'GBP' | 'JPY';

// Currency preference
export const currency = writable<Currency>('PHP');

// Load currency from localStorage
if (typeof window !== 'undefined') {
	const savedCurrency = localStorage.getItem('currency') as Currency | null;
	if (savedCurrency && ['PHP', 'USD', 'EUR', 'GBP', 'JPY'].includes(savedCurrency)) {
		currency.set(savedCurrency);
	}
}

// Save currency to localStorage on change
currency.subscribe((c) => {
	if (typeof window !== 'undefined') {
		localStorage.setItem('currency', c);
	}
});

/**
 * Currency locale mapping
 */
export const currencyLocales: Record<Currency, string> = {
	PHP: 'en-PH',
	USD: 'en-US',
	EUR: 'de-DE', // German locale for Euro formatting
	GBP: 'en-GB',
	JPY: 'ja-JP'
};

/**
 * Currency symbols
 */
export const currencySymbols: Record<Currency, string> = {
	PHP: '₱',
	USD: '$',
	EUR: '€',
	GBP: '£',
	JPY: '¥'
};

/**
 * Get locale for currency
 */
export function getCurrencyLocale(currencyCode: Currency): string {
	return currencyLocales[currencyCode];
}

/**
 * Format amount with currency
 * This is a reactive-friendly alternative to the pure formatCurrency function
 */
export function formatCurrencyWithCode(amount: number, currencyCode: Currency): string {
	const locale = getCurrencyLocale(currencyCode);
	return new Intl.NumberFormat(locale, {
		style: 'currency',
		currency: currencyCode,
		minimumFractionDigits: 2
	}).format(amount);
}
