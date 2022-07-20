package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	prom "github.com/prometheus/client_golang/prometheus"

	"github.rakops.com/rm/signal-api/stdlib/stats"
	"github.rakops.com/rm/signal-api/stdlib/stats/prometheus"
)

type (
	// StatsConfig configures this stats middleware.
	StatsConfig struct {
		Skipper middleware.Skipper
	}

	// SimpleStats uses the go http api to integrate with prometheus.
	SimpleStats struct {
		handler      http.Handler
		client       stats.Client
		urlSanitizer func(s string) string
	}
)

// NewSimpleStats returns an http middleware that records request durations.
// It uses the "simple_api_request_duration" collector.
//
// The param "urlSanitizer" is used to strip urls from sensitive information.
// Such as client names or ids. If you don't need to use this feature,
// just pass in nil. The provided string is simply the path part of the url.
func NewSimpleStats(handler http.Handler, client stats.Client, urlSanitizer func(s string) string) *SimpleStats {
	return &SimpleStats{
		handler:      handler,
		client:       client,
		urlSanitizer: urlSanitizer,
	}
}

// Stats is a middleware func that records request timing for every request.
func Stats(statsClient stats.Client) echo.MiddlewareFunc {
	return StatsWithConfig(statsClient, StatsConfig{})
}

// PrometheusCollectors is a prepopulated list of prometheus collectors.
// This must be used when using the stats collectors in this package
// with prometheus.
func PrometheusCollectors(app, team, env string) prometheus.Collectors {
	return prometheus.Collectors{
		"api_request_duration": prom.NewSummaryVec(
			prom.SummaryOpts{
				Name: "api_request_duration",
				Help: "The duration of each request",
				ConstLabels: prom.Labels{
					"app":  app,
					"team": team,
					"env":  env,
				},
			},
			[]string{"path", "code", "method"},
		),
		"simple_api_request_duration": prom.NewSummaryVec(
			prom.SummaryOpts{
				Name: "simple_api_request_duration",
				Help: "The duration of each request",
				ConstLabels: prom.Labels{
					"app":  app,
					"team": team,
					"env":  env,
				},
			},
			[]string{"path", "method"},
		),
	}
}

// StatsWithConfig returns a echo middleware that records the duration
// of any incoming HTTP request.
func StatsWithConfig(statsClient stats.Client, conf StatsConfig) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if conf.Skipper != nil && conf.Skipper(c) {
				return h(c)
			}

			start := time.Now()
			err := h(c)
			end := time.Now()
			response := c.Response()

			statsClient.Timing("api_request_duration", stats.Labels{
				"path", c.Path(),
				"code", strconv.Itoa(response.Status),
				"method", c.Request().Method,
			}, end.Sub(start))

			return err
		}
	}
}

func (s *SimpleStats) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	s.handler.ServeHTTP(w, r)
	end := time.Now()

	path := r.URL.Path

	if s.urlSanitizer != nil {
		path = s.urlSanitizer(path)
	}

	s.client.Timing("simple_api_request_duration", stats.Labels{
		"path", path,
		"method", r.Method,
	}, end.Sub(start))
}
