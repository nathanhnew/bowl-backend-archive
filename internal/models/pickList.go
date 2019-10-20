package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PickList struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	User        string             `json:"user" bson:"user"`
	PickDate    time.Time          `json:"pickDate" bson:"pickDate"`
	Season      string             `json:"season" bson:"season"`
	League      string             `json:"league" bson:"league"`
	IsFinalized bool               `json:"isFinalized" bson:"isFinalized"`
	Picks       []pick             `json:"picks" bson:"picks"`
	CreateTime  time.Time          `bson:"createTime"`
	UpdateTime  time.Time          `bson:"updateTime"`
}

type pick struct {
	Bowl        string  `json:"bowl" bson:"bowl"`
	Pick        string  `json:"pick" bson:"pick"`
	Confidence  int     `json:"confidence" bson:"confidence"`
	PointValue  float64 `json:"pointValue" bson:"pointValue"`
	HasOccurred bool    `json:"hasOccurred" bson:"hasOccurred"`
	IsCorrect   bool    `json:"isCorrect,omitempty" bson:"isCorrect,omitempty"`
}
