import type { Feature } from "./types";

export const features: Feature[] = [
  {
    icon: "shield",
    title: "Zero Panics",
    desc: "Every function returns errors. No Must* panic variants exist in the library.",
  },
  {
    icon: "type",
    title: "Type-Safe Flags",
    desc: "Struct tags define flags — no stringly-typed GetString lookups that fail at runtime.",
  },
  {
    icon: "syringe",
    title: "Dependency Injection",
    desc: "Built-in samber/do/v2 integration with lifecycle hooks, health checks, and graceful shutdown.",
  },
  {
    icon: "terminal",
    title: "Constructor Validation",
    desc: "Missing handlers, duplicate commands, invalid flags caught at AddCommand time, not runtime.",
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
