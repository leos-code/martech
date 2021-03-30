package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	ProcessErrCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "ug",
		Subsystem: "indexer",
		Name:      "process_err_count",
		Help:      "process err count",
	})

	ExceedMaxIntervalCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "ug",
		Subsystem: "indexer",
		Name:      "exceed_max_interval_count",
		Help:      "exceed max interval count",
	})

	SuccessCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "ug",
		Subsystem: "indexer",
		Name:      "success_count",
		Help:      "success count",
	})
)

func init() {
	prometheus.MustRegister(ProcessErrCount)
	prometheus.MustRegister(ExceedMaxIntervalCount)
	prometheus.MustRegister(SuccessCount)
}
