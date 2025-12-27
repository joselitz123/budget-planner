/**
 * Email validation regex
 */
export const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

/**
 * Validate email address
 */
export function isValidEmail(email: string): boolean {
	return EMAIL_REGEX.test(email);
}

/**
 * Validate amount (positive number)
 */
export function isValidAmount(amount: number | string): boolean {
	if (typeof amount === 'string') {
		return !isNaN(parseFloat(amount)) && parseFloat(amount) > 0;
	}
	return !isNaN(amount) && amount > 0;
}

/**
 * Validate required field
 */
export function isRequired(value: string | number | null | undefined): boolean {
	if (value === null || value === undefined) return false;
	if (typeof value === 'string') return value.trim().length > 0;
	return true;
}

/**
 * Validate date range
 */
export function isValidDateRange(startDate: Date, endDate: Date): boolean {
	return startDate <= endDate;
}

/**
 * Validate transaction amount against budget limit
 */
export function isWithinBudget(spent: number, limit: number): boolean {
	return spent <= limit;
}
