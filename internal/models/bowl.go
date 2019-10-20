package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Bowl struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Slug       string             `json:"slug" bson:"slug"`
	Location   Location           `json:"location" bson:"location"`
	Games      []Game             `json:"games" bson:"games"`
	CreateTime time.Time          `bson:"createTime"`
	UpdateTime time.Time          `bson:"updateTime"`
}

type Game struct {
	Season   string    `json:"season" bson:"season"`
	Teams    [2]string `json:"teams" bson:"teams"`
	Spreads  []spread  `json:"spreads" bson:"spreads"`
	GameDate time.Time `json:"gameDate" bson:"gameDate"`
	EspnID   string    `json:"espnID" bson:"espnID"`
}

type spread struct {
	QueryDate time.Time `json:"queryDate" bson:"queryDate"`
	Favorite  string    `json:"favorite" bson:"favorite"`
	Spread    float64   `json:"spread" bson:"spread"`
}
