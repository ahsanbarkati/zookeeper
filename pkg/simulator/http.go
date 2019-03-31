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

// Server is a simulator server + internal data for state.
type Server struct {
	router  *mux.Router
	port    string
	running sync.Mutex
	stop    pingChannel
}

// NewServer creates and returns a new Server
func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) setupRoutes() {
	r := mux.NewRouter()
	setupRoutes(s, r)
}

func (s *Server) serve() {
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

// SetupHTTP starts serving HTTP apis.
func (s *Server) SetupHTTP() {
	s.setupRoutes()
	s.serve()
}
