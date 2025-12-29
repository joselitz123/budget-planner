<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { cn } from '$lib/utils';

	interface Props extends HTMLAttributes<HTMLDivElement> {
		class?: string;
	}

	let { class: className = '', ...rest }: Props = $props();
</script>

<div class={cn('rounded-lg border bg-card text-card-foreground shadow-sm', className)} {...rest}>
	<slot />
</div>
