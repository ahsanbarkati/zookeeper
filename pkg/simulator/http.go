package simulator

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	localhost = "0.0.0.0"
)

type pingChannel chan struct{}

type server struct {
	router  *mux.Router
	port    string
	running sync.Mutex
	stop    pingChannel
}

// NewServer creates and returns a new server
func NewServer(r *mux.Router, port string) *server {
	return &server{
		port: port,
	}
}

func (s *server) setupRoutes() {
	r := mux.NewRouter()
	setupRoutes(s, r)
}

func (s *server) serve() {
	hs := &http.Server{
		Addr:           localhost + ":" + s.port,
		Handler:        s.router,
		ReadTimeout:    10 * time.Second, // keep sane timouts
		WriteTimeout:   10 * time.Second, // keep sane timouts
		MaxHeaderBytes: 1 << 20,
	}

	if err := hs.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal()
	}
}

func (s *server) SetupHTTP() {
	s.setupRoutes()
	s.serve()
}
