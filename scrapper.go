package metrics4gin

import (
	"github.com/VictoriaMetrics/metrics"
	"github.com/gin-gonic/gin"
)

// DefaultMetricsEndpointRelativePath is usual path to expose metrics for Prometheus scrappers
const DefaultMetricsEndpointRelativePath = "/metrics"

// ExposeMetrics expose metrics endpoint
func (h *Handler) ExposeMetrics(c *gin.Context) {
	if h.ExposeRuntimeMetrics {
		metrics.WriteProcessMetrics(c.Writer)
	}
	h.MetricSet.WritePrometheus(c.Writer)
}
