package metrics4gin

import "github.com/VictoriaMetrics/metrics"

const DefaultEndpoint = "http://localhost:8428/api/v1/import/prometheus"

type Handler struct {
	MetricSet            *metrics.Set
	ExposeRuntimeMetrics bool
	ExtraLabels          map[string]string
	ExtraHeaders         map[string]string
	Endpoint             string
}

func NewWithEmptyMetricsSet() *Handler {
	return &Handler{
		MetricSet:            metrics.NewSet(),
		ExposeRuntimeMetrics: true,
		ExtraLabels:          make(map[string]string, 0),
		ExtraHeaders:         make(map[string]string, 0),
		Endpoint:             DefaultEndpoint,
	}
}

func NewWithDefaultMetricsSet() *Handler {
	return &Handler{
		MetricSet:            metrics.GetDefaultSet(),
		ExposeRuntimeMetrics: true,
		ExtraLabels:          make(map[string]string, 0),
		ExtraHeaders:         make(map[string]string, 0),
		Endpoint:             DefaultEndpoint,
	}
}
