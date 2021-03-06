package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	farmCounter      *prometheus.CounterVec
}

func New() *Metrics {
	m := &Metrics{}
	m.farmCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "farm_api_api_farm_requests",
		Help: "counts calls to api farm",
	}, []string{})
	return m
}

func (m *Metrics) CountApiFarm() {
	m.farmCounter.WithLabelValues().Inc()
}