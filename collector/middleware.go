package collector

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	initialized     = false
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	responseSize    *prometheus.HistogramVec
	responseCode    *prometheus.CounterVec
	oscillationPeriod = flag.Duration("oscillation-period", 10*time.Minute, "The duration of the rate oscillation period.")//速率振荡周期的持续时间
)

func initializeMetricsMiddleware() {
	if initialized != true {
		requestCounter = createTotalRequestsCollector()
		prometheus.Register(requestCounter)
		requestDuration = createRequestDurationsCollector()
		prometheus.Register(requestDuration)
		responseSize = createRequestSizesCollector()
		prometheus.Register(responseSize)
		responseCode = createResponseCodesCollector()
		prometheus.Register(responseCode)
		initialized = true
	}
}

func InstrumentHandler(next http.Handler) http.Handler {
	if !initialized {
		initializeMetricsMiddleware()
	}
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		instrumentedWriter := &InstrumentedResponseWriter{writer, 0, 200}
		defer func(begun time.Time) {
			start := time.Now()
			oscillationFactor := func() float64 {
				return 2 + math.Sin(math.Sin(2*math.Pi*float64(time.Since(start))/float64(*oscillationPeriod)))
			}
			time.Sleep(time.Duration(75*oscillationFactor()) * time.Millisecond)
			size := float64(instrumentedWriter.Length())
			status := strconv.Itoa(instrumentedWriter.StatusCode())
			method := strings.ToLower(r.Method)
			route := r.URL.Path
			requestCounter.WithLabelValues(method, route, status).Inc()
			requestDuration.WithLabelValues(method, route, status).Observe(float64(time.Since(begun).Seconds() * 1000))
			responseCode.WithLabelValues(status).Inc()
			responseSize.WithLabelValues(method, route, status).Observe(size)
		}(time.Now())
		next.ServeHTTP(instrumentedWriter, r)
	})
}
