<script lang="ts">
	let { rating, max = 5 }: { rating: number; max?: number } = $props();
</script>

<span aria-label="{rating} out of {max} stars">{'★'.repeat(rating)}{'☆'.repeat(max - rating)}</span>
