module github.com/larsartmann/cmdguard

go 1.26.3

require (
	charm.land/fang/v2 v2.0.1
	charm.land/huh/v2 v2.0.3
	github.com/larsartmann/go-branded-id v0.3.0
	github.com/larsartmann/go-output v0.6.1
	github.com/larsartmann/go-output/d2 v0.6.1
	github.com/larsartmann/go-output/delimited v0.6.1
	github.com/larsartmann/go-output/graph v0.6.1
	github.com/larsartmann/go-output/markup v0.6.1
	github.com/larsartmann/go-output/serialization v0.6.1
	github.com/larsartmann/go-output/table v0.6.1
	github.com/muesli/mango v0.2.0
	github.com/muesli/mango-cobra v1.3.0
	github.com/muesli/roff v0.1.0
	github.com/pelletier/go-toml/v2 v2.3.1
	github.com/samber/do/v2 v2.0.0
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.10
	gopkg.in/yaml.v3 v3.0.1
)

require (
	charm.land/bubbles/v2 v2.1.0 // indirect
	charm.land/bubbletea/v2 v2.0.6 // indirect
	charm.land/lipgloss/v2 v2.0.3 // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/catppuccin/go v0.3.0 // indirect
	github.com/charmbracelet/colorprofile v0.4.3 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20260525132238-948f4557a654 // indirect
	github.com/charmbracelet/x/ansi v0.11.7 // indirect
	github.com/charmbracelet/x/exp/charmtone v0.0.0-20260531005911-0ca8ababeab2 // indirect
	github.com/charmbracelet/x/exp/ordered v0.1.0 // indirect
	github.com/charmbracelet/x/exp/strings v0.1.0 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-faster/jx v1.2.0 // indirect
	github.com/go-faster/yaml v0.4.6 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/larsartmann/go-output/enum v0.6.1 // indirect
	github.com/larsartmann/go-output/escape v0.6.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.24 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/mango-pflag v0.2.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/samber/go-type-to-string v1.8.0 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/term v0.43.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace (
	github.com/larsartmann/go-output => ../go-output
	github.com/larsartmann/go-output/d2 => ../go-output/d2
	github.com/larsartmann/go-output/delimited => ../go-output/delimited
	github.com/larsartmann/go-output/enum => ../go-output/enum
	github.com/larsartmann/go-output/escape => ../go-output/escape
	github.com/larsartmann/go-output/graph => ../go-output/graph
	github.com/larsartmann/go-output/markup => ../go-output/markup
	github.com/larsartmann/go-output/serialization => ../go-output/serialization
	github.com/larsartmann/go-output/table => ../go-output/table
	github.com/larsartmann/go-output/testhelpers => ../go-output/testhelpers
	github.com/larsartmann/go-output/testhelpers/graphtest => ../go-output/testhelpers/graphtest
)
