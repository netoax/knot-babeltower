package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CESARBR/knot-babeltower/pkg/logging"

	"github.com/gorilla/mux"
)

// Health represents the service's health status
type Health struct {
	Status string `json:"status"`
}

// Server represents the HTTP server
type Server struct {
	port   int
	logger logging.Logger
	srv    http.Server
}

// NewServer creates a new server instance
func NewServer(port int, logger logging.Logger) Server {
	return Server{port: port, logger: logger}
}

// Start starts the http server
func (s *Server) Start(started chan bool) {
	routers := s.createRouters()
	s.logger.Infof("Listening on %d", s.port)
	started <- true // Since ListenAndServe is blocking the channel must be set an instruction before call it
	s.srv = http.Server{Addr: fmt.Sprintf(":%d", s.port), Handler: s.logRequest(routers)}
	err := s.srv.ListenAndServe()
	if err != nil {
		s.logger.Error(err)
		started <- false
	}
}

// Stop stops the server
func (s *Server) Stop() {
	err := s.srv.Shutdown(context.TODO())
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Server) createRouters() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", s.healthcheckHandler)
	return r
}

func (s *Server) logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func (s *Server) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(&Health{Status: "online"})
	_, err := w.Write(response)
	if err != nil {
		s.logger.Errorf("Error sending response, %s\n", err)
	}
}
