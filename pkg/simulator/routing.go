package simulator

import (
	"net/http"

	"github.com/gorilla/mux"
)

func setupRoutes(s *Server) {
	s.router.HandleFunc("/start/{rideID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rideID := vars["rideID"]

		start(r, s, rideID)
	}).Methods(http.MethodGet, http.MethodOptions)

	s.router.HandleFunc("/stop/{rideID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rideID := vars["rideID"]

		stop(r, s, rideID)
	}).Methods(http.MethodGet, http.MethodOptions)
}
