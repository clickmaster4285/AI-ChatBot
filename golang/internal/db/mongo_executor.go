package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoExecutor struct {
	Client *mongo.Client
	DBName string
}

func NewMongoExecutor(uri string, dbName string) (*MongoExecutor, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri).SetDirect(true))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &MongoExecutor{
		Client: client,
		DBName: dbName,
	}, nil
}

// SIMPLE SAFE COUNT QUERY
func (m *MongoExecutor) CountDocuments(collection string) (int64, error) {
	coll := m.Client.Database(m.DBName).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GENERIC FETCH (SAFE READ ONLY)
func (m *MongoExecutor) FindAll(collection string) ([]bson.M, error) {

	coll := m.Client.Database(m.DBName).Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var results []bson.M
	err = cursor.All(ctx, &results)

	return results, err
}

func (m *MongoExecutor) Debug(msg string) {
	log.Println("[MONGO EXECUTOR]", msg)
}
