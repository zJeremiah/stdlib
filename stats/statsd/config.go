package statsd

import "time"

// Config defines config options for statsd.
type Config struct {
	Enabled  bool   `default:"false"`
	Address  string `default:"localhost:8125"`
	Interval int    `default:"2"`
}

// Duration returns the Interval on Statsd as a time.Duration in seconds.
func (s *Config) Duration() time.Duration {
	return time.Second * time.Duration(s.Interval)
}
