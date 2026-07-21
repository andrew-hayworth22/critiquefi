<script lang="ts">
	const stats = [
		{ value: '2.4M', label: 'members with opinions' },
		{ value: '17M', label: 'things logged' },
		{ value: '89M', label: 'ratings & hot takes' }
	];
</script>

<div
	class="mx-12 mt-16 grid grid-cols-3 gap-px overflow-hidden rounded-[14px] border border-white/8 bg-white/8"
>
	{#each stats as stat (stat.label)}
		<div class="bg-[#0e0e10] px-6 py-11 text-center">
			<div class="text-[54px] leading-none font-semibold tracking-tight">{stat.value}</div>
			<div class="mt-3 text-xs text-[#8a857c] font-['JetBrains_Mono']">{stat.label}</div>
		</div>
	{/each}
</div>
