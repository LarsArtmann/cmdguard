package v2

import (
	"context"
	"errors"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

type mockSpan struct {
	embedded.Span

	name           string
	attributes     []attribute.KeyValue
	recordedErrors []error
	statusCode     codes.Code
	statusDesc     string
	ended          bool
}

func (s *mockSpan) End(_ ...trace.SpanEndOption) {
	s.ended = true
}

func (s *mockSpan) AddEvent(_ string, _ ...trace.EventOption) {}

func (s *mockSpan) AddLink(_ trace.Link) {}

func (s *mockSpan) IsRecording() bool { return true }

func (s *mockSpan) RecordError(err error, _ ...trace.EventOption) {
	s.recordedErrors = append(s.recordedErrors, err)
}

func (s *mockSpan) SetStatus(code codes.Code, description string) {
	s.statusCode = code
	s.statusDesc = description
}

func (s *mockSpan) SetName(_ string) {}

func (s *mockSpan) SetAttributes(attrs ...attribute.KeyValue) {
	s.attributes = append(s.attributes, attrs...)
}

func (s *mockSpan) SpanContext() trace.SpanContext { return trace.SpanContext{} }

func (s *mockSpan) TracerProvider() trace.TracerProvider { return nil }

type mockTracer struct {
	embedded.Tracer

	spans []*mockSpan
}

func (t *mockTracer) Start(
	ctx context.Context,
	spanName string,
	_ ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	span := &mockSpan{name: spanName}
	t.spans = append(t.spans, span)

	return ctx, span
}

func TestTelemetryMiddleware(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	tracer := &mockTracer{}

	mw := TelemetryMiddleware[testConfig](tracer)
	err := mw(
		context.Background(),
		&testConfig{},
		CommandInfo{Name: "deploy", Phase: PhaseRun, HasRunE: true},
		func() error { return nil },
	)

	testutil.AssertNoError(t, err)

	if len(tracer.spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(tracer.spans))
	}

	span := tracer.spans[0]

	if span.name != "deploy run" {
		t.Errorf("expected span name %q, got %q", "deploy run", span.name)
	}

	if !span.ended {
		t.Error("span should be ended")
	}

	attrMap := make(map[string]attribute.Value)
	for _, a := range span.attributes {
		attrMap[string(a.Key)] = a.Value
	}

	if attrMap["command.name"].AsString() != "deploy" {
		t.Error("expected command.name attribute")
	}

	if attrMap["command.phase"].AsString() != "run" {
		t.Error("expected command.phase attribute")
	}

	if !attrMap["command.has_handler"].AsBool() {
		t.Error("expected command.has_handler to be true")
	}
}

func TestTelemetryMiddleware_RecordsError(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	tracer := &mockTracer{}
	expectedErr := errors.New("deploy failed")

	mw := TelemetryMiddleware[testConfig](tracer)
	err := mw(
		context.Background(),
		&testConfig{},
		CommandInfo{Name: "deploy", Phase: PhaseRun},
		func() error { return expectedErr },
	)

	testutil.AssertErrorIs(t, err, expectedErr)

	if len(tracer.spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(tracer.spans))
	}

	span := tracer.spans[0]

	if len(span.recordedErrors) != 1 {
		t.Fatalf("expected 1 recorded error, got %d", len(span.recordedErrors))
	}

	if !errors.Is(span.recordedErrors[0], expectedErr) {
		t.Error("recorded error should match expected error")
	}

	if span.statusCode != codes.Error {
		t.Error("expected error status code")
	}
}

func TestWithTelemetry(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	tracer := &mockTracer{}

	cli, err := NewCLI[testConfig](
		"test", "Test CLI", testConfig{},
		WithTelemetry[testConfig](tracer),
		WithFang[testConfig](false),
	)
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	if len(cli.middleware) != 1 {
		t.Fatalf("expected 1 middleware, got %d", len(cli.middleware))
	}
}

func TestTelemetryMiddleware_UsesFullPath(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	tracer := &mockTracer{}

	mw := TelemetryMiddleware[testConfig](tracer)
	err := mw(
		context.Background(),
		&testConfig{},
		CommandInfo{Name: "migrate", FullPath: "myapp database migrate", Phase: PhaseRun, HasRunE: true},
		func() error { return nil },
	)

	testutil.AssertNoError(t, err)

	if len(tracer.spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(tracer.spans))
	}

	span := tracer.spans[0]

	if span.name != "myapp database migrate run" {
		t.Errorf("expected span name %q, got %q", "myapp database migrate run", span.name)
	}

	attrMap := make(map[string]attribute.Value)
	for _, a := range span.attributes {
		attrMap[string(a.Key)] = a.Value
	}

	if attrMap["command.fullpath"].AsString() != "myapp database migrate" {
		t.Errorf(
			"expected command.fullpath %q, got %q",
			"myapp database migrate",
			attrMap["command.fullpath"].AsString(),
		)
	}
}
