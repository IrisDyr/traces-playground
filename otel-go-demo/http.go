package main

import (
  "context"
  "fmt"
  "io"
  "net/http"
  "time"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/codes"
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
  "go.opentelemetry.io/otel/sdk/resource"
  sdktrace "go.opentelemetry.io/otel/sdk/trace"
  "go.opentelemetry.io/otel/trace"
  semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func main() {
  ctx := context.Background()
  exporter, err := otlptracehttp.New(ctx,
    otlptracehttp.WithEndpoint("YOUR_COLLECTOR:4318"), // if using ingress, you might not need the port
    otlptracehttp.WithURLPath("/v1/traces"), // this is the default path, use your own path if your collector is expecting smth different
    otlptracehttp.WithInsecure(), // remove if using TLS
  )
  if err != nil {
    panic(fmt.Errorf("failed to create OTLP exporter: %w", err))
  }

  tp := sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exporter),
    sdktrace.WithResource(resource.NewWithAttributes(
      semconv.SchemaURL,
      semconv.ServiceNameKey.String("go-otel-demo"),
      semconv.DeploymentEnvironmentKey.String("dev"),
    )),
  )
  defer tp.Shutdown(ctx)

  otel.SetTracerProvider(tp)
  tracer := otel.Tracer("go-otel-tracer")

  ctx, mainSpan := tracer.Start(ctx, "main-operation",
    trace.WithAttributes(attribute.String("startup", "true")))
  defer mainSpan.End()
  mainSpan.AddEvent("begin main operation")

  if err := preCheck(ctx, tracer); err != nil {
    mainSpan.RecordError(err); mainSpan.SetStatus(codes.Error, err.Error()); return
  }
  if err := doWork(ctx, tracer); err != nil {
    mainSpan.RecordError(err); mainSpan.SetStatus(codes.Error, err.Error()); return
  }
  if err := postProcess(ctx, tracer); err != nil {
    mainSpan.RecordError(err); mainSpan.SetStatus(codes.Error, err.Error()); return
  }
  mainSpan.SetStatus(codes.Ok, "all steps completed")
  fmt.Println("Trace done.")
}

func preCheck(ctx context.Context, tracer trace.Tracer) error {
  ctx, span := tracer.Start(ctx, "pre-check",
    trace.WithAttributes(attribute.String("step", "precheck")))
  defer span.End()
  span.AddEvent("checking prerequisites")
  time.Sleep(20 * time.Millisecond)
  return nil
}

func doWork(ctx context.Context, tracer trace.Tracer) error {
  ctx, span := tracer.Start(ctx, "doMagic",
    trace.WithAttributes(attribute.String("operation", "db+api")))
  defer span.End()
  span.AddEvent("start work sequence")
  time.Sleep(30 * time.Millisecond)

  if err := cacheLookup(ctx, tracer); err != nil {
    span.AddEvent("cache miss")
  }
  if err := databaseCall(ctx, tracer); err != nil {
    span.RecordError(err); span.SetStatus(codes.Error, err.Error()); return err
  }
  if err := callExternalAPI(ctx, tracer); err != nil {
    span.RecordError(err); span.SetStatus(codes.Error, err.Error()); return err
  }
  span.SetStatus(codes.Ok, "Magic ok")
  return nil
}

func cacheLookup(ctx context.Context, tracer trace.Tracer) error {
  ctx, span := tracer.Start(ctx, "cache lookup",
    trace.WithAttributes(attribute.String("cache", "redis")))
  defer span.End()
  span.AddEvent("checking cache")
  time.Sleep(15 * time.Millisecond)
  return fmt.Errorf("cache miss")
}

func databaseCall(ctx context.Context, tracer trace.Tracer) error {
  ctx, span := tracer.Start(ctx, "database query",
    trace.WithAttributes(
      semconv.DBSystemKey.String("postgresql"),
      semconv.DBStatementKey.String("SELECT * FROM users WHERE id=123"),
    ))
  defer span.End()
  span.AddEvent("sending fancy db query")
  time.Sleep(40 * time.Millisecond)
  span.AddEvent("db query returned")
  return nil
}

func callExternalAPI(ctx context.Context, tracer trace.Tracer) error {
  ctx, span := tracer.Start(ctx, "HTTP GET example.com",
    trace.WithAttributes(
      semconv.HTTPMethodKey.String("GET"),
      semconv.HTTPTargetKey.String("/"),
      semconv.HTTPSchemeKey.String("https"),
      semconv.HTTPURLKey.String("https://www.example.com/"),
    ))
  defer span.End()

  span.AddEvent("sending HTTP request")
  req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.example.com/", nil)
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    span.RecordError(err); span.SetStatus(codes.Error, err.Error()); return err
  }
  defer resp.Body.Close()
  span.SetAttributes(semconv.HTTPStatusCodeKey.Int(resp.StatusCode))
  if resp.StatusCode >= 400 {
    span.SetStatus(codes.Error, http.StatusText(resp.StatusCode))
  } else {
    span.SetStatus(codes.Ok, "")
  }
  span.AddEvent("reading HTTP response")
  _, _ = io.ReadAll(resp.Body)
  span.AddEvent("response read")
  return nil
}

func postProcess(ctx context.Context, tracer trace.Tracer) error {
  ctx, span := tracer.Start(ctx, "post-process",
    trace.WithAttributes(attribute.String("step", "postprocess")))
  defer span.End()
  span.AddEvent("starting post-processing tasks")
  time.Sleep(25 * time.Millisecond)
  span.AddEvent("cleaning temporary data")
  time.Sleep(15 * time.Millisecond)
  span.AddEvent("finalizing operation")
  time.Sleep(10 * time.Millisecond)
  span.SetStatus(codes.Ok, "post-process done")
  return nil
}
