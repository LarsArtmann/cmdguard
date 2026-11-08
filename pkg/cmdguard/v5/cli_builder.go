package v5

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// Builder collects command-construction results and fails once at Build,
// killing the per-command `cmd, err := NewCommand(...); if err != nil` ladder.
//
//	b := v5.NewBuilder(cli)
//	v5.Add(b, "list", &ListFlags{}, listHandler, v5.WithShort("List tasks"))
//	v5.Add(b, "add", &AddFlags{}, addHandler, v5.WithShort("Add a task"))
//	if err := b.Build(); err != nil {
//		return err
//	}
//
// Add runs the constructor immediately: a nil error queues the command, a
// non-nil error is collected and the call moves on. Build attaches every
// queued command in one atomic AddCommands pass and reports all collected
// errors joined, so a construction failure names every broken command, not
// just the first. Callers that already hold a (Command, error) pair use Try.
type Builder[T any] struct {
	cli  *CLI[T]
	cmds []AddableCommand[T]
	errs []error
}

// NewBuilder starts collecting command-construction results for cli.
func NewBuilder[T any](cli *CLI[T]) *Builder[T] {
	return &Builder[T]{cli: cli}
}

// add queues cmd when err is nil, otherwise collects err.
func (b *Builder[T]) add(cmd AddableCommand[T], err error) {
	if err != nil {
		b.errs = append(b.errs, err)

		return
	}

	if cmd == nil {
		b.errs = append(b.errs, ErrInvalidCommand)

		return
	}

	b.cmds = append(b.cmds, cmd)
}

// Add constructs one command on b and queues it, collecting any construction
// error. It is a package-level function (not a Builder method) because Go
// methods cannot carry the per-command flags type parameter F.
func Add[T, F any](
	b *Builder[T],
	use string,
	flags F,
	runE func(context.Context, *T, F) error,
	opts ...CommandOption,
) {
	if b == nil {
		return
	}

	cmd, err := NewCommand(use, flags, runE, opts...)
	b.add(cmd, err)
}

// AddParent constructs a parent command on b and queues it, collecting any
// construction error.
func AddParent[T, F any](
	b *Builder[T],
	use, long string,
	flags F,
	opts ...CommandOption,
) {
	if b == nil {
		return
	}

	cmd, err := NewParentCommand[T](use, long, flags, opts...)
	b.add(cmd, err)
}

// Check collects a construction error without queueing anything. Use it for
// commands whose Command value is consumed by a later constructor (children
// passed to WithSubcommands) so they must not be queued for attachment twice:
//
//	child, err := v5.NewCommand(...)
//	b.Check(err)
//	v5.AddParent(b, "db", "Database operations", flags, v5.WithSubcommands(child))
func (b *Builder[T]) Check(err error) {
	if b == nil {
		return
	}

	if err != nil {
		b.errs = append(b.errs, err)
	}
}

// Try queues a (Command, error) pair the caller already holds: a nil error
// queues the command, a non-nil error is collected and the zero command
// dropped.
func Try[T, F any](b *Builder[T], cmd Command[T, F], err error) {
	if b == nil {
		return
	}

	b.add(cmd, err)
}

// Build attaches every queued command atomically and reports every failure
// joined: collected construction errors plus every prepare/validation error
// found across the batch (prepare is side-effect-free, so all commands are
// checked even after one fails). On error, nothing is attached.
func (b *Builder[T]) Build() error {
	errs := append([]error(nil), b.errs...)

	prepared := make([]*cobra.Command, 0, len(b.cmds))
	seen := make(map[string]struct{}, len(b.cmds))

	for _, cmd := range b.cmds {
		cobraCmd, err := cmd.prepareFor(b.cli)
		if err != nil {
			errs = append(errs, err)

			continue
		}

		name := cobraCmd.Name()

		if _, duplicate := seen[name]; duplicate {
			errs = append(errs,
				fmt.Errorf("%w: command %q appears twice in this Builder", ErrDuplicateCommand, name),
			)

			continue
		}

		seen[name] = struct{}{}

		prepared = append(prepared, cobraCmd)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	for _, cobraCmd := range prepared {
		b.cli.attachCommand(cobraCmd)
	}

	return nil
}
