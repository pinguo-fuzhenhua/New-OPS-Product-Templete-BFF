package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Level string
}

// NewFactory
// @TODO move thee config to yaml
func NewFactory() *Factory {
	cfg := &Config{
		Level: "debug",
	}
	return &Factory{
		cfg:          cfg,
		noopProvider: trace.NewNoopTracerProvider(),
	}
}

type Factory struct {
	cfg          *Config
	noopProvider trace.TracerProvider
}

func (s *Factory) StartNewTracer(ctx context.Context, name string) (context.Context, trace.Tracer, trace.Span) {
	var t trace.Tracer
	if s.cfg.Level != "debug" {
		t = s.noopProvider.Tracer(name)
	} else {
		t = otel.Tracer(name)
	}
	ctx2, span := t.Start(ctx, name)
	return ctx2, t, span
}
