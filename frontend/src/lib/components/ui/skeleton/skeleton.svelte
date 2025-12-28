<script lang="ts">
	import { tv, type VariantProps } from 'tailwind-variants';

	const skeletonVariants = tv({
		base: 'animate-shimmer bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 dark:from-gray-700 dark:via-gray-600 dark:to-gray-700 rounded',
		variants: {
			variant: {
				card: '',
				text: 'h-4',
				circle: 'rounded-full',
				custom: ''
			}
		},
		defaultVariants: {
			variant: 'custom'
		}
	});

	type SkeletonVariants = VariantProps<typeof skeletonVariants>;

	export let variant: SkeletonVariants['variant'] = 'custom';
	export let className = '';

	$: allClassName = skeletonVariants({ variant }) + (className ? ` ${className}` : '');
</script>

<div class={allClassName} role="status" aria-label="Loading content">
	<span class="sr-only">Loading...</span>
</div>
