package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

type Config struct {
	Name    string `flag:"name"    short:"n" default:"World" help:"Name to greet"`
	Shout   bool   `flag:"shout"            default:"false" help:"Shout the greeting"`
	Count   int    `flag:"count"   short:"c" default:"1"    help:"Number of greetings"`
	ConfigFile string `flag:"config"  short:"C" default:""      help:"Path to config file"`
}

func main() {
	cli, err := v2.NewCLI[Config]("greet", "Greeting CLI", Config{},
		v2.WithConfigFile[Config]("$HOME/.config/greet/config.json"),
	)
	if err != nil {
		panic(err)
	}

	cmd, err := v2.NewCommand[Config, v2.NoFlags]("hello",
		func(_ context.Context, cfg *Config, _ v2.NoFlags) error {
			for i := 0; i < cfg.Count; i++ {
				msg := fmt.Sprintf("Hello, %s!", cfg.Name)
				if cfg.Shout {
					msg = strings.ToUpper(msg)
				}
				fmt.Println(msg)
			}
			return nil
		},
		v2.WithShort[Config, v2.NoFlags]("Say hello"),
	)
	if err != nil {
		panic(err)
	}

	if err := v2.AddCommand(cli, cmd); err != nil {
		panic(err)
	}

	if err := cli.Execute(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
