package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/codegangsta/negroni"
	kitmetrics "github.com/go-kit/kit/metrics"
	"github.com/gorilla/mux"
	"github.com/hashknife/common/metrics"
	"github.com/julienschmidt/httprouter"
)

// DatadogMetricsMiddleware
type DatadogMetricsMiddleware struct {
	Total    kitmetrics.Counter
	Duration kitmetrics.Histogram
}

// NewDatadogMetricsMiddleware
func NewDatadogMetricsMiddleware(datadog metrics.Datadog, metric string) *DatadogMetricsMiddleware {
	var m DatadogMetricsMiddleware
	m.Total = datadog.NewCounter(metric)
	m.Duration = datadog.NewHistogram(metric)
	return &m
}

// AnnotateMux
func (m *DatadogMetricsMiddleware) AnnotateMux(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tags := []string{}
		for k, v := range mux.Vars(r) {
			// Whitelist fields we want tag on metrics
			switch k {
			case "protection":
				tags = append(tags, k, v)
			case "accountID":
				tags = append(tags, k, v)
			}
		}
		nw := negroni.NewResponseWriter(w)
		defer func(begin time.Time) {
			m.emitMetrics(begin, nw, tags...)
		}(time.Now())
		h.ServeHTTP(nw, r)
	})
}

// ServeHTTP
func (m *DatadogMetricsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func(begin time.Time) {
		nw := negroni.NewResponseWriter(w)
		next(nw, r)
		m.emitMetrics(begin, nw)
	}(time.Now())
}

// HTTPRouterMiddleware
func (m *DatadogMetricsMiddleware) HTTPRouterMiddleware(next httprouter.Handle, fields ...string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(begin time.Time) {
			nw := negroni.NewResponseWriter(w)
			next(nw, r, p)
			m.emitMetrics(begin, nw, fields...)
		}(time.Now())
	}
}

// Annotate
func (m *DatadogMetricsMiddleware) Annotate(h http.Handler, fields ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nw := negroni.NewResponseWriter(w)
		defer func(begin time.Time) {
			m.emitMetrics(begin, nw, fields...)
		}(time.Now())
		h.ServeHTTP(nw, r)
	})
}

// emitMetrics
func (m *DatadogMetricsMiddleware) emitMetrics(begin time.Time, w negroni.ResponseWriter, extraTags ...string) {
	var tags = []string{
		"http_status", strconv.FormatInt(int64(w.Status()), 10),
		"error", strconv.FormatBool(w.Status() >= 400),
	}
	var errorType = "none"
	if w.Status() >= 400 {
		errorType = "4xx"
	}
	if w.Status() >= 500 {
		errorType = "5xx"
	}
	if errorType != "none" {
		tags = append(tags, "error_type", errorType)
	}
	m.Total.
		With(tags...).
		With(extraTags...).
		Add(1)
	m.Duration.
		With(tags...).
		With(extraTags...).
		Observe(float64(time.Since(begin).Nanoseconds()))
}
