package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-advanced/pkg/observability"
	"log"
	"net/http"
	"time"
)

func ExampleLatencyMetric() {
	reg := prometheus.NewRegistry()
	latencySeconds := promauto.With(reg).
		NewHistogramVec(prometheus.HistogramOpts{
			Name:    "operation_duration_seconds",
			Help:    "Tracks the latency of operations in seconds.",
			Buckets: []float64{0.001, 0.01, 0.1, 1, 10, 100},
		}, []string{"error_type"})

	observability.Prepare()

	go func() {
		for i := 0; i < observability.XTimes; i++ {
			now := time.Now()
			err := observability.DoOperation() // Operation we want to measure and potentially optimize...
			elapsed := time.Since(now)

			// Prometheus metric.
			latencySeconds.WithLabelValues(observability.ErrorType(err)).
				Observe(elapsed.Seconds())

			if err != nil { /* Handle error... */
			}

			time.Sleep(1 * time.Second)
		}
	}()

	if err := http.ListenAndServe(
		":8484",
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
	); err != nil {
		log.Fatal(err)
	}

	observability.TearDown()

	observability.PrintPrometheusMetrics(reg)
}
