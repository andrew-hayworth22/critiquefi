<script lang="ts">
	import { resolve } from '$app/paths';

	const navLinks = [
		{ href: resolve('/') + '#how', label: 'how it works' },
		{ href: resolve('/') + '#feed', label: 'the feed' },
		{ href: resolve('/'), label: 'log in' }
	];
</script>

<nav
	class="sticky top-0 z-50 flex items-center justify-between border-b border-white/[.07] bg-[#0e0e10]/72 px-12 py-5 backdrop-blur-md"
>
	<div class="text-[22px] leading-none font-bold tracking-[-.02em]">critiquefi</div>
	<div class="flex items-center gap-7.5 text-sm text-[#a5a29c]">
		{#each navLinks as link (link.label)}
			<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
			<a href={link.href} class="hover:text-[#f2f0ec]">{link.label}</a>
		{/each}
		<a href="#" class="rounded-lg bg-[#f2f0ec] px-5 py-2.5 font-semibold text-[#0e0e10]"
			>sign up free</a
		>
	</div>
</nav>
