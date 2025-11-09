package metrics4gin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

// StartPushing starts pushing handler metrics into Victoria/Prometheus
func (h *Handler) StartPushing(ctx context.Context, interval time.Duration) error {
	labels := make([]string, len(h.ExtraLabels))
	headers := make([]string, len(h.ExtraHeaders))
	var i = 0
	for k := range h.ExtraLabels {
		labels[i] = fmt.Sprintf("%s=%q", k, h.ExtraLabels[k])
		i++
	}
	i = 0
	for k := range h.ExtraHeaders {
		headers[i] = fmt.Sprintf("%s : %s", k, h.ExtraHeaders)
		i++
	}
	opts := metrics.PushOptions{
		ExtraLabels:        strings.Join(labels, ","),
		Headers:            headers,
		DisableCompression: false,
		Method:             h.Method,
	}
	return h.MetricSet.InitPushWithOptions(ctx, h.Endpoint, interval, &opts)
}

// StartPushingDefaultMetrics starts pushing default (including runtime with 'process_*' and `go_*` if Handler.ExposeRuntimeMetrics enabled)
// metrics into Victoria/Prometheus
func (h *Handler) StartPushingDefaultMetrics(ctx context.Context, interval time.Duration) error {
	labels := make([]string, len(h.ExtraLabels))
	headers := make([]string, len(h.ExtraHeaders))
	var i = 0
	for k := range h.ExtraLabels {
		labels[i] = fmt.Sprintf("%s=%q", k, h.ExtraLabels[k])
		i++
	}
	i = 0
	for k := range h.ExtraHeaders {
		headers[i] = fmt.Sprintf("%s : %s", k, h.ExtraHeaders)
		i++
	}
	opts := metrics.PushOptions{
		ExtraLabels:        strings.Join(labels, ","),
		Headers:            headers,
		DisableCompression: false,
		Method:             h.Method,
	}
	return metrics.InitPushWithOptions(ctx, h.Endpoint, interval, h.ExposeRuntimeMetrics, &opts)
}
