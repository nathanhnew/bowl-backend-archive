package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Bowl struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Location   Location           `bson:"location"`
	Games      []Game             `bson:"games"`
	CreateTime time.Time          `bson:"createTime"`
	UpdateTime time.Time          `bson:"updateTime"`
}

type Game struct {
	Season   string                `bson:"season"`
	Teams    [2]primitive.ObjectID `bson:"teams"`
	Spreads  []spread              `bson:"spreads"`
	GameDate time.Time             `bson:"gameDate"`
	EspnID   string                `bson:"espnID"`
}

type spread struct {
	QueryDate time.Time          `bson:"queryDate"`
	Favorite  primitive.ObjectID `bson:"favorite"`
	Spread    float64            `bson:"spread"`
}
