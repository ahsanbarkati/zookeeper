package bimock

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

//struct sent to simulator, port is to specify the callback point to simulator
type location struct {
	Lat  float64
	Lon  float64
	Port string
}

func requestRide(r *http.Request, s *Server, rideID string) error {
	logrus.Info("ride requested")
	client := s.client
	sourceLoc := location{
		Lat:  15,
		Lon:  12,
		Port: s.port,
	}

	jsonData, err := json.Marshal(sourceLoc)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to do json marshal")
	}
	reqURL := "http://0.0.0.0:1000" + string((s.port)[3]) + "/start/" + rideID

	request, err := http.NewRequest(http.MethodGet, reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.WithError(err).Fatal("Unable to make start request")
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		logrus.WithError(err).Fatal("Bimock client's request failed")
	}

	return resp.Body.Close()
}
