<script lang="ts">
	import Nav from '../_components/Nav.svelte';
	import PageHeader from '../../_components/PageHeader.svelte';
	import ShortVersion from '../_components/ShortVersion.svelte';

	const highlights = [
		{
			hue: 155,
			title: 'we never sell your data',
			description:
				'not to advertisers, not to data brokers. your ratings power your recs - nothing else.'
		},
		{
			hue: 245,
			title: "you control what's public",
			description: 'set your whole profile - or any single log - to public, followers, or private.'
		},
		{
			hue: 320,
			title: 'export or delete anytime',
			description:
				"download everything you've logged, or wipe your account for good, from settings."
		}
	];

	const navItems = [
		{ id: 'collect', label: 'what we collect' },
		{ id: 'use', label: 'how we use it' },
		{ id: 'recs', label: 'recommendations' },
		{ id: 'sharing', label: 'sharing & visibility' },
		{ id: 'rights', label: 'your rights' },
		{ id: 'retention', label: 'retention' },
		{ id: 'contact', label: 'contact' }
	];

	const collectItems = [
		{
			title: 'account info',
			desc: 'your email, username, and password (hashed). if you sign in with google or apple, we receive a basic identifier - never your password there.'
		},
		{
			title: 'your activity',
			desc: 'the films, games, books, shows, and albums you log - plus ratings, reviews, lists, likes, follows, and the rating scale you prefer.'
		},
		{
			title: 'technical data',
			desc: "standard log data - device type, browser, and coarse region - used to keep the service secure and working. we don't build advertising profiles."
		}
	];

	const rightsChips = ['access your data', 'export everything', 'correct info', 'delete account'];
</script>

<div class="mx-12">
	<PageHeader>
		{#snippet section()}
			<span class="text-xs text-[#8a857c] uppercase font-['JetBrains_Mono'] tracking-widest">
				Legal
			</span>
		{/snippet}

		{#snippet header()}
			privacy policy
		{/snippet}

		{#snippet description()}
			the taste is yours. this is the plain-english version of what we collect, why, and the
			controls you have - the full legal text follows below.
		{/snippet}
	</PageHeader>

	<div class="mt-5.5 text-xs text-[#6b6863] font-['JetBrains_Mono']">
		last updated: july 21, 2026 · effective: august 1, 2026
	</div>

	<div class="mt-12">
		<ShortVersion {highlights} />
		<div class="mt-11 grid grid-cols-[220px_minmax(0,1fr)] items-start gap-12 pb-5">
			<Nav {navItems} />

			<div class="flex max-w-170 flex-col gap-11">
				<div id="collect" class="scroll-mt-23">
					<div class="mb-3 text-xs font-semibold text-[#8a857c] font-['JetBrains_Mono']">01</div>
					<h2 class="m-0 mb-4 text-[26px] leading-[1.1] font-semibold tracking-[-.02em]">
						what we collect
					</h2>
					<p class="m-0 mb-4 text-[15.5px] leading-[1.65] text-[#c9c6c0] text-pretty">
						to run critiquefi we collect a small set of information, in three buckets:
					</p>
					<div class="flex flex-col gap-3">
						{#each collectItems as item (item.title)}
							<div class="rounded-[11px] border border-white/8 px-4.5 py-4">
								<div class="mb-1.25 text-[14.5px] font-semibold">{item.title}</div>
								<div class="text-sm leading-[1.55] text-[#a5a29c]">{item.desc}</div>
							</div>
						{/each}
					</div>
				</div>

				<div id="use" class="scroll-mt-23">
					<div class="mb-3 text-xs font-semibold text-[#8a857c] font-['JetBrains_Mono']">02</div>
					<h2 class="m-0 mb-4 text-[26px] leading-[1.1] font-semibold tracking-[-.02em]">
						how we use it
					</h2>
					<p class="m-0 text-[15.5px] leading-[1.65] text-[#c9c6c0] text-pretty">
						we use your information only to provide and improve critiquefi: to run your profile and
						feed, generate cross-media recommendations, keep the platform safe, and send the
						notifications you've opted into. we do not sell personal data, and we don't share it
						with advertisers.
					</p>
				</div>

				<div id="recs" class="scroll-mt-23">
					<div class="mb-3 text-xs font-semibold text-[#8a857c] font-['JetBrains_Mono']">03</div>
					<h2 class="m-0 mb-4 text-[26px] leading-[1.1] font-semibold tracking-[-.02em]">
						recommendations
					</h2>
					<p class="m-0 text-[15.5px] leading-[1.65] text-[#c9c6c0] text-pretty">
						our recommendation engine reads across everything you rate - a moody detective game can
						surface a novel or an album that matches its mood. this processing stays within
						critiquefi and is used to personalize your experience. you can opt certain items out of
						your taste graph, or exclude mature-rated titles, from settings.
					</p>
				</div>

				<div id="sharing" class="scroll-mt-23">
					<div class="mb-3 text-xs font-semibold text-[#8a857c] font-['JetBrains_Mono']">04</div>
					<h2 class="m-0 mb-4 text-[26px] leading-[1.1] font-semibold tracking-[-.02em]">
						sharing &amp; visibility
					</h2>
					<p class="m-0 text-[15.5px] leading-[1.65] text-[#c9c6c0] text-pretty">
						what appears on your public profile and in the feed is controlled entirely by you.
						profile visibility can be set to public, followers-only, or private, and individual logs
						can be hidden. we share data with third parties only in narrow cases: service providers
						who host and secure the platform under contract, or when required by law.
					</p>
				</div>

				<div id="rights" class="scroll-mt-23">
					<div class="mb-3 text-xs font-semibold text-[#8a857c] font-['JetBrains_Mono']">05</div>
					<h2 class="m-0 mb-4 text-[26px] leading-[1.1] font-semibold tracking-[-.02em]">
						your rights &amp; controls
					</h2>
					<p class="m-0 mb-4 text-[15.5px] leading-[1.65] text-[#c9c6c0] text-pretty">
						you can access, correct, export, or delete your data at any time - most of it directly
						from settings. depending on where you live, you may have additional rights under laws
						like gdpr or ccpa, which we honor for all users.
					</p>
					<div class="flex flex-wrap gap-2">
						{#each rightsChips as chip (chip)}
							<span
								class="rounded-full border border-white/14 px-3.25 py-2 text-[12.5px] text-[#c9c6c0] font-['JetBrains_Mono']"
							>
								{chip}
							</span>
						{/each}
					</div>
				</div>

				<div id="retention" class="scroll-mt-23">
					<div class="mb-3 text-xs font-semibold text-[#8a857c] font-['JetBrains_Mono']">06</div>
					<h2 class="m-0 mb-4 text-[26px] leading-[1.1] font-semibold tracking-[-.02em]">
						retention
					</h2>
					<p class="m-0 text-[15.5px] leading-[1.65] text-[#c9c6c0] text-pretty">
						we keep your information for as long as your account is active. when you delete your
						account, your personal data and activity are permanently removed within 30 days, except
						where we're legally required to retain certain records. aggregated, anonymized
						statistics that can't identify you may be kept.
					</p>
				</div>

				<div id="contact" class="scroll-mt-23">
					<div class="mb-3 text-xs font-semibold text-[#8a857c] font-['JetBrains_Mono']">07</div>
					<h2 class="m-0 mb-4 text-[26px] leading-[1.1] font-semibold tracking-[-.02em]">
						contact
					</h2>
					<p class="m-0 mb-4.5 text-[15.5px] leading-[1.65] text-[#c9c6c0] text-pretty">
						questions about privacy, or want to exercise a right that isn't in settings? reach our
						privacy team - we aim to respond within a few business days.
					</p>
					<div class="flex items-center gap-4 rounded-xl border border-white/10 px-5.5 py-5">
						<div
							class="h-10 w-10 flex-none rounded-[10px]"
							style="background:oklch(.7 .16 245 / .18);border:1px solid oklch(.7 .16 245 / .5)"
						></div>
						<div>
							<div class="text-[14.5px] font-semibold">privacy@critiquefi.com</div>
							<div class="mt-1.25 text-xs text-[#8a857c] font-['JetBrains_Mono']">
								critiquefi, inc. · data protection team
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
