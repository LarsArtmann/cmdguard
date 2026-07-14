import type { ComparisonRow } from "./types";

export const comparisons: ComparisonRow[] = [
  {
    label: "Struct-tag flags",
    cobra: { level: "diy", note: "String lookups" },
    kong: { level: "full", note: "Yes" },
    urfave: { level: "diy", note: "String lookups" },
    cmdguard: { level: "full", note: "Validated at construction" },
  },
  {
    label: "Dependency injection",
    cobra: { level: "none", note: "—" },
    kong: { level: "none", note: "—" },
    urfave: { level: "none", note: "—" },
    cmdguard: { level: "full", note: "Lazy services, lifecycle, health checks" },
  },
  {
    label: "Graceful shutdown",
    cobra: { level: "diy", note: "Manual signal handling" },
    kong: { level: "diy", note: "Manual signal handling" },
    urfave: { level: "diy", note: "Manual signal handling" },
    cmdguard: { level: "full", note: "Reverse-order on SIGINT/SIGTERM" },
  },
  {
    label: "Zero panics",
    cobra: { level: "partial", note: "Run panics, RunE returns" },
    kong: { level: "partial", note: "Some panic variants" },
    urfave: { level: "partial", note: "Some panic variants" },
    cmdguard: { level: "full", note: "Only error returns, by construction" },
  },
  {
    label: "Error printing",
    cobra: { level: "diy", note: "Prints twice if not careful" },
    kong: { level: "diy", note: "Manual" },
    urfave: { level: "diy", note: "Manual" },
    cmdguard: { level: "full", note: "Exactly once, correct exit code" },
  },
  {
    label: "Usage on error",
    cobra: { level: "partial", note: "Prints full usage (footgun)" },
    kong: { level: "full", note: "Configurable" },
    urfave: { level: "full", note: "Configurable" },
    cmdguard: { level: "full", note: "Silenced by default" },
  },
  {
    label: "Styled output",
    cobra: { level: "none", note: "Plain by default" },
    kong: { level: "none", note: "Plain by default" },
    urfave: { level: "none", note: "Plain by default" },
    cmdguard: { level: "full", note: "fang + lipgloss by default" },
  },
  {
    label: "Cobra compatibility",
    cobra: { level: "native", note: "Native" },
    kong: { level: "none", note: "No" },
    urfave: { level: "none", note: "No" },
    cmdguard: { level: "full", note: "Escape hatch + gradual migration" },
  },
];
