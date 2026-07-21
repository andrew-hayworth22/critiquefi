<script lang="ts">
	import { MEDIA_TYPE_CLASSES, PARENT_MEDIA_TYPES } from '@critiquefi/ui';

	let { section, header, description } = $props();
</script>

<div class="mb-6.5 flex items-center gap-2.75">
	<div class="flex gap-1.5">
		{#each PARENT_MEDIA_TYPES as type (type.key)}
			<span class="h-3 w-3 rounded-[3px] {MEDIA_TYPE_CLASSES[type.key].dot}"></span>
		{/each}
	</div>
	{@render section()}
</div>
<h1 class="m-0 text-[74px] leading-[1.01] font-semibold text-balance tracking-[-.035em]">
	{@render header()}
</h1>
<p class="mt-7 mb-0 max-w-135 text-[19px] leading-[1.55] text-[#a5a29c]">
	{@render description()}
</p>
