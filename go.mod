module github.com/larsartmann/cmdguard/v3

go 1.26.5

require (
	charm.land/fang/v2 v2.0.1
	charm.land/lipgloss/v2 v2.0.5
	github.com/go-faster/yaml v0.4.6
	github.com/knadh/koanf/parsers/json v1.0.0
	github.com/knadh/koanf/parsers/yaml v1.1.0
	github.com/knadh/koanf/providers/file v1.2.1
	github.com/knadh/koanf/v2 v2.3.5
	github.com/larsartmann/go-output v0.31.1
	github.com/larsartmann/go-output/d2 v0.31.1 // indirect
	github.com/larsartmann/go-output/delimited v0.31.1 // indirect
	github.com/larsartmann/go-output/graph v0.31.1 // indirect
	github.com/larsartmann/go-output/markdown v0.31.1 // indirect
	github.com/larsartmann/go-output/markup v0.31.1 // indirect
	github.com/larsartmann/go-output/plantuml v0.31.1 // indirect
	github.com/larsartmann/go-output/serialization v0.31.1 // indirect
	github.com/larsartmann/go-output/table v0.31.1 // indirect
	github.com/larsartmann/go-output/tree v0.31.1 // indirect
	github.com/larsartmann/samber-do-auditlog v0.7.0
	github.com/pelletier/go-toml/v2 v2.4.3
	github.com/samber/do/v2 v2.1.0
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.10
	golang.org/x/term v0.45.0 // indirect
)

require github.com/larsartmann/go-output/escape v0.31.1 // indirect

require (
	charm.land/glamour/v2 v2.0.1 // indirect
	github.com/a-h/templ v0.3.1020 // indirect
	github.com/alecthomas/chroma/v2 v2.27.0 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/charmbracelet/colorprofile v0.4.3 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20260720091822-7cc6674724ac // indirect
	github.com/charmbracelet/x/ansi v0.11.7 // indirect
	github.com/charmbracelet/x/exp/charmtone v0.0.0-20260713092006-0d683c34c74b // indirect
	github.com/charmbracelet/x/exp/slice v0.0.0-20260713092006-0d683c34c74b // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/dlclark/regexp2/v2 v2.5.0 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/go-faster/errors v0.8.0 // indirect
	github.com/go-faster/jx v1.2.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/gofrs/flock v0.13.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/knadh/koanf/parsers/toml v0.1.0 // indirect
	github.com/larsartmann/go-atomic-write v0.3.0 // indirect
	github.com/larsartmann/go-branded-id v0.3.2 // indirect
	github.com/larsartmann/go-error-family v0.9.0 // indirect
	github.com/larsartmann/go-ndjson v0.0.1 // indirect
	github.com/larsartmann/go-output/daghtml v0.31.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.27 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/mango v0.2.0 // indirect
	github.com/muesli/mango-cobra v1.3.0 // indirect
	github.com/muesli/mango-pflag v0.2.0 // indirect
	github.com/muesli/roff v0.1.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.15.0 // indirect
	github.com/samber/go-type-to-string v1.8.0 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/yuin/goldmark v1.8.4 // indirect
	github.com/yuin/goldmark-emoji v1.0.6 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
)

replace github.com/larsartmann/cmdguard/glamour => ./glamour

require (
	github.com/larsartmann/cmdguard/glamour v0.0.0-00010101000000-000000000000
	github.com/larsartmann/cmdguard/spinner v0.0.0-00010101000000-000000000000
)

replace github.com/larsartmann/cmdguard/prompts => ./prompts

replace github.com/larsartmann/cmdguard/telemetry => ./telemetry

replace github.com/larsartmann/cmdguard/spinner => ./spinner
