package bimock

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Coordinates struct {
	Lat float64
	Lon float64
}
type RideInfo struct {
	Src       Coordinates
	Dest      Coordinates
	StartTime int64
	EndTime   int64
	RideID    string
	Distance  float64
}

//returns all the data in the collection as array of pointers
func composeData(collection *(mongo.Collection)) RideInfo {
	findOptions := options.Find()

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = cur.Close(context.Background()); err != nil {
			logrus.WithError(err).Fatal("Failed to close cursor")
		}
	}()

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

	rideinfo := RideInfo{
		RideID: (*results[0]).RideID,
		Src: Coordinates{
			Lat: (*results[0]).Lat,
			Lon: (*results[0]).Lon,
		},
		Dest: Coordinates{
			Lat: (*results[len(results)-1]).Lat,
			Lon: (*results[len(results)-1]).Lon,
		},
		StartTime: (*results[0]).TimeStamp,
		EndTime:   (*results[len(results)-1]).TimeStamp,
		Distance:  getDistance(collection),
	}

	return rideinfo
}
