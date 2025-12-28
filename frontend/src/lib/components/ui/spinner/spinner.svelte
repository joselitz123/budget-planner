<script lang="ts">
	import { tv, type VariantProps } from 'tailwind-variants';
	import { LoaderCircle } from 'lucide-svelte';

	const spinnerVariants = tv({
		base: 'animate-spin',
		variants: {
			size: {
				xs: 'w-4 h-4',
				sm: 'w-5 h-5',
				default: 'w-6 h-6',
				lg: 'w-8 h-8',
				xl: 'w-12 h-12'
			},
			color: {
				primary: 'text-primary dark:text-white',
				white: 'text-white',
				current: 'text-current'
			}
		},
		defaultVariants: {
			size: 'default',
			color: 'primary'
		}
	});

	type SpinnerVariants = VariantProps<typeof spinnerVariants>;

	export let size: SpinnerVariants['size'] = 'default';
	export let color: SpinnerVariants['color'] = 'primary';
	export let className = '';

	$: allClassName = spinnerVariants({ size, color }) + (className ? ` ${className}` : '');
</script>

<LoaderCircle class={allClassName} aria-hidden="true" />
