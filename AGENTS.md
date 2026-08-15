# metrics4gin - Agent Context

This is a Go library that generates metrics middleware for the [gin](https://github.com/gin-gonic/gin/) web framework. It records HTTP request metrics and supports Prometheus/Victoria Metrics time-series databases.

## Project Overview

- **Purpose**: A metrics middleware for Gin frameworks that tracks HTTP requests with request counts, durations, bytes read/written
- **Metrics recorded** (tagged with method, URI, and status code):
  - `gin_request_total{method,%q,uri=%q,status='%v'}` - total number of HTTP requests
  - `gin_request_duration_ms{method,%q,uri=%q,status='%v'}` - request duration in milliseconds
  - `gin_request_bytes_written{method,%q,uri=%q,status='%v'}` - bytes written in responses
  - `gin_request_bytes_read{method,%q,uri=%q,status='%v'}` - bytes read in request bodies
- **Push and pull mechanisms**: Active push to Victoria Metrics/Prometheus and a scrapper endpoint
- **License**: MIT

## Building and Running

### Build and run the example server
```bash
make start        # Runs the example main.go
make build        # Builds the library
make test         # Runs tests
make vet          # Runs go vet for type safety checks
make lint         # Runs golangci-lint
go build          # Manual build
go test           # Manual test run
```

### Build flags
```bash
go build -o binary
go build -ldflags "-s -w" # Production build without debug symbols
```

### Test running
```bash
go test ./...     # Run all tests
go test -v ./...  # Run all tests with verbose output
```

## Architecture

The codebase follows a simple, focused architecture:

### Core components
| File | Purpose |
|------|---------|
| `handler.go` | Central Handler struct containing metrics configuration and metrics handler state; contains constructor functions for Handler (`NewWithDefaultMetricsSet`, `NewWithEmptyMetricsSet`) and the `InjectMiddleware` method |
| `middleware.go` | Contains `InjectMiddleware` that adds metrics collection to gin routes using gin's `engine.Use()` method |
| `pusher.go` | Handles metric push via HTTP POST to Victoria Metrics/Prometheus endpoints and supports runtime metrics exposure at `/metrics` |
| `scrapper.go` | Handles `/metrics` endpoint for scrapper queries to pull metrics |
| `example/main.go` | Example application demonstrating full integration with custom configuration |

### Control flow
1. User creates `metrics4gin.NewWithDefaultMetricsSet()` (or `NewWithEmptyMetricsSet()`).
2. User configures `Handler` fields: `Endpoint`, `Method`, `ExtraLabels`, `ExtraHeaders`, `ExposeRuntimeMetrics`.
3. User calls `handler.InjectMiddleware(engine)` to inject metrics middleware into a Gin engine.
4. User optionally exposes `/metrics` endpoint via `handler.ExposeMetrics`.
5. User optionally pushes metrics periodically via `handler.StartPushing(ctx, interval)`.

### Data flow
- `InjectMiddleware`: Wraps each request with middleware, records request counts, durations, bytes written/read using VictoriaMetrics `metrics.Set`, and creates OTel spans for observability.
- `Pusher`: Periodically sends POST to configured Prometheus/Victoria Metrics endpoint with `Content-Type: application/prometheus`.
- `Scrapper`: Exposes `/metrics` endpoint for Prometheus-compatible pulling.
- `Metrics`: VictoriaMetrics library (`github.com/VictoriaMetrics/metrics`) provides counters.

## Handler Configuration

The `Handler` struct is a flexible container for metrics collection configuration:

```go
type Handler struct {
    MetricSet            *metrics.Set   // VictoriaMetrics.Set
    ExposeRuntimeMetrics bool            // Toggle exposing Go runtime metrics
    ExtraLabels          map[string]string // Label added to all metrics
    ExtraHeaders         map[string]string // HTTP headers added to push requests
    Endpoint             string          // Push destination URL (default: "http://localhost:8428/api/v1/import/prometheus")
    Method               string          // HTTP method for push (default: "GET")
}
```

### Key patterns
- `Endpoint` and `Method` control where and how metrics are pushed
- `ExtraLabels` must be a non-nil map; defaults to `{}` but can be modified by the user
- `ExposeRuntimeMetrics` defaults to `true`, enabling runtime metrics in `/metrics` endpoint
- OTel spans are automatically created for metrics, recording attributes with `span.SetAttributes(...)`

### Common configuration snippet from example:
```go
handler := metrics4gin.NewWithDefaultMetricsSet()
handler.Endpoint = metrics4gin.DefaultEndpoint // Can change this
handler.ExposeRuntimeMetrics = true
handler.ExtraHeaders["Authorization"] = "Basic dGVzdDp0ZXN0Cg==" // base64 encoded credentials
handler.ExtraLabels["job"] = "example_metrics4gin"
handler.Inject Middleware(engine)
engine.GET("/metrics", handler.ExposeMetrics)
go handler.StartPushing(context.Background(), 10*time.Second)
```

## Code Style Patterns

### Go version
- Go 1.24.9 (as declared in `go.mod`)
- Uses go modules (not GOPATH mode)

### Package organization
- All files are under `package metrics4gin`
- No sub-packages currently

### Conventions
- **Function naming**: `NewWith*` pattern for constructors (e.g., `NewWithDefaultMetricsSet`)
- **Type comments**: Inline comments for struct fields use `// FieldName ... description` format
- **Constant definitions**: `const` for global values like `DefaultEndpoint`
- **Struct comments**: Full `package metrics4gin` comment block with imports at top

## Integration Pattern

```go
func main() {
    engine := gin.Default()
    metricsHandler := metrics4gin.NewWithDefaultMetricsSet()
    // tune metricsHandler by setting parameters (Endpoint, ExposeRuntimeMetrics, ExtraLabels, etc.)
    metricsHandler.InjectMiddleware(engine)
    // optional custom endpoints: engine.GET("/metrics", handler.ExposeMetrics)
    // optional periodic pushing: go handler.StartPushing(ctx, interval)
    engine.Run("0.0.0.0:3000")
}
```

## Dependencies

Key dependencies:
- `github.com/gin-gonic/gin` v1.11.0 - Gin web framework
- `github.com/VictoriaMetrics/metrics` v1.40.2 - Metrics collection library
- `go.opentelemetry.io/otel` v1.38.0 and v1.38.0 (trace) - OpenTelemetry for distributed tracing
- `golang.org/x/sync` v0.17.0 - Context synchronization utilities

## Common Gotchas

1. **Metrics are always tagged**: Every metric automatically includes method, URI, and status code tags.
2. **Push URL must be accessible**: Ensure target VictoriaMetrics/Prometheus endpoint allows remote connections and CORS if needed.
3. **Headers must be base64 encoded**: `ExtraHeaders` uses HTTP Basic Authentication format.
4. **Runtime metrics are exposed at `/metrics`**: Default exposed at `http://localhost:3000/metrics` if `ExposeRuntimeMetrics=true`
5. **Push method default is GET**: `Method` defaults to `http.MethodGet`, but VictoriaMetrics typically uses POST
6. **ExtraHeaders/ExtraLabels require non-nil maps**: Using empty `make(map[string]string, 0)` is acceptable, but zero maps may cause failures

## File Structure

```
.
├── example/           # Demo application with configuration and running example
│   ├── index.html
├── handler.go         # Handler struct with metrics collection logic
├── middleware.go      # Gin middleware injection using engine.Use()
├── pusher.go          # Push to remote metrics endpoint and runtime metrics exposure
├── scrapper.go        # /metrics endpoint for Prometheus scrapper
├── Makefile           # Project build/run/test automation
├── go.mod             # Go module definition
```

## External Documentation

- [Gin Framework](https://github.com/gin-gonic/gin) - Web framework with clean syntax, routing, middleware support
- [VictoriaMetrics](https://github.com/VictoriaMetrics/VictoriaMetrics) - Time-series database
- [OpenTelemetry](https://opentelemetry.io/) - Observability tools (tracing attributes recorded per metric)

## Testing

The example in `example/main.go` serves as integration test demonstrating:
- Basic middleware injection
- Configuration customization
- Push functionality
- Runtime metrics exposure

## Deployment Notes

1. **Metrics endpoint**: Exposed at `/metrics` by default, follows Prometheus exposition format.
2. **Push endpoint**: POST to destination with `application/prometheus` content-type.
3. **Runtime metrics**: Enabled by default, includes Go runtime metrics (allocations, GC, etc.).
4. **Distributed tracing**: OTel spans enabled, capturing metrics as attributes.
