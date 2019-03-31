package simulator

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	geo "github.com/kellydunn/golang-geo"
	"github.com/sirupsen/logrus"
)

// StartRequest is the JSON format of request body.
type StartRequest struct {
	Lat float64
	Lon float64
}

// SimData is a instance of data which is sent back to bimock.
type SimData struct {
	RideID string
	Lat    float64
	Lon    float64
}

func start(r *http.Request, s *Server, rideID string) {
	logrus.Info("start recieved")
	s.running.Store(true)

	logrus.Infof("Received request from: %s", r.RemoteAddr)
	s.remoteAddr = r.RemoteAddr

	// read from request, lat lon
	decoder := json.NewDecoder(r.Body)
	var loc StartRequest
	err := decoder.Decode(&loc)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse JSON")
	}

	go runGenerator(s, rideID, loc.Lat, loc.Lon)
}

func runGenerator(s *Server, rideID string, lat, lon float64) {
	rand.Seed(time.Now().UnixNano())

	output := make(chan SimData)
	closer := make(pingChannel)

	go func(rideID string, lat, lon float64) {
		curLoc := geo.NewPoint(lat, lon)
		tick := time.NewTicker(10 * time.Nanosecond)

	FORLOOP:
		for {
			select {
			case <-closer:
				tick.Stop()
				break FORLOOP
			case <-tick.C:
				rDist := float64(rand.Intn(100000)) / 1000.0
				rBear := float64(rand.Intn(36000)) / 100.0
				curLoc = curLoc.PointAtDistanceAndBearing(rDist, rBear)
				output <- SimData{RideID: rideID, Lat: curLoc.Lat(), Lon: curLoc.Lng()}
			}
		}
	}(rideID, lat, lon)

	for {
		select {
		case data := <-output:
			if err := postJSON(s.client, s.remoteAddr, data); err != nil {
				logrus.WithError(err).Fatal("Failed to post JSON")
			}
		case <-s.stop:
			closer <- PING
		}
	}
}

func postJSON(c *http.Client, remote string, data SimData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to send request")
	}

	req, err := http.NewRequest(http.MethodPost, `http://`+remote+"/sensorData", bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.WithError(err).Warn("Failed to make request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to send request")
	}

	return resp.Body.Close()
}
