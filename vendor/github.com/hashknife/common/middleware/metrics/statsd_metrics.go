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

// StatsDMetricsMiddleware
type StatsDMetricsMiddleware struct {
	Total    kitmetrics.Counter
	Duration kitmetrics.Histogram
}

// NewStatsDMetricsMiddleware
func NewStatsDMetricsMiddleware(statsd metrics.StatsD, metric string) *StatsDMetricsMiddleware {
	var s StatsDMetricsMiddleware
	s.Total = statsd.NewCounter(metric)
	return &s
}

// AnnotateMux
func (s *StatsDMetricsMiddleware) AnnotateMux(h http.Handler) http.Handler {
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
			s.emitMetrics(begin, nw, tags...)
		}(time.Now())
		h.ServeHTTP(nw, r)
	})
}

// ServeHTTP
func (s *StatsDMetricsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func(begin time.Time) {
		nw := negroni.NewResponseWriter(w)
		next(nw, r)
		s.emitMetrics(begin, nw)
	}(time.Now())
}

// HTTPRouterMiddleware
func (s *StatsDMetricsMiddleware) HTTPRouterMiddleware(next httprouter.Handle, fields ...string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(begin time.Time) {
			nw := negroni.NewResponseWriter(w)
			next(nw, r, p)
			s.emitMetrics(begin, nw, fields...)
		}(time.Now())
	}
}

// Annotate
func (s *StatsDMetricsMiddleware) Annotate(h http.Handler, fields ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nw := negroni.NewResponseWriter(w)
		defer func(begin time.Time) {
			s.emitMetrics(begin, nw, fields...)
		}(time.Now())
		h.ServeHTTP(nw, r)
	})
}

// emitMetrics
func (s *StatsDMetricsMiddleware) emitMetrics(begin time.Time, w negroni.ResponseWriter, extraTags ...string) {
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
	s.Total.
		With(tags...).
		With(extraTags...).
		Add(1)
	s.Duration.
		With(tags...).
		With(extraTags...).
		Observe(float64(time.Since(begin).Nanoseconds()))
}
