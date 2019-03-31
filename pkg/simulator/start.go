package simulator

import (
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
	s.running.Lock()
	defer s.running.Unlock()

	// read from request, lat lon
	decoder := json.NewDecoder(r.Body)
	var loc StartRequest
	err := decoder.Decode(&loc)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse JSON")
	}

	generator, closer := makeGenerator(rideID, loc.Lat, loc.Lon)

	for {
		select {
		case data := <-generator:
			if err = postJSON(data); err != nil {
				logrus.WithError(err).Fatal("Failed to post JSON")
			}
		case <-s.stop:
			closer <- struct{}{}

		}
	}

}

func makeGenerator(rideID string, lat, lon float64) (chan SimData, chan struct{}) {
	rand.Seed(time.Now().UnixNano())

	output := make(chan SimData)
	closer := make(chan struct{})

	go func(rideID string, lat, lon float64) {
		curLoc := geo.NewPoint(lat, lon)

	FORLOOP:
		for {
			select {
			case <-closer:
				break FORLOOP
			default:
				rDist := float64(rand.Intn(10000)) / 1000.0
				rBear := float64(rand.Intn(36000)) / 100.0
				curLoc = curLoc.PointAtDistanceAndBearing(rDist, rBear)
				output <- SimData{RideID: rideID, Lat: curLoc.Lat(), Lon: curLoc.Lng()}
			}
		}
	}(rideID, lat, lon)

	return output, closer
}

func postJSON(data SimData) error {
	// TODO: APwhitehat
	return nil
}
