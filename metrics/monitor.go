package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	histogram *prometheus.HistogramVec
}

func New() *Metrics {
	m := &Metrics{}
	m.histogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "request_duration",
		Help: "counts calls to api",
	}, []string{"type", "status", "isError", "errorMessage", "method", "addr"})
	return m
}

func (m *Metrics) CountApiCall(typeReq string, status string, isError string, errorMessage string, method string, addr string, val float64) {
	m.histogram.WithLabelValues(typeReq, status, isError, errorMessage, method, addr).Observe(val)
}
