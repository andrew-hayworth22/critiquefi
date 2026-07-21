<script lang="ts">
	const marqueeItems = [
		{ stars: '★★★★☆', title: 'Dune: Part Two' },
		{ stars: '★★★★★', title: 'Disco Elysium' },
		{ stars: '★★★☆☆', title: 'Fourth Wing' },
		{ stars: '★★★★☆', title: 'Severance' },
		{ stars: '★★★★★', title: 'Brat' },
		{ stars: '★★★★☆', title: "Baldur's Gate 3" },
		{ stars: '★★★☆☆', title: 'The Bear S3' },
		{ stars: '★★★★★', title: 'Piranesi' }
	];
</script>

<div
	class="overflow-hidden border-t border-b border-white/8 py-4 mask-[linear-gradient(90deg,transparent,#000_8%,#000_92%,transparent)]"
>
	<div
		class="flex w-max animate-[cfmarquee_30s_linear_infinite] gap-9.5 font-['JetBrains_Mono'] text-[15px] font-medium text-nowrap text-[#6b6863]"
	>
		{#each [0, 1] as pass (pass)}
			{#each marqueeItems as item (item.title + pass)}
				<span>{item.stars} {item.title}</span>
			{/each}
		{/each}
	</div>
</div>

<style>
	@keyframes -global-cfmarquee {
		to {
			transform: translateX(-50%);
		}
	}
</style>
