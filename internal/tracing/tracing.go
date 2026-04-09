package tracing

import (
	"context"

	"urlshortener/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var globalTracer trace.Tracer

func InitTracing(cfg config.TracingConfig) error {
	if !cfg.Enabled {
		return nil
	}

	var exporter *jaeger.Exporter
	var err error

	if cfg.JaegerEndpoint != "" {
		exporter, err = jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(cfg.JaegerEndpoint),
			),
		)
	} else {
		exporter, err = jaeger.New(
			jaeger.WithAgentEndpoint(
				jaeger.WithAgentHost("localhost"),
				jaeger.WithAgentPort("6831"),
			),
		)
	}
	if err != nil {
		return err
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
		),
	)
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.2)),
	)

	otel.SetTracerProvider(tp)

	globalTracer = tp.Tracer(cfg.ServiceName)

	return nil
}

func Tracer() trace.Tracer {
	if globalTracer == nil {
		globalTracer = otel.Tracer("urlshortener")
	}
	return globalTracer
}

func AddTraceAttrsToCtx(ctx context.Context, attrs ...attribute.KeyValue) context.Context {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(attrs...)
	}
	return ctx
}

func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() && err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", err.Error()))
	}
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	t := Tracer()
	ctx, span := t.Start(ctx, name)
	return ctx, span
}

func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}
