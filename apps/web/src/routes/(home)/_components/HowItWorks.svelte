<script lang="ts">
	const steps = [
		{
			num: '01',
			title: 'log anything',
			desc: 'watched, played, read, or heard it? drop a rating and a one-liner. five seconds, tops.'
		},
		{
			num: '02',
			title: 'build one taste graph',
			desc: 'every rating feeds a single profile across all five media. no five apps, no five logins.'
		},
		{
			num: '03',
			title: 'get recs that cross over',
			desc: "loved a moody detective game? here's the novel and the record that match its exact vibe."
		}
	];
</script>

<div id="how" class="px-12 pt-24 pb-5">
	<div class="mb-3.5 text-xs text-[#8a857c] uppercase font-['JetBrains_Mono'] tracking-[.14em]">
		how it works
	</div>
	<h2 class="m-0 mb-10 max-w-160 text-[44px] leading-[1.05] font-semibold tracking-[-.03em]">
		three steps. then a lifetime of opinions.
	</h2>
	<div class="grid grid-cols-3 gap-6">
		{#each steps as step (step.num)}
			<div class="rounded-[14px] border border-white/10 p-7.5">
				<div class="text-[15px] font-semibold text-[#8a857c] font-['JetBrains_Mono']">
					{step.num}
				</div>
				<div class="mt-4 mb-2.5 text-xl font-semibold">{step.title}</div>
				<p class="m-0 text-[15px] leading-[1.55] text-[#a5a29c]">{step.desc}</p>
			</div>
		{/each}
	</div>
</div>
