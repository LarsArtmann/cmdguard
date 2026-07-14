import type { Feature } from "./types";

export const features: Feature[] = [
  {
    icon: "syringe",
    title: "Dependency Injection",
    desc: "Lazy services, lifecycle hooks, health checks, and graceful shutdown in reverse order. Powered by samber/do/v2.",
  },
  {
    icon: "shield",
    title: "Zero Panics",
    desc: "Every function returns errors. No Run, no Must*. Only error-returning handlers exist — by construction.",
  },
  {
    icon: "terminal",
    title: "Validated at Construction",
    desc: "Missing handlers, duplicate commands, invalid flags — caught at AddCommand time, not at runtime.",
  },
  {
    icon: "type",
    title: "Type-Safe Flags",
    desc: "Struct tags define flags with env vars, defaults, required, counting, and scoping. No string lookups.",
  },
  {
    icon: "file",
    title: "Config Files",
    desc: "JSON, YAML, TOML auto-loading with env var override and flag precedence chain.",
  },
  {
    icon: "bolt",
    title: "16 Output Formats",
    desc: "Table, JSON, CSV, YAML, Markdown, XML, HTML, D2, Mermaid, and more via go-output.",
  },
  {
    icon: "tag",
    title: "Rich Flag Tags",
    desc: "env, required, count, short, default, local, hidden — all via struct tags.",
  },
  {
    icon: "git",
    title: "Cobra Escape Hatch",
    desc: "Mix raw cobra.Command subcommands with cmdguard's typed runtime. Gradual migration supported.",
  },
];
