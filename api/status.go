package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ApplicationStatus contains information about a given application.
type ApplicationStatus struct {
	Version   string `json:"version"`
	BuildTime string `json:"buildTime"`
}

// StartStatus creates a handler on /status and listens on the given port.
// Use this if your application isn't a web app and just needs a simple status page.
func StartStatus(path, address string, port int, version, buildTime string) {
	http.HandleFunc(path, StatusHandler(version, buildTime))
	http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)
}

// StatusHandler returns an http.HandlerFunc that dumps out the status information of this application.
// Use this if you're running a web application with custom routing handlers.
func StatusHandler(version, buildTime string) http.HandlerFunc {
	httpStatus := http.StatusOK
	status := ApplicationStatus{
		Version:   version,
		BuildTime: buildTime,
	}

	b, err := json.Marshal(status)

	if err != nil {
		httpStatus = http.StatusInternalServerError
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		w.Write(b)
	}
}
