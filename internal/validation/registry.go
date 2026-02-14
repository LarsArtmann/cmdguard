// Package validation provides command and flag validation for cmdguard.
package validation

import (
	"fmt"
	"sync"

	"github.com/larsartmann/cmdguard/internal/config"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CommandInfo holds metadata about a registered command.
type CommandInfo struct {
	Name           string
	Handler        func(cmd *cobra.Command, args []string) error
	HasSubcommands bool
	Flags          []*FlagInfo
}

// FlagInfo holds metadata about a registered flag.
type FlagInfo struct {
	Name       string
	Shorthand  string
	Type       string
	IsBound    bool
	IsRequired bool
	Default    interface{}
}

// Registry tracks commands and flags for validation.
type Registry struct {
	mu       sync.RWMutex
	commands map[string]*CommandInfo
	flags    map[string][]*FlagInfo
	cfg      *config.Config
}

// NewRegistry creates a new validation registry.
func NewRegistry(i do.Injector) (*Registry, error) {
	cfg, err := do.Invoke[*config.Config](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke config: %w", err)
	}

	return &Registry{
		commands: make(map[string]*CommandInfo),
		flags:    make(map[string][]*FlagInfo),
		cfg:      cfg,
	}, nil
}

// RegisterCommand adds a command to the registry.
func (r *Registry) RegisterCommand(cmd *cobra.Command) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info := &CommandInfo{
		Name:           cmd.Name(),
		HasSubcommands: len(cmd.Commands()) > 0,
	}

	// Capture handler if present
	if cmd.RunE != nil {
		info.Handler = cmd.RunE
	} else if cmd.Run != nil {
		// Wrap Run in a function matching Handler signature
		runFunc := cmd.Run
		info.Handler = func(c *cobra.Command, args []string) error {
			runFunc(c, args)
			return nil
		}
	}

	// Register flags for this command
	info.Flags = r.extractFlags(cmd)
	r.flags[cmd.Name()] = info.Flags

	r.commands[cmd.Name()] = info
	return nil
}

// RegisterSubcommand adds a subcommand under a parent.
func (r *Registry) RegisterSubcommand(parent *cobra.Command, child *cobra.Command) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	fullName := parent.Name() + " " + child.Name()
	info := &CommandInfo{
		Name:           fullName,
		HasSubcommands: len(child.Commands()) > 0,
	}

	// Capture handler
	if child.RunE != nil {
		info.Handler = child.RunE
	} else if child.Run != nil {
		runFunc := child.Run
		info.Handler = func(c *cobra.Command, args []string) error {
			runFunc(c, args)
			return nil
		}
	}

	// Register flags
	info.Flags = r.extractFlags(child)
	r.flags[fullName] = info.Flags

	r.commands[fullName] = info
	return nil
}

// GetCommand retrieves command info by name.
func (r *Registry) GetCommand(name string) (*CommandInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.commands[name]
	return info, ok
}

// GetCommands returns all registered commands.
func (r *Registry) GetCommands() map[string]*CommandInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy
	result := make(map[string]*CommandInfo, len(r.commands))
	for k, v := range r.commands {
		result[k] = v
	}
	return result
}

// GetFlags returns all flags for a command.
func (r *Registry) GetFlags(commandName string) ([]*FlagInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	flags, ok := r.flags[commandName]
	return flags, ok
}

// extractFlags extracts flag information from a cobra command.
func (r *Registry) extractFlags(cmd *cobra.Command) []*FlagInfo {
	var flags []*FlagInfo

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		flag := &FlagInfo{
			Name:       f.Name,
			Shorthand:  f.Shorthand,
			Type:       f.Value.Type(),
			IsRequired: len(f.Annotations[cobra.BashCompOneRequiredFlag]) > 0,
			Default:    f.DefValue,
			IsBound:    true, // Assume bound, validate later
		}
		flags = append(flags, flag)
	})

	return flags
}

// MarkFlagBound marks a flag as properly bound.
func (r *Registry) MarkFlagBound(commandName, flagName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if flags, ok := r.flags[commandName]; ok {
		for _, f := range flags {
			if f.Name == flagName {
				f.IsBound = true
				return
			}
		}
	}
}

// UnmarkFlagBound marks a flag as not bound (for validation).
func (r *Registry) UnmarkFlagBound(commandName, flagName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if flags, ok := r.flags[commandName]; ok {
		for _, f := range flags {
			if f.Name == flagName {
				f.IsBound = false
				return
			}
		}
	}
}

// HealthCheck implements the Healthchecker interface.
func (r *Registry) HealthCheck() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check for commands without handlers
	for name, cmd := range r.commands {
		if cmd.Handler == nil && !cmd.HasSubcommands {
			return fmt.Errorf("command %q has no handler", name)
		}
	}

	return nil
}

// CommandCount returns the number of registered commands.
func (r *Registry) CommandCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.commands)
}

// FlagCount returns the total number of registered flags.
func (r *Registry) FlagCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, flags := range r.flags {
		count += len(flags)
	}
	return count
}
