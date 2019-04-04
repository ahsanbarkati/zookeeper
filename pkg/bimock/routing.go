package bimock

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
)

func setupRoutes(s *Server) {
	s.router.HandleFunc("/sensorData/{rideID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rideID := vars["rideID"]

		readData(r, s, rideID)
	}).Methods(http.MethodPost, http.MethodOptions)

	s.router.HandleFunc("/requestRide/{rideID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rideID := vars["rideID"]

		if err := requestRide(r, s, rideID); err != nil {
				logrus.WithError(err).Fatal("Failed to request ride")
			}
	}).Methods(http.MethodGet, http.MethodOptions)
}
