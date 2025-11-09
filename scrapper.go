package metrics4gin

import (
	"github.com/VictoriaMetrics/metrics"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ExposeMetrics(c *gin.Context) {
	if h.ExposeRuntimeMetrics {
		metrics.WriteProcessMetrics(c.Writer)
	}
	h.MetricSet.WritePrometheus(c.Writer)
}
