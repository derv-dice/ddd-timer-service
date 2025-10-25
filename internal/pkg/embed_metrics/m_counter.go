package embed_metrics

import "sync/atomic"

type Counter struct {
	v atomic.Uint64
}

func NewCounter(startValue uint64) *Counter {
	if startValue == 0 {
		return new(Counter)
	}

	c := new(Counter)
	c.v.Store(startValue)

	return c
}

func (c *Counter) Inc() {
	c.v.Add(1)
}

func (c *Counter) Value() uint64 {
	return c.v.Load()
}
