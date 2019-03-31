package simulator

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	localhost = "0.0.0.0"
)

var (
	// PING can be used to ping a pingChannel
	PING = struct{}{}
)

type pingChannel chan struct{}

// Server is a simulator server + internal data for state.
type Server struct {
	router     *mux.Router
	client     *http.Client // keep a common http client to make requests
	port       string
	running    atomic.Value
	stop       pingChannel
	remoteAddr string
}

// NewServer creates and returns a new Server
func NewServer(port string) *Server {
	return &Server{
		port: port,
		client: &http.Client{
			Timeout: 5 * time.Second, // keep sane timeouts
		},
		stop: make(pingChannel),
	}
}

// SetupHTTP starts serving HTTP apis.
func (s *Server) SetupHTTP() {
	s.router = mux.NewRouter()
	setupRoutes(s)

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
