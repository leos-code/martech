package metricutil

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Observed
type Observed interface {
	labels() []string
}

// CalcTimeUsedMicro 计算使用时间，单位微秒
func CalcTimeUsedMicro(startTime time.Time) float64 {
	return float64(time.Since(startTime) / time.Microsecond)
}

// CalcTimeUsedMilli 计算使用时间，单位毫秒
func CalcTimeUsedMilli(startTime time.Time) float64 {
	return float64(time.Since(startTime) / time.Millisecond)
}

// ServeMetrics 启动metrics服务
func ServeMetrics(metricsServerAddress string) error {
	go func() {
		muxProm := http.NewServeMux()
		muxProm.Handle("/metrics", promhttp.Handler())
		var err error
		if err = http.ListenAndServe(metricsServerAddress, muxProm); err != nil {
			glog.Errorf("failed to listen prometheus, err: %v", err)
		}
	}()

	return nil
}

// CollectActionMetrics 收集错误、延时信息，vec
func CollectActionMetrics(errCounter *prometheus.CounterVec,
	cost *prometheus.HistogramVec, name string, startTime time.Time, err error) {

	cost.WithLabelValues(name).Observe(CalcTimeUsedMilli(startTime))
	if err != nil {
		errCounter.WithLabelValues(name).Add(1)
	}
}

// CollectMetrics 收集错误、延时信息
func CollectMetrics(errCounter prometheus.Counter,
	cost prometheus.Observer, startTime time.Time, err error) {

	cost.Observe(CalcTimeUsedMilli(startTime))
	if err != nil {
		errCounter.Add(1)
	}
}
