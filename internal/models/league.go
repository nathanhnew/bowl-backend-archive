package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type League struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Slug         string             `json:"slug" bson:"slug"`
	Commissioner string             `json:"commissioner" bson:"commissioner"`
	Seasons      []season           `json:"seasons" bson:"seasons,omitempty"`
	Winners      []winner           `json:"winners" bson:"winners,omitempty"`
	CreateTime   time.Time          `json:"createTime" bson:"createTime"`
	UpdateTime   time.Time          `json:"updateTime" bson:"updateTime"`
}

type winner struct {
	Winner primitive.ObjectID `json:"winner" bson:"winner"`
	Season string             `json:"season" bson:"season"`
}

type season struct {
	Season             string               `json:"season" bson:"season"`
	Players            []player             `json:"players" bson:"players,omitempty"`
	Bowls              []primitive.ObjectID `json:"bowls" bson:"bowls,omitempty"`
	SubmissionLockDate time.Time            `json:"submissionLockDate" bson:"submissionLockDate"`
	Spread             spreadConfig         `json:"spread" bson:"spread"`
	Dues               duesConfig           `json:"dues" bson:"dues"`
}

type spreadConfig struct {
	UseSpread bool      `json:"useSpread" bson:"useSpread"`
	Simplify  string    `json:"simplify" bson:"simplify"`
	LockDate  time.Time `json:"lockDate" bson:"lockDate"`
}

type duesConfig struct {
	HasDues   bool    `json:"hasDues" bson:"hasDues"`
	Amount    float64 `json:"amount" bson:"amount"`
	UseOnline bool    `json:"useOnline" bson:"useOnline"`
	Payout    payout  `json:"payout" bson:"payout,omitempty"`
}

type payout struct {
	Type    string `bson:"type"` // fixed - fixed dollar value, percentage - pct of total money
	Amounts []struct {
		Rank   float64     `json:"rank" bson:"rank"`
		Amount interface{} `json:"amount" bson:"amount"`
	} `json:"amounts" bson:"amounts"`
}

type player struct {
	User           string    `json:"user" bson:"user"`
	Status         string    `json:"status" bson:"status"`
	HasSubmitted   bool      `json:"hasSubmitted" bson:"hasSubmitted"`
	SubmissionDate time.Time `json:"submissionDate" bson:"submissionDate,omitempty"`
	HasPaid        bool      `json:"hasPaid" bson:"hasPaid,omitempty"`
}
