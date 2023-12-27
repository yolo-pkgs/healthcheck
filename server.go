package healthcheck

import (
	"net/http"
	"sync/atomic"
)

type Server struct {
	s         *http.Server
	readyChan <-chan bool
	ready     atomic.Bool
}

func (s *Server) Serve(addr string) {
	go func(readyChan <-chan bool) {
		for ready := range readyChan {
			s.ready.Store(ready)
		}
	}(s.readyChan)

	s.s.ListenAndServe()
}

func New(addr string, ready <-chan bool) *Server {
	srv := Server{
		s:         &http.Server{Addr: addr},
		readyChan: ready,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/live", srv.liveHandle)
	mux.HandleFunc("/ready", srv.readyHandle)
	srv.s.Handler = mux

	return &srv
}

func (s *Server) liveHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) readyHandle(w http.ResponseWriter, r *http.Request) {
	if s.ready.Load() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
