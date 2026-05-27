# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| v2.x    | :white_check_mark: |
| v1.x    | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in cmdguard, please report it privately.

**Do not open a public issue.**

Instead, send an email to **lars@larsartmann.com** with:

- A description of the vulnerability
- Steps to reproduce (if applicable)
- Potential impact
- Suggested fix (if you have one)

You can expect an initial response within 48 hours. If the vulnerability is confirmed, we will work with you to coordinate a fix and disclosure timeline.

## Security Best Practices

When using cmdguard in production:

- Validate all inputs with `required:"true"` or `WithPreRunE` hooks
- Use `WithConfigValidation` to enforce invariants after parsing
- Keep dependencies up to date (run `go mod tidy` regularly)
- Avoid logging sensitive flag values (passwords, tokens) unless explicitly required

## Disclosure Policy

We follow a coordinated disclosure process:

1. Reporter submits vulnerability privately
2. Maintainers acknowledge receipt within 48 hours
3. Maintainers investigate and develop a fix
4. Fix is released in a patch version
5. Public disclosure after users have had reasonable time to upgrade
