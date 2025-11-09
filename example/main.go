package main

import (
	"context"
	_ "embed" // it is ok
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vodolaz095/metrics4gin"
	"golang.org/x/sync/errgroup"
)

//go:embed index.html
var index string

func main() {
	mainCtx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT)

	eg, ctx := errgroup.WithContext(mainCtx)
	defer cancel()

	engine := gin.Default()

	metricsHandler := metrics4gin.NewWithDefaultMetricsSet()

	// set parameters
	metricsHandler.Endpoint = metrics4gin.DefaultEndpoint
	metricsHandler.ExposeRuntimeMetrics = true
	metricsHandler.ExtraHeaders["Authorization"] = "Basic dGVzdDp0ZXN0Cg==" // test:test
	metricsHandler.ExtraLabels["job"] = "example_metrics4gin"

	// inject middleware into gin engine
	metricsHandler.InjectMiddleware(engine)

	// business logic endpoints
	engine.GET("/", func(c *gin.Context) {
		val := metricsHandler.MetricSet.GetOrCreateCounter("example_counter").Get()
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, index, val)
	})
	engine.POST("/", func(c *gin.Context) {
		// pretend we record something into database
		time.Sleep(100 * time.Millisecond)
		c.Redirect(http.StatusFound, "/")
	})

	engine.GET("/incr", func(c *gin.Context) {
		metricsHandler.MetricSet.GetOrCreateCounter("example_counter").Inc()
		c.Redirect(http.StatusFound, "/")
	})
	engine.GET("/decr", func(c *gin.Context) {
		metricsHandler.MetricSet.GetOrCreateCounter("example_counter").Dec()
		c.Redirect(http.StatusFound, "/")
	})
	// expose prometheus/victoria metrics compatible scrapper endpoint
	engine.GET("/metrics", metricsHandler.ExposeMetrics)

	// starting application goroutines
	listener, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		log.Fatalf("error starting gin on 3000 port: %s", err)
	}
	eg.Go(func() error {
		<-ctx.Done()
		return listener.Close()
	})
	eg.Go(func() error {
		return engine.RunListener(listener)
	})

	eg.Go(func() error {
		return metricsHandler.StartPushing(ctx, 10*time.Second)
	})
	err = eg.Wait()
	if err != nil {
		log.Fatalf("error starting application: %s", err)
	}
}
