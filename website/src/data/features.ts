import type { Feature } from "./types";

export const features: Feature[] = [
  {
    icon: "shield",
    title: "Batteries On By Default",
    desc: "Strict validation, signal handling, graceful DI shutdown, strict config keys, error hints, and doctor/version are enabled at NewCLI. Each has a named opt-out.",
  },
  {
    icon: "type",
    title: "Type-Safe Flags",
    desc: "Struct tags define the flag name, shorthand, default, env var, required state, enum values, counting, scoping, validation, and secret redaction. No string lookups.",
  },
  {
    icon: "file",
    title: "Machine Contract",
    desc: "Every CLI answers describe with a versioned JSON contract: commands, flags, env vars, enums, required flags, and inheritance. Built for agents and tooling.",
  },
  {
    icon: "syringe",
    title: "Dependency Injection",
    desc: "Lazy services, lifecycle hooks, health checks, and dependency-order shutdown on a signal. Powered by samber/do/v2.",
  },
  {
    icon: "terminal",
    title: "Zero Panics, Real Errors",
    desc: "Every constructor returns an error. Classified errors print a Hint line, carry a stable error_code, and map to the right exit code.",
  },
  {
    icon: "tag",
    title: "Config You Can Explain",
    desc: "JSON, YAML, and TOML with flag > env > file > default precedence. Unknown keys fail with a suggestion, and config explain shows where each value came from.",
  },
  {
    icon: "bolt",
    title: "16 Output Formats",
    desc: "Table, JSON, CSV, TSV, YAML, TOML, Markdown, XML, HTML, tree, D2, DOT, Mermaid, PlantUML, JSONL, and AsciiDoc via go-output.",
  },
  {
    icon: "git",
    title: "Cobra Escape Hatch",
    desc: "Keep raw cobra.Command subcommands beside typed ones. Gradual migration, no rewrite, no lock-in.",
  },
];
