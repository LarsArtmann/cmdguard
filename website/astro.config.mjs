import { defineConfig, fontProviders } from "astro/config";
import starlight from "@astrojs/starlight";
import sitemap from "@astrojs/sitemap";

import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
	site: "https://cmdguard.lars.software",
	security: {
		csp: {
			scriptDirective: {
				resources: ["'self'"],
			},
			styleDirective: {
				resources: ["'self'", "'unsafe-inline'"],
			},
		},
	},

	compressHTML: true,
	trailingSlash: "always",

	prefetch: {
		prefetchAll: false,
		defaultStrategy: "hover",
	},

	fonts: [
		{
			provider: fontProviders.google(),
			name: "Space Grotesk",
			cssVariable: "--font-space-grotesk",
			weights: [300, 400, 500, 600, 700],
			styles: ["normal"],
			subsets: ["latin"],
			fallbacks: ["sans-serif"],
		},
		{
			provider: fontProviders.fontsource(),
			name: "JetBrains Mono",
			cssVariable: "--font-jetbrains-mono",
			weights: [400, 500, 600, 700],
			styles: ["normal"],
			subsets: ["latin"],
			fallbacks: ["monospace"],
		},
	],

	integrations: [
		sitemap(),
		starlight({
			title: "cmdguard",
			favicon: "/favicon.svg",
			customCss: ["./src/styles/starlight.css"],
			lastUpdated: true,
			editLink: {
				baseUrl: "https://github.com/LarsArtmann/cmdguard/edit/master/website/src/content/docs/",
			},
			expressiveCode: {
				themes: ["github-light", "github-dark"],
				frames: {
					showCopyToClipboardButton: true,
				},
			},
			sidebar: [
				{
					label: "Getting Started",
					items: [
						{ label: "Installation", slug: "getting-started/installation" },
						{ label: "Quick Start", slug: "getting-started/quick-start" },
					],
				},
				{
					label: "Core Concepts",
					items: [
						{ label: "Type-Safe Flags", slug: "guides/type-safe-flags" },
						{ label: "Custom Value Types", slug: "guides/custom-types" },
						{ label: "Dependency Injection", slug: "guides/dependency-injection" },
						{ label: "Lifecycle & Signals", slug: "guides/lifecycle" },
						{ label: "Production Batteries", slug: "guides/production-batteries" },
					],
				},
				{
					label: "Guides",
					items: [
						{ label: "Error Handling", slug: "guides/error-handling" },
						{ label: "Config Files", slug: "guides/config-files" },
						{ label: "Rich Output", slug: "guides/rich-output" },
						{ label: "Middleware", slug: "guides/middleware" },
						{ label: "Flow Context", slug: "guides/flow-context" },
						{ label: "Plugin System", slug: "guides/plugins" },
						{ label: "Audit Log", slug: "guides/audit-log" },
						{ label: "Shell Completion", slug: "guides/shell-completion" },
						{ label: "Machine Contract", slug: "guides/machine-contract" },
					],
				},
				{
					label: "Sub-Modules & Tooling",
					items: [
						{ label: "Sub-Modules", slug: "guides/sub-modules" },
						{ label: "Schema Generation", slug: "guides/schema" },
						{ label: "Usage Linter", slug: "guides/lint" },
						{ label: "Doctor & Health Checks", slug: "guides/doctor" },
					],
				},
				{
					label: "Migration",
					items: [
						{ label: "Migrating from Cobra", slug: "guides/migrating-from-cobra" },
						{ label: "Migrating v4 to v5", slug: "guides/migrating-v4-to-v5" },
					],
				},
				{
					label: "API Reference",
					items: [
						{ label: "Overview", slug: "api-reference" },
						{
							label: "Full API on pkg.go.dev",
							link: "https://pkg.go.dev/github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5",
						},
					],
				},
				{
					label: "Community",
					items: [
						{ label: "Changelog", slug: "changelog" },
						{ label: "Contributing", slug: "contributing" },
						{ label: "Related Tools", slug: "related-tools" },
					],
				},
			],
			social: [
				{
					icon: "github",
					label: "GitHub",
					href: "https://github.com/LarsArtmann/cmdguard",
				},
			],
			head: [
				{
					tag: "meta",
					attrs: {
						name: "description",
						content:
							"The only Go CLI framework that unifies type-safe flags, dependency injection with lifecycle management, and zero-panic error contracts. Wraps Cobra.",
					},
				},
			],
		}),
	],

	vite: {
		plugins: [tailwindcss()],
	},
});
