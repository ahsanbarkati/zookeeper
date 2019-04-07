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
	Lat  float64
	Lon  float64
	Port string
}

// SimData is a instance of data which is sent back to bimock.
type SimData struct {
	RideID string
	Lat    float64
	Lon    float64
}

func start(r *http.Request, s *Server, rideID string) {
	s.running.Store(true)

	s.remoteAddr = r.RemoteAddr

	// read from request, lat lon
	decoder := json.NewDecoder(r.Body)
	var sReq StartRequest
	err := decoder.Decode(&sReq)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse JSON")
	}

	logrus.Infof("Received request from port: %s, rideID: %s", sReq.Port, rideID)

	go runGenerator(s, rideID, sReq.Port, sReq.Lat, sReq.Lon)
}

func runGenerator(s *Server, rideID string, bimockPort string, lat, lon float64) {
	rand.Seed(time.Now().UnixNano())

	output := make(chan SimData)
	closer := make(pingChannel)

	go func() {
		curLoc := geo.NewPoint(lat, lon)
		tick := time.NewTicker(2 * time.Second)

	FORGENERATOR:
		for {
			select {
			case <-closer:
				tick.Stop()
				break FORGENERATOR
			case <-tick.C:
				rDist := float64(rand.Intn(100000)) / 1000.0
				rBear := float64(rand.Intn(36000)) / 100.0
				curLoc = curLoc.PointAtDistanceAndBearing(rDist, rBear)
				output <- SimData{RideID: rideID, Lat: curLoc.Lat(), Lon: curLoc.Lng()}
			}
		}
	}()

FORSELECT:
	for {
		select {
		case data := <-output:
			if err := postJSON(s.client, bimockPort, data); err != nil {
				logrus.WithError(err).Fatal("Failed to post JSON")
			}
		case <-s.stop:
			closer <- PING
			break FORSELECT
		}
	}
}

func postJSON(c *http.Client, bimockPort string, data SimData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to send request")
	}
	reqURL := "http://0.0.0.0:" + bimockPort + "/sensorData/" + data.RideID
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(jsonData))
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
