package embed_metrics

import "sync"

type Collector struct {
	mx sync.Mutex
	c  map[string]*Counter
	g  map[string]*Gauge
}

func NewCollector() *Collector {
	return &Collector{
		c: make(map[string]*Counter),
		g: make(map[string]*Gauge),
	}
}

func (c *Collector) RegisterCounter(name string, startValue uint64) *Counter {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.c[name] = NewCounter(startValue)

	panic("implement me later") // TODO
}

func (c *Collector) RegisterGauge(name string) *Gauge {

	panic("implement me later") // TODO
}
