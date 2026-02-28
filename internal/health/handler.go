package health

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

type State struct {
	ready atomic.Bool
}

func NewState() *State {
	return &State{}
}

func (s *State) SetReady(ready bool) {
	s.ready.Store(ready)
}

func (s *State) IsReady() bool {
	return s.ready.Load()
}

type statusResponse struct {
	Checks    map[string]string `json:"checks,omitempty"`
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
}

func NewHealthServer(state *State) *http.Server {
	if state == nil {
		state = NewState()
	}

	addr := os.Getenv("HEALTH_LISTEN_ADDR")
	if addr == "" {
		addr = ":8081"
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/livez",
		getRequestOnly(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(w, 200, statusResponse{
				Status:    "ok",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})
		})),
	)
	mux.Handle(
		"/readyz",
		getRequestOnly(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if state.IsReady() {
				writeJSON(w, 200, statusResponse{
					Status:    "ok",
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Checks: map[string]string{
						"readiness": "ready",
					},
				})
				return
			}

			writeJSON(w, 503, statusResponse{
				Status:    "error",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Checks: map[string]string{
					"readiness": "not ready",
				},
			})
		})),
	)
	mux.Handle(
		"/healthz",
		getRequestOnly(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if state.IsReady() {
				writeJSON(w, 200, statusResponse{
					Status:    "ok",
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Checks: map[string]string{
						"liveness":  "alive",
						"readiness": "ready",
					},
				})
				return
			}

			writeJSON(w, 503, statusResponse{
				Status:    "error",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Checks: map[string]string{
					"liveness":  "alive",
					"readiness": "not ready",
				},
			})
		})),
	)

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func getRequestOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			writeJSON(w, 405, statusResponse{
				Status:    "error",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Checks: map[string]string{
					"method": "only GET is allowed",
				},
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, payload statusResponse) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, _ = w.Write(buf.Bytes())
}
