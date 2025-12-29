<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import type { SharedBudget } from '$lib/api/shares';
	import { formatMonthYear, formatMediumDate } from '$lib/utils/format';

	export let sharedBudget: SharedBudget;

	function getPermissionVariant(permission: 'view' | 'edit') {
		return permission === 'edit' ? 'default' : 'secondary';
	}

	function getBudgetName(): string {
		if (sharedBudget.budgetName) {
			return sharedBudget.budgetName;
		}
		if (sharedBudget.budgetMonth) {
			return formatMonthYear(new Date(sharedBudget.budgetMonth + '-01'));
		}
		return 'Untitled Budget';
	}

	function getOwnerName(): string {
		if (sharedBudget.ownerName) {
			return sharedBudget.ownerName;
		}
		if (sharedBudget.ownerEmail) {
			return sharedBudget.ownerEmail;
		}
		return 'Unknown';
	}
</script>

<a
	href="/budgets/{sharedBudget.budgetId}"
	class="block bg-paper-light dark:bg-paper-dark border border-line-light dark:border-line-dark rounded-lg p-4 hover:shadow-md hover:border-primary/50 dark:hover:border-primary/50 transition-all duration-200 group"
>
	<div class="flex items-start justify-between gap-3">
		<div class="flex-1 min-w-0">
			<!-- Budget Name -->
			<div class="flex items-center gap-2 mb-2">
				<span class="material-icons-outlined text-primary dark:text-white">folder_shared</span>
				<h4 class="font-display font-bold text-primary dark:text-white text-lg truncate group-hover:text-primary/80 dark:group-hover:text-white/80 transition-colors">
					{getBudgetName()}
				</h4>
			</div>

			<!-- Owner Attribution -->
			<p class="text-sm text-gray-600 dark:text-gray-400 mb-2">
				Shared by <span class="font-medium text-gray-800 dark:text-gray-200">{getOwnerName()}</span>
			</p>

			<!-- Permission & Date -->
			<div class="flex items-center gap-3 text-sm">
				<Badge variant={getPermissionVariant(sharedBudget.permission)}>
					{sharedBudget.permission}
				</Badge>
				<span class="text-xs text-gray-500 dark:text-gray-400">
					{formatMediumDate(sharedBudget.createdAt)}
				</span>
			</div>
		</div>

		<!-- Chevron Icon -->
		<span class="material-icons-outlined text-gray-400 group-hover:text-primary dark:group-hover:text-white transition-colors shrink-0">
			chevron_right
		</span>
	</div>
</a>
