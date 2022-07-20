package prometheus

// Config defines config options for prometheus.
// The tags used here are meant to be used with
// the library github.com/koding/multiconfig
type Config struct {
	Enabled bool   `default:"true"`
	Path    string `default:"/prometheus_metrics"`
}
