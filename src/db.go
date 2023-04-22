package src

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitDB(mongoURI, dbName string, timeout int) (*mongo.Client, error) {
	uri := fmt.Sprintf("%s/%s", mongoURI, dbName)
	return Connect(uri, timeout)
}

func Connect(mongoURI string, timeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	readPreference := readpref.SecondaryPreferred(readpref.WithMaxStaleness(90 * time.Second))

	clientOption := &options.ClientOptions{
		ReadPreference: readPreference,
	}

	client, err := mongo.Connect(ctx, clientOption.ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readPreference); err != nil {
		return nil, err
	}

	return client, nil
}
