package v2

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"charm.land/lipgloss/v2"
	"golang.org/x/term"
)

// defaultSpinnerInterval is the default delay between spinner frames.
const defaultSpinnerInterval = 100 * time.Millisecond

// spinnerFrames returns the default braille pattern frames for the spinner.
func spinnerFrames() []string {
	return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
}

// SpinnerConfig configures the visual appearance of the spinner middleware.
type SpinnerConfig struct {
	Title    string
	Writer   io.Writer
	Frames   []string
	Interval time.Duration
}

// DefaultSpinnerConfig returns a SpinnerConfig with sensible defaults.
func DefaultSpinnerConfig(title string) SpinnerConfig {
	return SpinnerConfig{
		Title:    title,
		Writer:   os.Stderr,
		Frames:   spinnerFrames(),
		Interval: defaultSpinnerInterval,
	}
}

// SpinnerMiddleware returns a middleware that displays a terminal spinner
// while the command handler runs. The spinner writes to stderr and is
// automatically cleared when the command completes.
//
// If stderr is not a terminal, the spinner is silently skipped to avoid
// polluting logs or redirected output.
//
// Usage:
//
//	cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
//	    v2.WithMiddleware(v2.SpinnerMiddleware[Config]("Loading...")),
//	)
func SpinnerMiddleware[T any](title string) Middleware[T] {
	return func(_ context.Context, _ *T, _ CommandInfo, next func() error) error {
		cfg := DefaultSpinnerConfig(title)
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

// textSpinner is a simple goroutine-based spinner for terminals.
type textSpinner struct {
	cfg      SpinnerConfig
	stopCh   chan struct{}
	stopOnce sync.Once
	mu       sync.Mutex
}

func newTextSpinner(cfg SpinnerConfig) *textSpinner {
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
		// Clear the spinner line.
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
		frames = spinnerFrames()
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
