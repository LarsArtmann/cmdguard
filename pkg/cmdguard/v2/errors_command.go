package v2

import "errors"

// Command-related sentinel errors.
var (
	// ErrInvalidCommand indicates a command is malformed or missing required fields.
	ErrInvalidCommand = errors.New("invalid command")

	// ErrMissingHandler indicates a command has no RunE handler and no subcommands.
	ErrMissingHandler = errors.New("command has no handler")

	// ErrMissingName indicates a command has no name (Use field is empty).
	ErrMissingName = errors.New("command has no name")

	// ErrMissingLong indicates a parent command has no long description.
	ErrMissingLong = errors.New("command has no long description")

	// ErrMissingShort indicates a command has no short description.
	ErrMissingShort = errors.New("command has no short description")

	// ErrDuplicateCommand indicates a command with the same name already exists.
	ErrDuplicateCommand = errors.New("duplicate command")

	// ErrCommandPanic indicates a command handler panicked during execution.
	ErrCommandPanic = errors.New("command panicked")

	// ErrMissingExample indicates a leaf command has no example in draconian mode.
	ErrMissingExample = errors.New("command has no example")

	// ErrTooFewArgs indicates a command received fewer positional arguments than required.
	ErrTooFewArgs = errors.New("too few arguments")

	// ErrTooManyArgs indicates a command received more positional arguments than allowed.
	ErrTooManyArgs = errors.New("too many arguments")

	// ErrNegativeArgCount indicates a negative argument count was provided.
	ErrNegativeArgCount = errors.New("argument count must not be negative")

	// ErrInvalidArgRange indicates min > max in a range argument validator.
	ErrInvalidArgRange = errors.New("minimum argument count must not exceed maximum")
)
