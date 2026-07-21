<script lang="ts">
	interface NavItem {
		id: string;
		label: string;
	}

	interface Props {
		navItems: NavItem[];
	}

	let { navItems }: Props = $props();
</script>

<div class="sticky top-23 flex flex-col gap-0.5">
	<div
		class="px-3 pb-3 text-[10.5px] text-[#6b6863] uppercase font-['JetBrains_Mono'] tracking-[.12em]"
	>
		on this page
	</div>
	{#each navItems as item, i (item.id)}
		<a
			href="#{item.id}"
			class="rounded-lg px-3 py-2.25 text-[13.5px] no-underline text-[#c9c6c0] hover:text-[#f2f0ec]"
		>
			{i + 1} · {item.label}
		</a>
	{/each}
</div>
