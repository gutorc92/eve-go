package dao

import (
	"context"
	"fmt"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/gutorc92/api-farm/collections"
)

type DataMongo struct {
	Uri string
	Client *mongo.Client
	Database string
}

func NewDataMongo (uri string) (*DataMongo, error) {
	var dt DataMongo
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dt.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Error to connect")
		return nil, err
	}
	err = dt.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("Error to connect")
		return nil, err
	}
	return &dt, nil
}

func (dt *DataMongo) InsertOne() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := dt.Client.Database("testing").Collection("numbers")
	res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
	if err != nil {
		fmt.Println("Error to error")
	}
	id := res.InsertedID
	fmt.Println("Id insert: ", id)
}

func (dt *DataMongo) InsertFarm(farm collections.Farm) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := dt.Client.Database("testing").Collection("farm")
	res, err := collection.InsertOne(ctx, farm)
	if err != nil {
		fmt.Println("Error to error")
		return primitive.NewObjectID(), err
	}
	id := res.InsertedID.(primitive.ObjectID)
	// if err != nil {
	// 	return primitive.NewObjectID(), nil
	// }
	fmt.Println("Id insert: ", id)
	return id, nil
}

func (dt *DataMongo) FindFarm() ([]collections.Farm, error) {
	var farms []collections.Farm
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := dt.Client.Database("testing").Collection("farm")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &farms); err != nil {
		return nil, err
	}
	fmt.Println(farms)
	return farms, nil
}

