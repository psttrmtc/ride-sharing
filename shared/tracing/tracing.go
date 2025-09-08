package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Config struct {
	ServiceName    string
	Environment    string
	JaegerEndpoint string
}

func InitTracer(cfg Config) (func(context.Context) error, error) {
	// Exporter
	var traceExporter sdktrace.SpanExporter

	// Trace Provider
	traceProvider, err := newTraceProvider(cfg, traceExporter)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(traceProvider)

	//Propagator
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)
	return traceProvider.Shutdown, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(cfg Config, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	res, err := resource.New(context.Background(), resource.WithAttributes(
		semconv.ServiceNameKey.String(cfg.ServiceName),
		semconv.DeploymentEnvironmentKey.String(cfg.Environment),
	),
	)
	if err != nil {
		return nil, fmt.Errorf("faild to create resource: %v", err)
	}
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	return traceProvider, nil
}
