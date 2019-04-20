package bimock

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupRoutes(s *Server, collection *(mongo.Collection)) {
	s.router.HandleFunc("/sensorData/{rideID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rideID := vars["rideID"]
		readData(r, s, rideID, collection)
	}).Methods(http.MethodPost, http.MethodOptions)

	s.router.HandleFunc("/requestRide/{rideID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rideID := vars["rideID"]

		if err := requestRide(r, s, rideID); err != nil {
			logrus.WithError(err).Fatal("Failed to request ride")
		}

		go endRide(s, collection, rideID)

	}).Methods(http.MethodGet, http.MethodOptions)
}
