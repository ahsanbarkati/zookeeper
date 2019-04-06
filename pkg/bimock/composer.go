package bimock

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type coordinates struct {
	Lat float64
	Lon float64
}
type rideInfo struct {
	src       coordinates
	dest      coordinates
	startTime int64
	endTime   int64
	rideID    string
	distance  float64
}

//returns all the data in the collection as array of pointers
func composeData(collection *(mongo.Collection)) rideInfo {
	findOptions := options.Find()

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	var results []*insData

	for cur.Next(context.TODO()) {
		var elem insData
		err = cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}

	err = cur.Close(context.TODO())
	if err != nil {
		logrus.WithError(err).Fatal("Failed to close cursor")
	}

	rideinfo := rideInfo{
		rideID: (*results[0]).RideID,
		src: coordinates{
			Lat: (*results[0]).Lat,
			Lon: (*results[0]).Lon,
		},
		dest: coordinates{
			Lat: (*results[len(results)-1]).Lat,
			Lon: (*results[len(results)-1]).Lon,
		},
		startTime: (*results[0]).TimeStamp,
		endTime:   (*results[len(results)-1]).TimeStamp,
		distance:  getDistance(collection),
	}

	return rideinfo

}
