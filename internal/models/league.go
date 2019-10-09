package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type League struct {
	ID           primitive.ObjectID `bson: "_id"`
	Name         string             `bson: "name"`
	Commissioner primitive.ObjectID `bson: "commissioner"`
	Seasons      []season           `bson: "seasons"`
	Winners      []winner           `bson: "winners"`
}

type winner struct {
	Winner primitive.ObjectID `bson: "winner"`
	Season string             `bson: "season"`
}

type season struct {
	Season             string               `bson: "season"`
	Spread             spreadConfig         `bson: "spread"`
	Dues               duesConfig           `bson: "dues"`
	SubmissionLockDate time.Time            `bson: "submissionLockDate"`
	Bowls              []primitive.ObjectID `bson: "bowls"`
	Players            []player             `bson: "players"`
}

type spreadConfig struct {
	UseSpread bool      `bson: "useSpread"`
	Simplify  string    `bson: "simplify"`
	LockDate  time.Time `bson: "lockDate"`
}

type duesConfig struct {
	HasDues   bool    `bson: "hasDues"`
	Amount    float64 `bson: "amount"`
	UseOnline bool    `bson: "useOnline"`
}

type player struct {
	User           primitive.ObjectID `bson: "user"`
	Status         string             `bson: "status"`
	HasSubmitted   bool               `bson: "hasSubmitted"`
	SubmissionDate time.Time          `bson: "submissionDate"`
	HasPaid        bool               `bson: "hasPaid"`
}
