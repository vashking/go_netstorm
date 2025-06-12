package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	RequestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "flooder_requests_total",
		Help: "Total number of requests sent.",
	})
	SuccessTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "flooder_success_total",
		Help: "Total number of successful requests.",
	})
	FailTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "flooder_fail_total",
		Help: "Total number of failed requests.",
	})
	LatencySum = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "flooder_latency_sum_seconds",
		Help: "Total latency in seconds.",
	})
)

func InitPrometheus() {
	prometheus.MustRegister(RequestsTotal, SuccessTotal, FailTotal, LatencySum)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)
}

func IncRequests()  { RequestsTotal.Inc() }
func IncSuccess()   { SuccessTotal.Inc() }
func IncFail()      { FailTotal.Inc() }
func AddLatency(s float64) { LatencySum.Add(s) } 