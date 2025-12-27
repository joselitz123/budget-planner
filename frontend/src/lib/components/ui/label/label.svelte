<script lang="ts">
	export let className = '';
</script>

<label
	class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 {className}"
	{...$$restProps}
>
	<slot />
</label>
