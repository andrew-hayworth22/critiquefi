<script lang="ts">
	import { resolve } from '$app/paths';

	const footerColumns = [
		{
			title: 'product',
			links: [
				{ href: resolve('/') + '#features', label: 'features' },
				{ href: resolve('/') + '#how', label: 'how it works' },
				{ href: resolve('/') + '#feed', label: 'the feed' }
			]
		},
		{
			title: 'company',
			links: [
				{ href: '#', label: 'about' },
				{ href: '#', label: 'careers' },
				{ href: '#', label: 'blog' }
			]
		},
		{
			title: 'legal',
			links: [
				{ href: resolve('/(home)/(legal)/privacy'), label: 'privacy' },
				{ href: resolve('/(home)/(legal)/terms'), label: 'terms' }
			]
		}
	];
</script>

<div class="flex flex-wrap justify-between gap-6 border-t border-white/8 p-12">
	<div>
		<div class="mb-2.5 text-xl leading-none font-bold tracking-[-.02em]">critiquefi</div>
		<div class="text-xs leading-[1.6] text-[#6b6863] font-['JetBrains_Mono']">
			one shelf for everything.<br />web now · apps coming.
		</div>
	</div>
	<div class="flex gap-14 text-sm text-[#a5a29c]">
		{#each footerColumns as col (col.title)}
			<div class="flex flex-col gap-2.75">
				<span
					class="mb-0.5 text-[11px] text-[#6b6863] uppercase font-['JetBrains_Mono'] tracking-widest"
					>{col.title}</span
				>
				{#each col.links as link (link.label)}
					<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
					<a href={link.href} class="hover:text-[#f2f0ec]">{link.label}</a>
				{/each}
			</div>
		{/each}
	</div>
</div>
