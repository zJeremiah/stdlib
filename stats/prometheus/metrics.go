package prometheus

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

func EchoHttpDuration(constLabels map[string]string) *prom.HistogramVec {
	return prom.NewHistogramVec(prom.HistogramOpts{
		Name:        "http_requests_seconds",
		Help:        "Seconds an HTTP request took to complete",
		ConstLabels: constLabels,
	}, []string{"code", "method"})
}

func SqlxDuration(constLabels map[string]string) *prom.HistogramVec {
	return prom.NewHistogramVec(prom.HistogramOpts{
		Name:        "db_requests_seconds",
		Help:        "Seconds a SQL request took to complete",
		ConstLabels: constLabels,
	}, []string{"query"})
}
