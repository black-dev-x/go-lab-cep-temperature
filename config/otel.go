package config

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type OtelTools struct {
	Tracer trace.Tracer
	Meter  metric.Meter
	Logger *slog.Logger
}

var OTEL *OtelTools = &OtelTools{}
var Context = context.Background()

func SetupOTEL(name string) {

	println("Setting up OpenTelemetry.")
	prop := NewPropagator()
	otel.SetTextMapPropagator(prop)
	tracerProvider, err := NewTracerProvider(name)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tracerProvider)
	meterProvider, err := NewMeterProvider()
	if err != nil {
		panic(err)
	}
	otel.SetMeterProvider(meterProvider)
	loggerProvider, err := NewLoggerProvider()
	if err != nil {
		panic(err)
	}
	global.SetLoggerProvider(loggerProvider)

	OTEL.Tracer = otel.Tracer(name)
	OTEL.Meter = otel.Meter(name)
	OTEL.Logger = otelslog.NewLogger(name)

	println("OpenTelemetry setup completed.")
}

func NewPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func NewTracerProvider(name string) (trace.TracerProvider, error) {
	zipkinEndpoint := "http://zipkin:9411/api/v2/spans"
	zipkinExporter, err := zipkin.New(
		zipkinEndpoint,
	)
	if err != nil {
		return nil, err
	}
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(name), // or "input-service" for the other service
		),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(zipkinExporter),
		sdktrace.WithResource(res),
	)
	return tracerProvider, nil
}

func NewMeterProvider() (metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	)
	return meterProvider, nil
}

func NewLoggerProvider() (*sdklog.LoggerProvider, error) {
	logExporter, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}
