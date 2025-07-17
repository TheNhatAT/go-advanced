package tracing

import (
	"context"
	"go-advanced/pkg/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func ExampleLatencyTrace() {
	initCtx := context.Background()
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint("http://localhost:14268/api/traces"),
	))
	if err != nil {
		return
	}
	defer exporter.Shutdown(initCtx)

	res, err := resource.New(initCtx,
		resource.WithAttributes(
			semconv.ServiceName("go-advanced-service"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	tracer := otel.Tracer("example-tracer")

	observability.Prepare()
	for i := 0; i < observability.XTimes; i++ {
		ctx, span := tracer.Start(initCtx, "doOperation")
		err := observability.DoOperationWithCtx(ctx)
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}

	observability.TearDown()
}
