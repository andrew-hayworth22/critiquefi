<script lang="ts">
	const features = [
		{
			hue: 25,
			title: 'one profile, five worlds',
			desc: 'films next to games next to albums. your taste, finally uncompartmentalised.'
		},
		{
			hue: 155,
			title: 'recs that cross over',
			desc: 'the engine reads across everything you love, so it never runs out of good ideas.'
		},
		{
			hue: 245,
			title: 'follow the right people',
			desc: 'friends with great taste in one thing usually have it in others. find them here.'
		},
		{
			hue: 75,
			title: 'lists worth arguing over',
			desc: 'build ranked lists that mix every medium. "best comfort media, period."'
		},
		{
			hue: 320,
			title: 'your year, reviewed',
			desc: 'a wrapped-style recap that spans films, games, books, shows and albums at once.'
		},
		{
			hue: 130,
			title: 'import in one click',
			desc: "bring your history from the apps you're leaving behind. we'll take it from here."
		}
	];
</script>

<div id="features" class="px-12 pt-24 pb-5">
	<div class="grid grid-cols-3 gap-6">
		{#each features as feature (feature.title)}
			<div class="rounded-[14px] border border-white/10 p-7.5">
				<div
					class="mb-4.5 h-9.5 w-9.5 rounded-[9px]"
					style="background:oklch(.7 .16 {feature.hue} / .18);border:1px solid oklch(.7 .16 {feature.hue} / .5)"
				></div>
				<div class="mb-2.25 text-lg font-semibold">{feature.title}</div>
				<p class="m-0 text-[14.5px] leading-[1.55] text-[#a5a29c]">{feature.desc}</p>
			</div>
		{/each}
	</div>
</div>
