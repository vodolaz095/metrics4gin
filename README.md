Metrics for Gin
==========================

[![PkgGoDev](https://pkg.go.dev/badge/github.com/vodolaz095/metrics4gin)](https://pkg.go.dev/github.com/vodolaz095/metrics4gin?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/vodolaz095/metrics4gin)](https://goreportcard.com/report/github.com/vodolaz095/metrics4gin)


Metric Middleware for [gin](https://github.com/gin-gonic/gin/) framework with request counters
for Prometheus/Victoria Metrics time-series databases.
Push and pull (via scrapper endpoint) mechanisms are supported

Metrics recorded
=================================
These counters are recorded, tagged with method, uri and status code:

- `gin_request_total{method=%q,uri=%q,status="%v"}` - integer - total number of HTTP requests served.
- `gin_request_duration_ms{method=%q,uri=%q,status="%v"}` - float - counter for duration of requests in milliseconds.
- `gin_request_bytes_written{method=%q,uri=%q,status="%v"}` - integer - counter for bytes written as responses.
- `gin_request_bytes_read{method=%q,uri=%q,status="%v"}` - integer - counter for bytes read as requests body.

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
	// enable exposing default golang runtime metrics on scrapper endpoint
	metricsHandler.ExposeRuntimeMetrics = true
	// headers to be send with every request for pushing data in Victoria Metrics/Prometheus
	metricsHandler.ExtraHeaders["Authorization"] = "Basic dGVzdDp0ZXN0Cg==" // test:test
	// extra labels to be added to all metrics send to Victoria Metrics/Prometheus
	metricsHandler.ExtraLabels["job"] = "example_metrics4gin"
	metricsHandler.ExtraLabels["instance"] = "server.local"
	
	// expose prometheus compatible scrapper endpoint
	engine.GET("/metrics", metricsHandler.ExposeMetrics) 

	// periodically push metrics to victoria metrics database
	go metricsHandler.StartPushing(context.Background(), 10*time.Second)

	// start http server
	engine.Run("0.0.0.0:3000")
}


```

MIT License
=============================

Copyright (c) 2025 Остроумов Анатолий

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
