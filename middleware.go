package metrics4gin

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
)

// InjectMiddleware injects middleware into gin.Engine that counts requests
func (h *Handler) InjectMiddleware(engine *gin.Engine) {
	engine.Use(func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		startedAt := time.Now()

		requestsTotalCounter := h.MetricSet.GetOrCreateCounter(fmt.Sprintf("gin_request_total{method=%q,uri=%q,status=\"%v\"}",
			c.Request.Method, c.FullPath(), c.Writer.Status()))
		requestsTotalCounter.Add(1)

		bytesWrittenCounter := h.MetricSet.GetOrCreateCounter(fmt.Sprintf("gin_request_bytes_written{method=%q,uri=%q,status=\"%v\"}",
			c.Request.Method, c.FullPath(), c.Writer.Status()))
		bytesReadCounter := h.MetricSet.GetOrCreateCounter(fmt.Sprintf("gin_request_bytes_read{method=%q,uri=%q,status=\"%v\"}",
			c.Request.Method, c.FullPath(), c.Writer.Status()))

		span.SetAttributes(attribute.Int64("metrics.gin_request_total", int64(requestsTotalCounter.Get())))

		c.Next()

		durationCounter := h.MetricSet.GetOrCreateFloatCounter(fmt.Sprintf("gin_request_duration_ms{method=%q,uri=%q,status=\"%v\"}",
			c.Request.Method, c.FullPath(), c.Writer.Status()))

		durationCounter.Add(float64(time.Since(startedAt).Microseconds() / 1000))
		span.SetAttributes(attribute.Float64("metrics.gin_request_duration_ms", durationCounter.Get()))

		if c.Writer.Written() {
			bytesWrittenCounter.Add(c.Writer.Size())
			span.SetAttributes(attribute.Int64("metrics.gin_request_bytes_written", int64(bytesWrittenCounter.Get())))
		}
		if c.Request.ContentLength > 0 {
			bytesReadCounter.Add(int(c.Request.ContentLength))
			span.SetAttributes(attribute.Int64("metrics.gin_request_bytes_read", int64(bytesReadCounter.Get())))
		}
		span.AddEvent("metrics recorded", trace.WithAttributes(
			attribute.Int64("metrics.gin_request_total", int64(requestsTotalCounter.Get())),
			attribute.Float64("metrics.gin_request_duration_ms", durationCounter.Get()),
			attribute.Int64("metrics.gin_request_bytes_written", int64(bytesWrittenCounter.Get())),
			attribute.Int64("metrics.gin_request_bytes_read", int64(bytesReadCounter.Get())),
		))
	})
}
