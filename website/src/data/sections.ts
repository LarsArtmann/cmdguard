import type { ComparisonRow, UseCase } from "./types";

export const comparisons: ComparisonRow[] = [
  {
    label: "Struct-tag flags",
    cobra: "String lookups",
    kong: "Yes",
    urfave: "String lookups",
    cmdguard: "Yes, validated at construction",
  },
  {
    label: "Dependency injection",
    cobra: "None",
    kong: "None",
    urfave: "None",
    cmdguard: "Lazy services, lifecycle, health checks",
  },
  {
    label: "Graceful shutdown",
    cobra: "DIY",
    kong: "DIY",
    urfave: "DIY",
    cmdguard: "Reverse-order on SIGINT/SIGTERM",
  },
  {
    label: "Zero panics",
    cobra: "Run panics, RunE returns",
    kong: "Some panic variants",
    urfave: "Some panic variants",
    cmdguard: "Only error returns, by construction",
  },
  {
    label: "Error printing",
    cobra: "Prints twice (cobra + main)",
    kong: "Manual",
    urfave: "Manual",
    cmdguard: "Exactly once, correct exit code",
  },
  {
    label: "Usage on error",
    cobra: "Prints full usage (footgun)",
    kong: "Configurable",
    urfave: "Configurable",
    cmdguard: "Silenced by default",
  },
  {
    label: "Styled output",
    cobra: "Plain",
    kong: "Plain",
    urfave: "Plain",
    cmdguard: "fang + lipgloss by default",
  },
  {
    label: "Cobra compatibility",
    cobra: "Native",
    kong: "No",
    urfave: "No",
    cmdguard: "Escape hatch + gradual migration",
  },
];

export const useCases: UseCase[] = [
  {
    title: "Production CLIs",
    desc: "Ship robust command-line tools with proper error handling and exit codes out of the box.",
    icon: "cog",
  },
  {
    title: "DevTools",
    desc: "Build developer utilities with typed flags, config loading, and rich output formats.",
    icon: "chart",
  },
  {
    title: "Microservice CLIs",
    desc: "Wire up DI services, health checks, graceful shutdown, and telemetry in one framework.",
    icon: "refresh",
  },
];
