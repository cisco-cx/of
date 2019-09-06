package v1alpha1

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/cisco-cx/of/lib/v1alpha1"
)

// Ensures Counter implements v1alpha1.Counter
var _ v1alpha1.Counter = &Counter{}

type Counter struct {
	Namespace string
	Name      string
	Help      string
	cntr      prometheus.Counter
}

// Create a new counter.
func (c *Counter) Create() error {
	c.cntr = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: c.Namespace,
		Name:      c.Name,
		Help:      c.Help,
	})

	if c.cntr == nil {
		return v1alpha1.ErrCounterCreateFailed
	}

	return prometheus.Register(c.cntr)
}

// Increment the counter by 1.
func (c *Counter) Incr() error {
	c.cntr.Inc()
	return nil
}

// Remove counter..
func (c *Counter) Destroy() error {
	if prometheus.Unregister(c.cntr) {
		return nil
	}
	return v1alpha1.ErrCounterDestroyFailed
}
