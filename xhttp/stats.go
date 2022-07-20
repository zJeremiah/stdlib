package xhttp

import (
	"io"
	"net/http"
	"net/url"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.rakops.com/rm/signal-api/stdlib/stats"
	"github.rakops.com/rm/signal-api/stdlib/stats/prometheus"
)

// StatsClient is a wrapper around Client that records length
// of requests and reports them to statsd.
type StatsClient struct {
	httpClient  Client
	statsClient stats.Client
}

// NewDefaultStatsClient returns a new StatsClient with the
// default options.
func NewDefaultStatsClient(statsClient stats.Client) *StatsClient {
	return &StatsClient{
		statsClient: statsClient,
		httpClient:  NewDefaultClient(),
	}
}

// NewStatsClient returns a new StatsClient with the given
// options.
func NewStatsClient(statsClient stats.Client, transport http.RoundTripper, jar http.CookieJar, timeout time.Duration) *StatsClient {
	return &StatsClient{
		httpClient:  NewClient(transport, jar, timeout),
		statsClient: statsClient,
	}
}

// PrometheusCollectors is a prepopulated list of prometheus collectors.
// This must be used when using the stats collectors in this package
// with prometheus.
func PrometheusCollectors(app, team, env string) prometheus.Collectors {
	return prometheus.Collectors{
		"http_client_request": prom.NewSummaryVec(
			prom.SummaryOpts{
				Name: "http_client_request",
				Help: "The duration of an http client request",
				ConstLabels: prom.Labels{
					"app":  app,
					"team": team,
					"env":  env,
				},
			},
			[]string{"hostname", "path", "method"},
		),
	}
}

// Do calls the underlying http client's Post method and
// records how long that request took.
func (s *StatsClient) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := s.httpClient.Do(req)
	end := time.Now()

	s.statsClient.Timing("http_client_request", s.labels(req.Method, req.URL), end.Sub(start))

	return resp, err
}

// Get calls the underlying http client's Get method and
// records how long that request took.
func (s *StatsClient) Get(u string) (*http.Response, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	resp, err := s.httpClient.Get(u)
	end := time.Now()

	s.statsClient.Timing("http_client_request", s.labels("GET", parsedURL), end.Sub(start))

	return resp, err
}

// Head calls the underlying http client's Head method and
// records how long that request took.
func (s *StatsClient) Head(u string) (*http.Response, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	resp, err := s.httpClient.Head(u)
	end := time.Now()

	s.statsClient.Timing("http_client_request", s.labels("HEAD", parsedURL), end.Sub(start))

	return resp, err
}

// Post calls the underlying http client's Post method and
// records how long that request took.
func (s *StatsClient) Post(u, bodyType string, body io.Reader) (*http.Response, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	resp, err := s.httpClient.Post(u, bodyType, body)
	end := time.Now()

	s.statsClient.Timing("http_client_request", s.labels("POST", parsedURL), end.Sub(start))

	return resp, err
}

// PostForm calls the underlying http client's PostForm method
// and records how long that request took.
func (s *StatsClient) PostForm(u string, data url.Values) (*http.Response, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	resp, err := s.httpClient.PostForm(u, data)
	end := time.Now()

	s.statsClient.Timing("http_client_request", s.labels("POST", parsedURL), end.Sub(start))

	return resp, err
}

func (s *StatsClient) labels(method string, u *url.URL) stats.Labels {
	path := u.Path

	if len(path) == 0 {
		path = "/"
	}

	hostname := u.Hostname()

	if len(hostname) == 0 {
		hostname = "unknown-host"
	}

	return stats.Labels{
		"hostname", hostname,
		"path", path,
		"method", method,
	}
}
