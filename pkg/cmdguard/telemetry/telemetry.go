// Package telemetry provides OpenTelemetry middleware for cmdguard CLIs.
// It is an optional module — import it only when you need OTel tracing,
// to keep your dependency tree lean.
//
// Usage:
//
//	import (
//	    v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
//	    "github.com/larsartmann/cmdguard/telemetry"
//	)
//
//	tracer := otel.Tracer("myapp")
//	cli, _ := v3.NewCLI[Config]("app", "My app", Config{},
//	    v3.WithMiddleware(telemetry.Middleware[Config](tracer)),
//	)
package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

// Middleware returns a cmdguard middleware that creates an OpenTelemetry span
// for each command execution. The span captures the command name, phase, and
// any error returned by the handler.
//
// Each phase (pre-run, run, post-run) gets a uniquely-named span so traces
// are unambiguous: "deploy pre-run", "deploy run", "deploy post-run".
func Middleware[T any](tracer trace.Tracer) v3.Middleware[T] {
	return func(ctx context.Context, _ *T, info v3.CommandInfo, next func() error) error {
		if tracer == nil {
			return next()
		}

		spanName := info.Name + " " + string(info.Phase)
		if info.FullPath != "" {
			spanName = info.FullPath + " " + string(info.Phase)
		}

		_, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		attrs := []attribute.KeyValue{
			attribute.String("command.name", info.Name),
			attribute.String("command.phase", string(info.Phase)),
			attribute.Bool("command.has_handler", info.HasRunE),
		}

		if info.FullPath != "" {
			attrs = append(attrs, attribute.String("command.fullpath", info.FullPath))
		}

		span.SetAttributes(attrs...)

		err := next()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		return err
	}
}

// WithTelemetry is a convenience wrapper that registers telemetry middleware
// via v3.WithMiddleware.
func WithTelemetry[T any](tracer trace.Tracer) v3.CLIOption[T] {
	return v3.WithMiddleware(Middleware[T](tracer))
}
