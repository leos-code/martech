package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	RTAProcessErrorCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "ug",
			Subsystem:   "rta",
			Name:        "process_err_count",
			Help:        "process err count",
		})
)

func init()  {
	prometheus.MustRegister(RTAProcessErrorCount)
}
