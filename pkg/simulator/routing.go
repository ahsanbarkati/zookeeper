package simulator

import (
	"net/http"

	"github.com/gorilla/mux"
)

func setupRoutes(s *Server, r *mux.Router) {
	r.HandleFunc("/start/{rideID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rideID := vars["rideID"]

		start(r, s, rideID)
	})

	r.HandleFunc("/stop/{rideID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rideID := vars["rideID"]

		stop(r, s, rideID)
	})
}
