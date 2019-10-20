package db

import (
	"context"
	"fmt"
	"github.com/nathanhnew/bowl-backend/internal/app/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var cfg, _ = config.GetConfig(config.DefaultConfigLocation)
var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var Client = Connect(cfg.GetMongoUri())
var database = cfg.GetMongoDB()
var userCollection = "User"
var schoolCollection = "School"
var leagueCollection = "League"
var bowlCollection = "Bowl"
var pickCollection = "PickList"

func Connect(uri string) *mongo.Client {
	ctx = getContext()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
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
	ctx, _ = context.WithTimeout(context.Background(), cfg.GetMongoTimeout()*time.Second)
	return ctx
}
