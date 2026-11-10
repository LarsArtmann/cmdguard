package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type AppConfig struct {
	Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
	Config  string `flag:"config" help:"Path to the config file" local:"true"`
}

type GreetFlags struct {
	Name  string `flag:"name"  short:"n" default:"World" help:"Name to greet"`
	Loud  bool   `flag:"loud"  default:"false" help:"Uppercase the greeting"`
	Times int    `flag:"times" default:"1" help:"Repetitions" validate:"min=1,max=100"`
}

func main() {
	cli, err := v5.NewCLI[AppConfig]("demo", "A greeting CLI", AppConfig{})
	if err != nil {
		log.Fatal(err)
	}

	b := v5.NewBuilder(cli)
	v5.Add(b, "greet", GreetFlags{},
		func(_ context.Context, _ *AppConfig, f GreetFlags) error {
			msg := fmt.Sprintf("Hello, %s!", f.Name)
			if f.Loud {
				msg = strings.ToUpper(msg)
			}
			for range f.Times {
				fmt.Println(msg)
			}
			return nil
		},
		v5.WithShort("Greet someone"),
		v5.WithExample("demo greet --name Alice --times 2"),
	)
	if err := b.Build(); err != nil {
		log.Fatal(err)
	}

	cli.ExecuteAndExit(context.Background())
}
