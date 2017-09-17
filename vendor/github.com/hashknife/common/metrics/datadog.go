package metrics

import (
	"io/ioutil"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/dogstatsd"
	"github.com/go-kit/kit/util/conn"
)

type Datadog struct {
	Dogstatsd  *dogstatsd.Dogstatsd
	Interval   *time.Ticker
	SampleRate float64
}

// NewDatadog
func NewDatadog(prefix, address string, interval time.Duration) Datadog {
	var d = Datadog{
		Dogstatsd:  dogstatsd.New(prefix+".", log.NewNopLogger()),
		Interval:   time.NewTicker(interval),
		SampleRate: 1,
	}
	var sink = ioutil.Discard
	if address != "" {
		sink = conn.NewDefaultManager("udp", address, log.NewNopLogger())
	}
	go d.Dogstatsd.WriteLoop(d.Interval.C, sink)
	return d
}

// NewCounter
func (d Datadog) NewCounter(name string) *dogstatsd.Counter {
	return d.Dogstatsd.NewCounter(name+".count", d.SampleRate)
}

// NewGauge
func (d Datadog) NewGauge(name string) *dogstatsd.Gauge {
	return d.Dogstatsd.NewGauge(name + ".gauge")
}

// NewHistogram
func (d Datadog) NewHistogram(name string) *dogstatsd.Histogram {
	// Histograms have an unspecified unit
	return d.Dogstatsd.NewHistogram(name+".duration", d.SampleRate)
}

// NewTiming
func (d Datadog) NewTiming(name string) *dogstatsd.Timing {
	// Timings are interpreted as a value in milliseconds
	return d.Dogstatsd.NewTiming(name+".duration", d.SampleRate)
}
