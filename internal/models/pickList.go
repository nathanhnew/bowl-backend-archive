package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PickList struct {
	ID          primitive.ObjectID `bson: "_id"`
	User        primitive.ObjectID `bson: "user"`
	PickDate    time.Time          `bson: "pickDate"`
	Season      string             `bson: "season"`
	League      primitive.ObjectID `bson: "league"`
	IsFinalized *bool              `bson: "isFinalized"`
	Picks       []pick             `bson: "picks"`
}

type pick struct {
	Bowl        primitive.ObjectID `bson: "bowl"`
	Pick        primitive.ObjectID `bson: "pick"`
	Confidence  int                `bson: "confidence"`
	PointValue  float64            `bson: "pointValue"`
	HasOccurred bool               `bson: "hasOccurred"`
	IsCorrect   bool               `bson: "isCorrect"`
}
