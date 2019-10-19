package db

import (
	"context"
	"fmt"
	"github.com/nathanhnew/bowl-backend/internal/app/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var cfg, _ = config.GetConfig(config.DefaultConfigLocation)
var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var Client *mongo.Client = Connect(cfg.GetMongoUri())
var database = cfg.Values["mongoDatabase"].(string)
var userCollection = "User"
var schoolCollection = "School"
var leagueCollection = "League"
var bowlCollection = "Bowl"
var pickCollection = "PickList"

func Connect(uri string) *mongo.Client {
	ctx = getContext()
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println(err)
	}
	return client
}

func Disconnect(client *mongo.Client) {
	ctx = getContext()
	client.Disconnect(ctx)
}

func getContext() context.Context {
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}
