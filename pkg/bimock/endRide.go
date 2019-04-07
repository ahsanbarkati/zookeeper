package bimock

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const backendURL = "localhost:8080"

func endRide(s *Server, collection *(mongo.Collection), rideID string) {
	time.Sleep(10 * time.Second)

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

	// defer resp close to avoid leaks.
	defer func() {
		if err = resp.Body.Close(); err != nil {
			logrus.WithError(err).Fatal("Failed to close Response Body")
		}
	}()

	// save data to be sent to backend
	data := composeData(collection)
	_, err = collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to delete collection")
	}

	logrus.Info("Ride Ended")

	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to do marshal json")
	}

	reqURL = backendURL
	req, err = http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.WithError(err).Fatal("http request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = c.Do(req)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to POST Ride JSON")
	}

	if err = resp.Body.Close(); err != nil {
		logrus.WithError(err).Warn("Failed to close response, possible memmory leak")
	}

	logrus.Info("Ride JSON POST succeeded")
}
