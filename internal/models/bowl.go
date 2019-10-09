package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Bowl struct {
	ID       primitive.ObjectID `bson: "_id"`
	Name     string             `bson: "name"`
	Location location           `bson: "location"`
	Games    []game             `bson: "games"`
}

type game struct {
	Season   string                `bson: "season"`
	Teams    [2]primitive.ObjectID `bson: "teams"`
	Spreads  []spread              `bson: "spreads"`
	GameDate time.Time             `bson: "gameDate"`
}

type spread struct {
	QueryDate time.Time          `bson: "queryDate"`
	Favorite  primitive.ObjectID `bson: "favorite"`
	Spread    float64            `bson: "spread"`
}
