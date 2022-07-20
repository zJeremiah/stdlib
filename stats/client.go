package stats

import "time"

type (
	// Client is a generic stats collecting interface.
	Client interface {
		Timing(key string, labels Labels, d time.Duration) error
		Incr(key string, labels Labels, value int64) error
		Gauge(key string, labels Labels, value float64) error
	}

	// NoOpClient is an implementation of Client that does nothing.
	NoOpClient struct{}
)

// Timing does nothing.
func (n *NoOpClient) Timing(key string, labels Labels, d time.Duration) error {
	return nil
}

// Incr does nothing.
func (n *NoOpClient) Incr(key string, labels Labels, value int64) error {
	return nil
}

// Gauge does nothing.
func (n *NoOpClient) Gauge(key string, labels Labels, value float64) error {
	return nil
}
