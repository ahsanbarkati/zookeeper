package bimock

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	localhost = "0.0.0.0"
)

// Server is a bimock server + internal data for state.
type Server struct {
	router *mux.Router
	client *http.Client // keep a common http client to make requests
	port   string
}

// NewServer creates and returns a new Server
func NewServer(port string) *Server {
	return &Server{
		port: port,
		client: &http.Client{
			Timeout: 5 * time.Second, // keep sane timeouts
		},
	}
}

// SetupHTTP starts serving HTTP apis.
func (s *Server) SetupHTTP() {
	s.router = mux.NewRouter()
	collection := setupDB(s.port)
	setupRoutes(s, collection)

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
