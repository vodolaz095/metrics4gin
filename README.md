Metrics for Gin
==========================

Middleware for git framework with integration with counters/gauges for incoming requests and Prometheus/Victoria Metrics scrapper endpoints.

Example
===========================

See [main.go](example%2Fmain.go) for full working example.

Basically:

```go
package main

import (
	"context"
    "time"
	
	"github.com/gin-gonic/gin"
	"github.com/vodolaz095/metrics4gin"
)

func main() {
	engine := gin.Default()
	metricsHandler := metrics4gin.NewWithDefaultMetricsSet()

	// tune metrics handler by setting parameters
	metricsHandler.Endpoint = metrics4gin.DefaultEndpoint
	metricsHandler.ExposeRuntimeMetrics = true
	metricsHandler.ExtraHeaders["Authorization"] = "Basic dGVzdDp0ZXN0Cg==" // test:test
	metricsHandler.ExtraLabels["job"] = "example_metrics4gin"
	
	// expose prometheus compatible scrapper endpoint
	engine.GET("/metrics", metricsHandler.ExposeMetrics) 

	// periodically push metrics to victoria metrics database
	go metricsHandler.StartPushing(context.Background(), 10*time.Second)

	// start http server
	engine.Run("0.0.0.0:3000")
}


```

Metrics recorded
=================================
There counters are recorded:

- gin_request_total{method=%q,uri=%q,status="%v"} - integer
- gin_request_duration_ms{method=%q,uri=%q,status="%v"} - float
- gin_request_bytes_written{method=%q,uri=%q,status="%v"} - integer
- gin_request_bytes_read{method=%q,uri=%q,status="%v"} - integer
