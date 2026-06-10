package v2

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TelemetryMiddleware returns a middleware that creates an OpenTelemetry span
// for each command execution. The span captures the command name, phase, and
// any error returned by the handler.
//
// Each phase (pre-run, run, post-run) gets a uniquely-named span so traces
// are unambiguous: "deploy pre-run", "deploy run", "deploy post-run".
//
// Usage:
//
//	tracer := otel.Tracer("myapp")
//	cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
//	    v2.WithMiddleware(v2.TelemetryMiddleware[Config](tracer)),
//	)
func TelemetryMiddleware[T any](tracer trace.Tracer) Middleware[T] {
	return func(ctx context.Context, _ *T, info CommandInfo, next func() error) error {
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

// WithTelemetry adds OpenTelemetry span tracking to every command.
// It is a convenience wrapper around WithMiddleware(TelemetryMiddleware(tracer)).
//
// Usage:
//
//	tracer := otel.Tracer("myapp")
//	cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
//	    v2.WithTelemetry[Config](tracer),
//	)
func WithTelemetry[T any](tracer trace.Tracer) CLIOption[T] {
	return WithMiddleware[T](TelemetryMiddleware[T](tracer))
}
