package prometheus

import (
	"net/http"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.rakops.com/rm/signal-api/stdlib/stats"
)

// Client holds a prometheus registry and allows for the setting up of the metrics endpoint
type Client struct {
	registry    prom.Registerer
	gatherer    prom.Gatherer
	collectors  Collectors
	handlerPath string
}

// NewClient returns a Prometheus Client while registering the given collectors
func NewClient(handlerPath string, collectors Collectors) *Client {
	registry := prom.NewRegistry()

	for _, collector := range collectors {
		registry.MustRegister(collector)
	}

	return &Client{
		registry:    registry,
		gatherer:    registry,
		collectors:  collectors,
		handlerPath: handlerPath,
	}
}

// NewDefaultClient returns a Prometheus Client that uses prometheus' built in DefaultRegisterer
func NewDefaultClient(handlerPath string, collectors Collectors) *Client {
	for _, collector := range collectors {
		prom.MustRegister(collector)
	}

	return &Client{
		registry:    prom.DefaultRegisterer,
		gatherer:    prom.DefaultGatherer,
		collectors:  collectors,
		handlerPath: handlerPath,
	}
}

// AddCollectors will add all given metrics collectors to the client's registry
func (c *Client) AddCollectors(collectors Collectors) {
	for key, collector := range collectors {
		c.collectors[key] = collector
		c.registry.MustRegister(collector)
	}
}

// RemoveCollectors will remove all given metrics collectors from the client's registry
func (c *Client) RemoveCollectors(collectors Collectors) {
	for key, col := range collectors {
		delete(c.collectors, key)
		c.registry.Unregister(col)
	}
}

// AddHandler will hand back the expected metrics endpoint and http handler for the client's registry
// This allows for any http handler to expose the endpoint. This is useful for Echo, Gin, pure, etc...
//
// Here's an example showing the setup process for Echo:
//
//  import (
//  	"github.com/labstack/echo/v4"
//
//  	"github.rakops.com/rm/signal-api/stdlib/stats/prometheus"
//  )
//
//  func main() {
//  	e := new(echo.Echo)
//  	statsClient := prometheus.NewClient("/prometheus_metrics", prometheus.Collectors{})
//  	statsClient.AddHandler(func(path string, handler http.Handler) {
//  		e.GET(path, echo.WrapHandler(handler))
//  	})
//
//  	e.Start("0.0.0.0:8080")
//  }
func (c *Client) AddHandler(callback func(string, http.Handler)) {
	callback(c.handlerPath, promhttp.HandlerFor(c.gatherer, promhttp.HandlerOpts{}))
}

func (c *Client) Timing(key string, labels stats.Labels, d time.Duration) error {
	collector, ok := c.collectors[key]
	if !ok {
		return ErrNoKey
	}

	l, err := labels.AsMap()
	if err != nil {
		return err
	}

	var observer prom.Observer

	switch c := collector.(type) {
	case *prom.HistogramVec:
		observer = c.With(l)
	case *prom.SummaryVec:
		observer = c.With(l)
	case prom.Summary:
		observer = c
	case prom.Histogram:
		observer = c
	default:
		return ErrInvalidType
	}

	observer.Observe(d.Seconds())

	return nil
}

func (c *Client) Incr(key string, labels stats.Labels, value int64) error {
	collector, ok := c.collectors[key]
	if !ok {
		return ErrNoKey
	}

	l, err := labels.AsMap()
	if err != nil {
		return err
	}

	var counter prom.Counter

	switch c := collector.(type) {
	case *prom.CounterVec:
		counter = c.With(l)
	case prom.Counter:
		counter = c
	default:
		return ErrInvalidType
	}

	counter.Add(float64(value))

	return nil
}

func (c *Client) Gauge(key string, labels stats.Labels, value float64) error {
	collector, ok := c.collectors[key]
	if !ok {
		return ErrNoKey
	}

	l, err := labels.AsMap()
	if err != nil {
		return err
	}

	var gauge prom.Gauge

	switch c := collector.(type) {
	case *prom.GaugeVec:
		gauge = c.With(l)
	case prom.Gauge:
		gauge = c
	default:
		return ErrInvalidType
	}

	gauge.Set(value)

	return nil
}
