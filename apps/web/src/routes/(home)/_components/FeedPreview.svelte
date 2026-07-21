<script lang="ts">
	const feedPreview = [
		{
			hue: 25,
			user: '@mira',
			verb: 'rated',
			title: 'Disco Elysium',
			rating: '★★★★★',
			quote: 'Ruined every other game for me. Worth it.'
		},
		{
			hue: 245,
			user: '@dev',
			verb: 'logged',
			title: 'Severance S2',
			rating: '★★★★☆',
			quote: 'The elevator does more character work than most films.'
		},
		{
			hue: 320,
			user: '@sol',
			verb: 'rated',
			title: 'Brat',
			rating: '★★★★★',
			quote: "Album of the summer, no notes, don't @ me."
		}
	];
</script>

<div id="feed" class="px-12 pt-24 pb-5">
	<div class="grid grid-cols-2 items-center gap-12">
		<div>
			<div class="mb-3.5 text-xs text-[#8a857c] uppercase font-['JetBrains_Mono'] tracking-[.14em]">
				The feed
			</div>
			<h2 class="m-0 mb-4.5 max-w-110 text-[44px] leading-[1.05] font-semibold tracking-[-.03em]">
				a timeline of everyone's freshest takes.
			</h2>
			<p class="m-0 max-w-105 text-base leading-[1.6] text-[#a5a29c]">
				see what your people are rating in real time — a five-star album at breakfast, a two-star
				blockbuster by midnight. reply, quote, and pile on.
			</p>
		</div>
		<div class="flex flex-col gap-3">
			{#each feedPreview as entry (entry.user)}
				<div class="flex items-center gap-3.5 rounded-xl border border-white/10 px-5 py-4.5">
					<span
						class="h-10 w-10 flex-none rounded-full"
						style="background:oklch(.7 .16 {entry.hue} / .3)"
					></span>
					<div>
						<div class="text-sm">
							<b class="font-semibold">{entry.user}</b>
							<span class="text-[#8a857c]">{entry.verb}</span>
							<b class="font-semibold">{entry.title}</b>
							<span style="color:oklch(.8 .13 75)">{entry.rating}</span>
						</div>
						<div class="mt-0.75 text-[13px] text-[#a5a29c]">"{entry.quote}"</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>
