<script lang="ts">
	import { PARENT_MEDIA_TYPES, MEDIA_TYPE_CLASSES } from '@critiquefi/ui';
</script>

<div class="px-12 pt-24 pb-5">
	<div class="mb-3.5 text-xs text-[#8a857c] uppercase font-['JetBrains_Mono'] tracking-[.14em]">
		Five worlds, one shelf
	</div>
	<h2 class="m-0 mb-10 max-w-160 text-[44px] leading-[1.05] font-semibold tracking-[-.03em]">
		whatever you're into, it counts here.
	</h2>
	<div class="grid grid-cols-5 gap-3.5">
		{#each PARENT_MEDIA_TYPES as type (type.key)}
			<div
				class="flex aspect-3/4 flex-col justify-between rounded-xl p-4.5 {MEDIA_TYPE_CLASSES[
					type.key
				].card}"
			>
				<span
					class="text-[11px] font-semibold font-['JetBrains_Mono'] {MEDIA_TYPE_CLASSES[type.key]
						.label}">{type.label_plural}</span
				>
				<div>
					<div class="text-[22px] leading-none font-semibold">{type.count}</div>
					<div class="mt-1.25 text-[11px] text-[#8a857c] font-['JetBrains_Mono']">logged</div>
				</div>
			</div>
		{/each}
	</div>
</div>
