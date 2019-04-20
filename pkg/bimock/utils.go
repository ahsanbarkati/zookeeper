package bimock

import (
	"context"
	"math"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Given latitude and longitude of two points : Returns distance
// Set unit parameter to "K" for kilometers or  "N" for nautical miles
func distance(lat1 float64, lon1 float64, lat2 float64, lon2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lon1 - lon2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}

// Given the relevant collection, returns the distance travelled from the start
func getDistance(collection *(mongo.Collection)) float64 {

	findOptions := options.Find()

	// Finding multiple documents returns a cursor
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create cursor")
	}
	defer func() {
		// Close the cursor once finished
		if err = cur.Close(context.Background()); err != nil {
			logrus.WithError(err).Fatal("Failed to close cursor")
		}
	}()

	// Iterate through the cursor
	var prev, curr insData
	dist := 0.0

	for cur.Next(context.TODO()) {
		err = cur.Decode(&curr)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to decode to BSON")
		}
		if prev != (insData{}) {
			dist += distance(prev.Lat, prev.Lon, curr.Lat, curr.Lon, "K")
		}
		prev = curr
	}

	err = cur.Err()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to close cursor")
	}

	return dist
}
