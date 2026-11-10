package main

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
	"github.com/spf13/pflag"
)

type Env string

func ParseEnv(s string) (Env, error) {
	switch s {
	case "dev", "staging", "prod":
		return Env(s), nil
	}
	return "", fmt.Errorf("unknown environment %q", s)
}

type DeployFlags struct {
	Target Env `flag:"target" default:"dev" help:"Deployment target"`
}

type AppConfig struct{}

func main() {
	err := v5.RegisterTypeCodec(reflect.TypeFor[Env](), v5.TypeCodecFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag v5.FlagTag) error {
			flags.StringP(tag.Name, tag.Short, tag.Default, tag.Help)
			return nil
		},
		ParseFunc: func(value string, _ v5.FlagTag) (any, error) {
			return ParseEnv(value)
		},
		DefaultFunc: func(tag v5.FlagTag) any {
			env, _ := ParseEnv(tag.Default)
			return env
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	cli, err := v5.NewCLI[AppConfig]("deploy", "Custom type example", AppConfig{})
	if err != nil {
		log.Fatal(err)
	}

	b := v5.NewBuilder(cli)
	v5.Add(b, "plan", DeployFlags{},
		func(_ context.Context, _ *AppConfig, f DeployFlags) error {
			fmt.Println("target:", f.Target)
			return nil
		},
		v5.WithShort("Plan a deployment"),
	)
	if err := b.Build(); err != nil {
		log.Fatal(err)
	}

	cli.ExecuteAndExit(context.Background())
}
