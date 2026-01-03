/**
 * Supported currency codes with their symbols
 */
export type CurrencyCode = 'PHP' | 'USD' | 'EUR' | 'GBP' | 'JPY' | 'CNY' | 'CAD' | 'AUD';

/**
 * Currency symbol mapping for supported currencies
 */
const CURRENCY_SYMBOLS: Record<CurrencyCode, string> = {
	PHP: '₱',
	USD: '$',
	EUR: '€',
	GBP: '£',
	JPY: '¥',
	CNY: '¥',
	CAD: 'C$',
	AUD: 'A$'
};

/**
 * Get currency symbol for a given currency code
 * @param currencyCode - ISO 4217 currency code (e.g., 'USD', 'EUR')
 * @returns Currency symbol or empty string if currency is not supported
 *
 * @example
 * getCurrencySymbol('USD') // returns '$'
 * getCurrencySymbol('EUR') // returns '€'
 * getCurrencySymbol('INVALID') // returns ''
 */
export function getCurrencySymbol(currencyCode: string): string {
	// Validate input
	if (!currencyCode || typeof currencyCode !== 'string') {
		console.warn('getCurrencySymbol: Invalid currency code provided');
		return '';
	}
	
	// Normalize to uppercase for case-insensitive matching
	const normalizedCode = currencyCode.toUpperCase() as CurrencyCode;
	
	// Return symbol if currency is supported
	return CURRENCY_SYMBOLS[normalizedCode] || '';
}

/**
 * Format a number as currency using the specified currency code
 * @param amount - The amount to format
 * @param currencyCode - ISO 4217 currency code (defaults to 'PHP')
 * @returns Formatted currency string
 *
 * @example
 * formatCurrencyWithCode(1000, 'USD') // returns '$1,000.00'
 * formatCurrencyWithCode(1000, 'EUR') // returns '€1,000.00'
 */
export function formatCurrencyWithCode(amount: number, currencyCode: CurrencyCode = 'PHP'): string {
	// Validate amount
	if (typeof amount !== 'number' || isNaN(amount)) {
		console.warn('formatCurrencyWithCode: Invalid amount provided');
		return '';
	}
	
	return new Intl.NumberFormat('en-US', {
		style: 'currency',
		currency: currencyCode,
		minimumFractionDigits: 2,
		maximumFractionDigits: 2
	}).format(amount);
}

/**
 * Format a number as currency (Philippine Peso)
 * @deprecated Use formatCurrencyWithCode instead for multi-currency support
 */
export function formatCurrency(amount: number): string {
	return formatCurrencyWithCode(amount, 'PHP');
}

/**
 * Format a date as a readable string (e.g., "January 2024")
 */
export function formatMonthYear(date: Date): string {
	return new Intl.DateTimeFormat('en-US', {
		month: 'long',
		year: 'numeric'
	}).format(date);
}

/**
 * Format a date as a short string (e.g., "01/15/2024")
 */
export function formatShortDate(date: Date | string): string {
	const d = typeof date === 'string' ? new Date(date) : date;
	return new Intl.DateTimeFormat('en-US', {
		month: '2-digit',
		day: '2-digit',
		year: 'numeric'
	}).format(d);
}

/**
 * Format a date as a medium string (e.g., "Jan 15, 2024")
 */
export function formatMediumDate(date: Date | string): string {
	const d = typeof date === 'string' ? new Date(date) : date;
	return new Intl.DateTimeFormat('en-US', {
		month: 'short',
		day: 'numeric',
		year: 'numeric'
	}).format(d);
}

/**
 * Get category color for display
 */
export function getCategoryColor(categoryName: string): string {
	const colors: Record<string, string> = {
		'Housing': 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200',
		'Food': 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200',
		'Transportation': 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
		'Health Care': 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
		'Personal': 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
		'Entertainment': 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200',
		'Bills': 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
		'Social': 'bg-pink-100 text-pink-800 dark:bg-pink-900 dark:text-pink-200',
		'Health': 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
	};

	return colors[categoryName] || 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200';
}

/**
 * Calculate percentage with precision
 */
export function calculatePercentage(value: number, total: number): number {
	if (total === 0) return 0;
	return Math.round((value / total) * 100 * 100) / 100; // Round to 2 decimal places
}

/**
 * Get month key for storage (e.g., "2024-01")
 */
export function getMonthKey(date: Date): string {
	return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;
}

/**
 * Parse month key to Date
 */
export function parseMonthKey(monthKey: string): Date {
	const [year, month] = monthKey.split('-').map(Number);
	return new Date(year, month - 1, 1);
}
