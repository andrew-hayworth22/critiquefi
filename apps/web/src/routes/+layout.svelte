<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import '@fontsource/space-grotesk/400.css';
	import '@fontsource/space-grotesk/500.css';
	import '@fontsource/space-grotesk/600.css';
	import '@fontsource/space-grotesk/700.css';
	import '@fontsource/jetbrains-mono/400.css';
	import '@fontsource/jetbrains-mono/500.css';
	import '@fontsource/jetbrains-mono/600.css';
	import spaceGrotesk400 from '@fontsource/space-grotesk/files/space-grotesk-latin-400-normal.woff2?url';
	import spaceGrotesk500 from '@fontsource/space-grotesk/files/space-grotesk-latin-500-normal.woff2?url';
	import spaceGrotesk600 from '@fontsource/space-grotesk/files/space-grotesk-latin-600-normal.woff2?url';
	import spaceGrotesk700 from '@fontsource/space-grotesk/files/space-grotesk-latin-700-normal.woff2?url';
	import jetbrainsMono400 from '@fontsource/jetbrains-mono/files/jetbrains-mono-latin-400-normal.woff2?url';
	import jetbrainsMono500 from '@fontsource/jetbrains-mono/files/jetbrains-mono-latin-500-normal.woff2?url';
	import jetbrainsMono600 from '@fontsource/jetbrains-mono/files/jetbrains-mono-latin-600-normal.woff2?url';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="preload" href={spaceGrotesk400} as="font" type="font/woff2" crossorigin="anonymous" />
	<link rel="preload" href={spaceGrotesk500} as="font" type="font/woff2" crossorigin="anonymous" />
	<link rel="preload" href={spaceGrotesk600} as="font" type="font/woff2" crossorigin="anonymous" />
	<link rel="preload" href={spaceGrotesk700} as="font" type="font/woff2" crossorigin="anonymous" />
	<link rel="preload" href={jetbrainsMono400} as="font" type="font/woff2" crossorigin="anonymous" />
	<link rel="preload" href={jetbrainsMono500} as="font" type="font/woff2" crossorigin="anonymous" />
	<link rel="preload" href={jetbrainsMono600} as="font" type="font/woff2" crossorigin="anonymous" />
</svelte:head>
{@render children()}
