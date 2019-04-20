package bimock

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type sensorData struct {
	RideID string
	Lat    float64
	Lon    float64
}

func readData(r *http.Request, s *Server, rideID string, collection *(mongo.Collection)) {
	decoder := json.NewDecoder(r.Body)
	var data sensorData
	err := decoder.Decode(&data)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse JSON")
	}
	dataToInsert := insData{
		RideID:    data.RideID,
		TimeStamp: int64(time.Now().Unix()),
		Lat:       data.Lat,
		Lon:       data.Lon,
	}

	insertToDB(collection, dataToInsert)
}
