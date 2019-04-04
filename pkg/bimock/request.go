package bimock

import (
	"net/http"
	"encoding/json"
	"bytes"
	"github.com/sirupsen/logrus"
)

type location struct {
	lat float64
	lon float64
}

func requestRide(r *http.Request, s *Server, rideID string) error  {
	logrus.Info("ride requested")
	client := s.client
	sourceLoc := location{
		lat:15,
		lon:12}

	jsonData, err := json.Marshal(sourceLoc)
	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:10000/start/" + rideID, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.WithError(err).Fatal("Unable to make start request")
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to make start request client")
	}
	return resp.Body.Close()
	// TODO
}