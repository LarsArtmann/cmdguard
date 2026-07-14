import type { ComparisonRow, UseCase } from "./types";

export const comparisons: ComparisonRow[] = [
  {
    label: "Usage on error",
    cobra: "Prints full usage after every error",
    cmdguard: "Silenced by default — help still works",
    cobraBad: true,
  },
  {
    label: "Error output",
    cobra: "Prints twice (cobra + main)",
    cmdguard: "Printed exactly once",
    cobraBad: true,
  },
  {
    label: "Exit codes",
    cobra: "Manual os.Exit with right code",
    cmdguard: "ExecuteAndExit handles it",
    cobraBad: true,
  },
  {
    label: "Panics",
    cobra: "Run panics, RunE returns error",
    cmdguard: "Only error-returning handlers",
    cobraBad: true,
  },
  {
    label: "Flag lookups",
    cobra: "String-based, fail at runtime",
    cmdguard: "Typed structs, validated at construction",
    cobraBad: true,
  },
  {
    label: "Missing handler",
    cobra: "Found at runtime",
    cmdguard: "Caught at AddCommand time",
    cobraBad: true,
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
