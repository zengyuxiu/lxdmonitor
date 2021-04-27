package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	Netstats = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "netstats",
		Help: "record net status of instance",
	},
		[]string{"instance", "interface", "type"})
	Source = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "source",
		Help: "record cpu&mem status of instance",
	}, []string{"instance", "sourcetype"})
)

func prometheus_srv() {
	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":2112", nil)
}
