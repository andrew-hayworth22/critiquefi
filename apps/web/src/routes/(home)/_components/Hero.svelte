<script lang="ts">
	import PageHeader from './PageHeader.svelte';
</script>

<div class="px-12 pb-17">
	<PageHeader>
		{#snippet section()}
			<span class="text-xs text-[#8a857c] uppercase font-['JetBrains_Mono'] tracking-widest">
				Film · TV · Games · Books · Music
			</span>
		{/snippet}

		{#snippet header()}
			every film, show, game, book and album you've ever loved.
			<span class="text-[#8a857c]">one profile.</span>
		{/snippet}

		{#snippet description()}
			log it, rate it, argue about it. critiquefi is the first place your whole taste actually lives
			together - and the only one smart enough to recommend across all of it.
		{/snippet}
	</PageHeader>
	<div class="mt-9 flex flex-wrap items-center gap-3">
		<a
			href="#"
			class="rounded-[9px] bg-[#f2f0ec] px-7.5 py-4 text-[15px] font-semibold text-[#0e0e10]"
			>Sign up free →</a
		>
		<a
			href="#how"
			class="rounded-[9px] border border-white/16 px-6.5 py-4 text-[15px] font-medium text-[#f2f0ec]"
			>Take the tour</a
		>
		<span class="ml-1.5 text-xs text-[#6b6863] font-['JetBrains_Mono']">no card · free forever</span
		>
	</div>
</div>
