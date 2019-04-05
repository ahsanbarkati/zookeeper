package bimock

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type insData struct {
	RideID    string
	TimeStamp int64
	Lat       float64
	Lon       float64
}

func setupDB(port string) *(mongo.Collection) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to Connect to MongoDB")
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to ping to MongoDB")
	}

	logrus.Info("Connected to MongoDB")

	// Get a handle for your collection
	collection := client.Database("bikeDB" + port).Collection("bikeHistory")
	_, err = collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to delete")
	}

	logrus.Info("Created collection")

	return collection
}

func insertToDB(collection *(mongo.Collection), data insData) {
	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to insert data")
	}
	logrus.Info("Inserted: ", insertResult.InsertedID)
}
