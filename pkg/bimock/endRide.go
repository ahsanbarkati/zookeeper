package bimock

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func endRide(s *Server, collection *(mongo.Collection), rideID string) {

	timer := time.NewTimer(10 * time.Second)
	<-timer.C

	// data := composeData(collection)
	// Send this to backend

	reqURL := "http://0.0.0.0:1000" + string((s.port)[3]) + "/stop/" + rideID
	req, err := http.NewRequest(http.MethodGet, reqURL, bytes.NewBuffer(nil))
	if err != nil {
		logrus.WithError(err).Warn("Failed to make request")
	}

	c := s.client
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to send request")
	}

	_, err = collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to delete")
	}
	err = resp.Body.Close()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to close Response Body")
	}
	logrus.Info("Ride Ended")
}
