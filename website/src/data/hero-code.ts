const importPath = "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3";

export const heroCode = `package main

import (
    "context"
    "fmt"
    "strings"

    "${importPath}"
)

type AppConfig struct {
    Verbose bool   \`flag:"verbose" short:"v" default:"false" help:"Enable verbose output"\`
    Output  string \`flag:"output" short:"o" default:"text" help:"Output format"\`
}

type GreetFlags struct {
    Name  string \`flag:"name" short:"n" default:"World" help:"Name to greet"\`
    Shout bool   \`flag:"shout" short:"s" default:"false" help:"Uppercase output"\`
}

func main() {
    cli, err := v3.NewCLI[AppConfig]("myapp", "My CLI app", AppConfig{})
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed: %v\\n", err)
        os.Exit(1)
    }

    greetCmd, _ := v3.NewCommand("greet", &GreetFlags{},
        func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
            msg := fmt.Sprintf("Hello, %s!", flags.Name)
            if flags.Shout { msg = strings.ToUpper(msg) }
            fmt.Println(msg)
            return nil
        },
        v3.WithShort("Greet someone"),
    )
    v3.AddCommand(cli, greetCmd)
    cli.ExecuteAndExit(context.Background())
}`;
