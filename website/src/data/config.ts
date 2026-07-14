export const siteConfig = {
  name: "cmdguard",
  title: "cmdguard — Type-safe CLI Framework for Go",
  description:
    "Build production Go CLIs with type-safe flags, dependency injection, and zero panics. Wraps Cobra with struct-tag-driven configuration.",
  ogDescription:
    "Type-safe CLI framework for Go. Struct-tag flags, dependency injection, constructor validation, and zero panics. Wraps Cobra.",
  siteUrl: "https://cmdguard.lars.software",
  github: "https://github.com/LarsArtmann/cmdguard",
  author: {
    name: "LarsArtmann",
    url: "https://larsartmann.com/",
  },
  pkgGoDev: "https://pkg.go.dev/github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3",
} as const;
