package telemetry

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// HTTPMiddleware returns an HTTP middleware that instruments requests with OpenTelemetry
func HTTPMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "http_request",
			otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		)
	}
}

// StartSpan starts a new span with the given name and attributes
func StartSpan(ctx context.Context, spanName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	tracer := otel.Tracer("goExpertOtel")
	return tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
}

// SetSpanAttributes adds attributes to the current span
func SetSpanAttributes(span trace.Span, attrs ...attribute.KeyValue) {
	span.SetAttributes(attrs...)
}

// RecordError records an error in the current span
func RecordError(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

// SetSpanStatus sets the status of the current span
func SetSpanStatus(span trace.Span, code codes.Code, description string) {
	span.SetStatus(code, description)
}