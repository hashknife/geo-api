package metrics

import (
	"io/ioutil"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/statsd"
	"github.com/go-kit/kit/util/conn"
)

// StatsD
type StatsD struct {
	StatsD     *statsd.Statsd
	Interval   *time.Ticker
	SampleRate float64
}

// NewStatsD
func NewStatsD(prefix, address string, interval time.Duration) StatsD {
	var s = StatsD{
		StatsD:     statsd.New(prefix+".", log.NewNopLogger()),
		Interval:   time.NewTicker(interval),
		SampleRate: 1,
	}
	sink := ioutil.Discard
	if address != "" {
		sink = conn.NewDefaultManager("udp", address, log.NewNopLogger())
	}
	go s.StatsD.WriteLoop(s.Interval.C, sink)
	return s
}

// NewCounter
func (s StatsD) NewCounter(name string) *statsd.Counter {
	return s.StatsD.NewCounter(name+".count", s.SampleRate)
}

// NewGauge
func (s StatsD) NewGauge(name string) *statsd.Gauge {
	return s.StatsD.NewGauge(name + ".gauge")
}

// NewTiming
func (s StatsD) NewTiming(name string) *statsd.Timing {
	// Timings are interpreted as a value in milliseconds
	return s.StatsD.NewTiming(name+".duration", s.SampleRate)
}
