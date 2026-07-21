<script lang="ts">
	interface Highlight {
		hue: number;
		title: string;
		description: string;
	}
	interface Props {
		highlights: Highlight[];
	}

	const { highlights }: Props = $props();
</script>

<div class="rounded-[16px] border border-white/10 bg-white/2 p-7 px-7.5">
	<div class="mb-5 text-[11px] text-[#8a857c] uppercase font-['JetBrains_Mono'] tracking-[.12em]">
		the short version
	</div>
	<div class="grid grid-cols-3 gap-6">
		{#each highlights as item (item.title)}
			<div>
				<div
					class="mb-3.5 h-8.5 w-8.5 rounded-[9px]"
					style="background:oklch(.7 .16 {item.hue} / .18);border:1px solid oklch(.7 .16 {item.hue} / .5)"
				></div>
				<div class="mb-1.75 text-[15px] font-semibold">{item.title}</div>
				<p class="m-0 text-[13.5px] leading-[1.5] text-[#a5a29c]">{item.description}</p>
			</div>
		{/each}
	</div>
</div>
