// Package spinner provides a terminal spinner middleware for cmdguard CLIs.
// It is an optional module — import it only when you need a visual spinner,
// to keep your dependency tree lean.
//
// Usage:
//
//	import (
//	    v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
//	    "github.com/larsartmann/cmdguard/spinner"
//	)
//
//	cli, _ := v3.NewCLI[Config]("app", "My app", Config{},
//	    v3.WithMiddleware(spinner.Middleware[Config]("Loading...")),
//	)
package spinner

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"charm.land/lipgloss/v2"
	"golang.org/x/term"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

const defaultInterval = 100 * time.Millisecond

var defaultFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Config configures the visual appearance of the spinner middleware.
type Config struct {
	Title    string
	Writer   io.Writer
	Frames   []string
	Interval time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig(title string) Config {
	return Config{
		Title:    title,
		Writer:   os.Stderr,
		Frames:   defaultFrames,
		Interval: defaultInterval,
	}
}

// Middleware returns a cmdguard middleware that displays a terminal spinner
// while the command handler runs. The spinner writes to stderr and is
// automatically cleared when the command completes.
//
// If stderr is not a terminal, the spinner is silently skipped.
func Middleware[T any](title string) v3.Middleware[T] {
	return MiddlewareWithConfig[T](DefaultConfig(title))
}

// MiddlewareWithConfig returns spinner middleware using the provided configuration.
func MiddlewareWithConfig[T any](cfg Config) v3.Middleware[T] {
	return func(_ context.Context, _ *T, _ v3.CommandInfo, next func() error) error {
		if cfg.Writer == nil || cfg.Interval <= 0 || len(cfg.Frames) == 0 {
			return next()
		}

		if !isTerminal(cfg.Writer) {
			return next()
		}

		s := newTextSpinner(cfg)
		s.Start()

		err := next()

		s.Stop()

		return err
	}
}

type textSpinner struct {
	cfg      Config
	stopCh   chan struct{}
	stopOnce sync.Once
	mu       sync.Mutex
}

func newTextSpinner(cfg Config) *textSpinner {
	return &textSpinner{
		cfg:    cfg,
		stopCh: make(chan struct{}),
	}
}

func (s *textSpinner) Start() {
	go s.run()
}

func (s *textSpinner) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
		s.mu.Lock()
		_, _ = fmt.Fprintf(s.cfg.Writer, "\r\033[K")
		s.mu.Unlock()
	})
}

func (s *textSpinner) run() {
	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF"))
	frameIndex := 0

	frames := s.cfg.Frames
	if len(frames) == 0 {
		frames = defaultFrames
	}

	for {
		select {
		case <-s.stopCh:
			return

		case <-ticker.C:
			frame := frames[frameIndex%len(frames)]
			frameIndex++

			s.mu.Lock()
			_, _ = fmt.Fprintf(
				s.cfg.Writer,
				"\r%s %s",
				style.Render(frame),
				s.cfg.Title,
			)
			s.mu.Unlock()
		}
	}
}

func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}

	return false
}
