package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func Connect(uri string) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	err := client.Ping(ctx, readpref.Primary())
	fmt.Println(err)
	return client
}

func GetContext(timeout time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), timeout*time.Second)
	return ctx
}
