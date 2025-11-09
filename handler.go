package metrics4gin

import (
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

// DefaultEndpoint is default endpoint to send metrics too
const DefaultEndpoint = "http://localhost:8428/api/v1/import/prometheus"

// Handler stores logic in it
type Handler struct {
	// MetricSet is
	MetricSet *metrics.Set
	// ExposeRuntimeMetrics toggles exposing default golang runtime metrics at scrapper endpoint
	ExposeRuntimeMetrics bool
	// ExtraLabels are added to all metrics being sent via PUSH mechanism into Prometheus/Victoria Metrics database
	ExtraLabels map[string]string
	// ExtraHeaders are added to HTTP requests being sent
	ExtraHeaders map[string]string
	// Endpoint sets destination for metrics, by default it is  DefaultEndpoint
	Endpoint string
	// Method defines http method used to deliver metrics, default is http.MethodGet
	Method string
}

// NewWithEmptyMetricsSet creates new Handler with empty metrics set
func NewWithEmptyMetricsSet() *Handler {
	return &Handler{
		MetricSet:            metrics.NewSet(),
		ExposeRuntimeMetrics: true,
		ExtraLabels:          make(map[string]string, 0),
		ExtraHeaders:         make(map[string]string, 0),
		Endpoint:             DefaultEndpoint,
		Method:               http.MethodGet,
	}
}

// NewWithDefaultMetricsSet creates new Handler with default golang runtime metrics set
func NewWithDefaultMetricsSet() *Handler {
	return &Handler{
		MetricSet:            metrics.GetDefaultSet(),
		ExposeRuntimeMetrics: true,
		ExtraLabels:          make(map[string]string, 0),
		ExtraHeaders:         make(map[string]string, 0),
		Endpoint:             DefaultEndpoint,
		Method:               http.MethodGet,
	}
}
