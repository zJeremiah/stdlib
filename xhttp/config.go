package xhttp

import "fmt"

// Config contains configuration to start an http server.
type Config struct {
	Address string `default:"0.0.0.0"`
	Port    int    `default:"8080"`
}

// ListenAddress returns a formatted string to be passed into
// something like http.ListenAndServe. It's of the format:
// 0.0.0.0:8080
func (h *Config) ListenAddress() string {
	return fmt.Sprintf("%s:%d", h.Address, h.Port)
}
