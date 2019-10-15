package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type League struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	Commissioner primitive.ObjectID `bson:"commissioner"`
	Seasons      []season           `bson:"seasons,omitempty"`
	Winners      []winner           `bson:"winners,omitempty"`
	CreateTime   time.Time          `bson:"createTime"`
	UpdateTime   time.Time          `bson:"updateTime"`
}

type winner struct {
	Winner primitive.ObjectID `bson:"winner"`
	Season string             `bson:"season"`
}

type season struct {
	Season             string               `bson:"season"`
	Players            []player             `bson:"players,omitempty"`
	Bowls              []primitive.ObjectID `bson:"bowls,omitempty"`
	SubmissionLockDate time.Time            `bson:"submissionLockDate"`
	Spread             spreadConfig         `bson:"spread"`
	Dues               duesConfig           `bson:"dues"`
}

type spreadConfig struct {
	UseSpread bool      `bson:"useSpread"`
	Simplify  string    `bson:"simplify"`
	LockDate  time.Time `bson:"lockDate"`
}

type duesConfig struct {
	HasDues   bool    `bson:"hasDues"`
	Amount    float64 `bson:"amount"`
	UseOnline bool    `bson:"useOnline"`
	Payout    payout  `bson:"payout,omitempty"`
}

type payout struct {
	Type    string `bson:"type"` // Fixed - fixed dollar value, percentage - pct of total money
	Amounts []struct {
		Rank   float64     `bson:"rank"`
		Amount interface{} `bson:"amount"`
	} `bson:"amounts"`
}

type player struct {
	User           primitive.ObjectID `bson:"user"`
	Status         string             `bson:"status"`
	HasSubmitted   bool               `bson:"hasSubmitted"`
	SubmissionDate time.Time          `bson:"submissionDate,omitempty"`
	HasPaid        bool               `bson:"hasPaid,omitempty"`
}
