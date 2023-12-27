package healthcheck

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type Server struct {
	s         *http.Server
	readyChan <-chan bool
	ready     atomic.Bool
}

func (s *Server) Serve() error {
	go func(readyChan <-chan bool) {
		for ready := range readyChan {
			s.ready.Store(ready)
		}
	}(s.readyChan)

	return s.s.ListenAndServe()
}

// GET /live for liveness.
// GET /ready for readiness.
func New(addr string, prefix string, ready <-chan bool) *Server {
	srv := Server{
		s:         &http.Server{Addr: addr},
		readyChan: ready,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("%s/live", prefix), srv.liveHandle)
	mux.HandleFunc(fmt.Sprintf("%s/ready", prefix), srv.readyHandle)
	srv.s.Handler = mux

	return &srv
}

func (s *Server) liveHandle(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) readyHandle(w http.ResponseWriter, _ *http.Request) {
	if s.ready.Load() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
