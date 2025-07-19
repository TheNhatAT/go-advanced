package cpu

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-advanced/pkg/observability"
	"log"
	"net/http"
	"runtime"
)

func ExampleCPUTimeMetric() {
	// limit for using 2 CPU cores
	runtime.GOMAXPROCS(2)
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	observability.Prepare()

	go func() {
		for i := 0; i < observability.XTimes; i++ {
			err := observability.DoOperation()
			// ...
			_ = err
		}

		observability.TearDown()
	}()

	if err := http.ListenAndServe(":8484", promhttp.HandlerFor(reg, promhttp.HandlerOpts{})); err != nil {
		log.Fatal(err)
	}

	observability.PrintPrometheusMetrics(reg)
}

// use "docker compose up --build" to run this example
