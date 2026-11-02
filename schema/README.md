# cmdguard/schema

Emit a [JSON Schema draft 2020-12](https://json-schema.org/draft/2020-12/schema)
document describing a cmdguard config struct or a whole CLI, from the same
struct tags that drive flags, environment variables, and strict config-key
checking.

**Experimental surface** — the API may still change before the first tagged
release.

## What you get

- One property per accepted config key, including nested objects and arrays
  (array element shapes under `items`).
- Constraints translated from `validate` tags (`min`/`max`/`pattern`,
  string length bounds), `values` tags as `enum`, and `default` tags as
  `default`.
- cmdguard facts JSON Schema has no home for under `x-cmdguard`: the flag,
  the env var, cmdguard-only formats (`duration`, `hostport`, `path`),
  and the `secret` marker.
- `additionalProperties: false` everywhere: a key the schema rejects is the
  same key strict config loading rejects. The schema and the binary can never
  disagree, because both are generated from `v4.ConfigKeys`.

## Usage

```go
package main

import (
    "os"

    "github.com/larsartmann/cmdguard/schema"
    "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func main() {
    cli, _ := v4.NewCLI[AppConfig]("myapp", "My app", AppConfig{})

    if err := schema.Emit(os.Stdout, cli.Config()); err != nil {
        panic(err)
    }
}
```

`FromConfig(cfg)` builds the document, `FromCLI(cli)` adds the CLI's name and
short description as title/description, `Emit` writes indented JSON to an
`io.Writer`, and `Marshal` returns it as bytes. Output is deterministic:
properties are emitted in sorted key order.

A typical pattern is a hidden `schema` subcommand (or a build-time generator)
that writes `config.schema.json` next to your binary.

## Editor and CI integration

Point your editor at the emitted file with a `$schema` reference, or pass it to
a validator in CI:

```yaml
# GitHub Actions example using check-jsonschema (pip install check-jsonschema)
- name: Validate example config
  run: check-jsonschema --schemafile config.schema.json config.example.yaml
```

Committing an emitted schema lets consumers, editors, and CI validate configs
without running your binary.

## Development

```bash
cd schema
go test ./...
# regenerate the golden fixture after changing the emitter:
go test ./... -run TestFromConfig_Golden -update-golden
```

The golden test (`testdata/config.schema.json`) fails on any emitter change;
regenerate intentionally and review the diff. If the test reports schema drift
on an unchanged emitter, suspect nondeterministic output — plain Go maps
marshal in random order under `encoding/json/v2`, which is why properties use
the sorted `PropertyMap`.
