package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func Register() {
	if err := prometheus.Register(functionLatency); err != nil {
		log.Fatal(err)
	}
}

const (
	Namespace = "httpserver"
)

func NewTimer() *ExecutionTimer {
	return NewExecutionTimer(functionLatency)
}

var (
	functionLatency = CreateExecutionTimeMetric(Namespace, "Time spent.")
)

func NewExecutionTimer(histogram *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histogram: histogram,
		start:     now,
		last:      now,
	}
}

func (t *ExecutionTimer) ObserveTotal() {
	(*t.histogram).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}

func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"},
	)
}

type ExecutionTimer struct {
	histogram *prometheus.HistogramVec
	start     time.Time
	last      time.Time
}
