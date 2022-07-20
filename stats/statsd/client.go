package statsd

import (
	"fmt"
	"strings"
	"time"

	"github.com/quipo/statsd"

	"github.rakops.com/rm/signal-api/stdlib/stats"
)

// Client implements the stats.Client interface. This will send
// stats data to a StatsD server.
type Client struct {
	address  string
	interval time.Duration
	prefix   string

	client statsd.Statsd
	buffer statsd.Statsd
}

// NewClient returns an implementation of Client that
// writes stats to a statsd service.
func NewClient(statsdAddress, prefix string, interval time.Duration) *Client {
	if !strings.HasSuffix(prefix, ".") {
		prefix = prefix + "."
	}

	return &Client{
		address:  statsdAddress,
		interval: interval,
		prefix:   prefix,
	}
}

// Open opens the connection to the StatsD server.
func (s *Client) Open() error {
	client := statsd.NewStatsdClient(s.address, s.prefix)

	if err := client.CreateSocket(); err != nil {
		return err
	}

	buffer := statsd.NewStatsdBuffer(s.interval, client)

	s.client = client
	s.buffer = buffer

	return nil
}

// Statsd returns the underlying StatsD implementation.
func (s *Client) Statsd() statsd.Statsd {
	return s.buffer
}

// Close cleans up connections.
func (s *Client) Close() error {
	if err := s.buffer.Close(); err != nil {
		return err
	}

	return s.client.Close()
}

// Timing sends timing data to the StatsD server.
func (s *Client) Timing(key string, labels stats.Labels, d time.Duration) error {
	return s.Statsd().PrecisionTiming(s.key(key, labels), d)
}

// Incr sends counter data to the StatsD server.
func (s *Client) Incr(key string, labels stats.Labels, value int64) error {
	return s.Statsd().Incr(s.key(key, labels), value)
}

// Gauge sends gauge data to the StatsD server.
func (s *Client) Gauge(key string, labels stats.Labels, value float64) error {
	return s.Statsd().Gauge(s.key(key, labels), int64(value))
}

func (s *Client) key(key string, labels stats.Labels) string {
	key = strings.ToLower(key)

	if len(labels) > 0 {
		replacer := strings.NewReplacer(".", "-", ":", "-")

		for i, label := range labels {
			labels[i] = strings.ToLower(replacer.Replace(label))
		}

		return fmt.Sprintf("%s.%s", key, strings.Join(labels, "."))
	}

	return key
}
