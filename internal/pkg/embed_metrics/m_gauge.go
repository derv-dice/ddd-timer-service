package embed_metrics

import "sync/atomic"

type Gauge struct {
	v atomic.Int64
}

func NewGauge(startValue int64) *Gauge {
	if startValue == 0 {
		return new(Gauge)
	}

	g := new(Gauge)
	g.v.Store(startValue)

	return g
}

func (g *Gauge) Set(v int64) {
	g.v.Store(v)
}

func (g *Gauge) Inc() {
	g.v.Add(1)
}

func (g *Gauge) Dec() {
	g.v.Add(-1)
}
