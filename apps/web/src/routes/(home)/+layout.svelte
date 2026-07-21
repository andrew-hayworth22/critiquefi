<script lang="ts">
	import Footer from './_components/Footer.svelte';
	import Nav from './_components/Nav.svelte';

	let { children } = $props();
</script>

<div
	class="min-h-screen bg-[#0e0e10] text-[#f2f0ec] antialiased font-['Space_Grotesk',system-ui,sans-serif]"
>
	<div class="mx-auto max-w-310">
		<Nav />
		<main class="pt-24">
			{@render children()}
		</main>
		<Footer />
	</div>
</div>
