package main

import (
	"context"
	"fmt"
	"github.com/matthewboyd/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoResult struct {
	Geometry   models.Geometry   `json:"geometry" bson:"geometry"`
	Properties models.Properties `json:"properties" bson:"properties"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalln("Unable to connect to the mongo database: %v", err)
	}
	collection := client.Database("countries").Collection("NorthernIreland")
	collection2 := client.Database("countries").Collection("NorthernIrelandAttractions")

	coords := &MongoResult{}
	coords.getCoordinates(ctx, collection, "BT61 9JG")
	coords.insertActivity(ctx, collection2, "Matthew's house")
	fmt.Println("coordinates", coords)
}

func (co *MongoResult) getCoordinates(ctx context.Context, collection *mongo.Collection, postcode string) {
	filter := bson.D{{"properties.Postcode", postcode}}
	rawBytes, err := collection.FindOne(ctx, filter).DecodeBytes()
	//Some of the data isn't unmarshalling properly.
	err = bson.Unmarshal(rawBytes, &co)
	if err != nil {
		log.Fatalf("Error fetching postcode, %v", err)
	}
	log.Println("Mongo results", *co)
}

func (co *MongoResult) insertActivity(ctx context.Context, collection *mongo.Collection, name string) {

	document := bson.D{
		{"geometry", bson.D{
			{"type", "Point"},
			{"coordinates", bson.A{
				co.Geometry.Coordinates[0], co.Geometry.Coordinates[1],
			}},
		}},
		{"properties", bson.D{
			{"name", name},
		}},
	}
	one, err := collection.InsertOne(ctx, document)

	log.Println("The id of the object that was inserted %v", one.InsertedID)
	if err != nil {
		fmt.Sprintf("There was an error when inserting into the db: %v", err)
	}
}
