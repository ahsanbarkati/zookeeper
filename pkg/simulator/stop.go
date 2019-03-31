package simulator

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func stop(r *http.Request, s *Server, rideID string) {
	logrus.Info(s.running.Load())
	if s.running.Load() != nil && s.running.Load().(bool) {
		s.stop <- PING
		s.running.Store(false)
	} else {
		logrus.Info("Simulator not running")
	}
}
