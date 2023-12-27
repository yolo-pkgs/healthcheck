package healthcheck

import (
	"net/http"
)

type Server struct {
	s         *http.Server
	readyChan <-chan bool
	ready     bool
}

func (s *Server) Serve(addr string) {
	go func(readyChan <-chan bool) {
		for ready := range readyChan {
			s.ready = ready
		}
	}(s.readyChan)

	s.s.ListenAndServe()
}

func New(addr string, ready <-chan bool) *Server {
	srv := Server{
		s:         &http.Server{Addr: addr},
		readyChan: ready,
		ready:     false,
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
	if s.ready {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
