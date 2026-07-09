package main

import (
	"log"
	"net/http"
	"time"
)

// statusRecorder wraps http.ResponseWriter to capture the status code and the
// number of bytes written, so the access log can report them.
type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	if s.status == 0 {
		s.status = http.StatusOK
	}
	n, err := s.ResponseWriter.Write(b)
	s.bytes += n
	return n, err
}

// noisyProbe reports paths that fire on every scrape/health-check tick, so
// they don't drown the access log at the default verbosity.
func noisyProbe(path string) bool {
	switch path {
	case "/health", "/readyz", "/metrics":
		return true
	}
	return false
}

// accessLog wraps a handler with one structured line per request:
//
//	atlas access: 200 POST /issue 412µs 127.0.0.1 1834b
//
// Probe/scrape endpoints are skipped unless verbose logging is enabled, so a
// health-checked deployment doesn't flood its logs.
func (a *App) accessLog(next http.Handler) http.Handler {
	if !a.logRequests {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if noisyProbe(r.URL.Path) && !a.logVerbose {
			next.ServeHTTP(w, r)
			return
		}
		start := a.clock.Now()
		rec := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)
		log.Printf("atlas access: %d %s %s %s %s %db",
			rec.status, r.Method, r.URL.Path,
			time.Since(start).Round(time.Microsecond), clientIP(r), rec.bytes)
	})
}
