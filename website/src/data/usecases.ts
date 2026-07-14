import type { UseCase } from "./types";

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
