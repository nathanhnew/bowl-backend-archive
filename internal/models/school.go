package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type School struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        schoolName         `json:"name" bson:"name"`
	TeamName    string             `json:"teamName" bson:"teamName"`
	Logo        string             `json:"logo" bson:"logo"`
	Conferences []conference       `json:"conferences" bson:"conferences"`
	Colors      colorScheme        `json:"colors" bson:"colors"`
	Location    Location           `json:"location" bson:"location"`
	CreateTime  time.Time          `bson:"createTime"`
	UpdateTime  time.Time          `bson:"updateTime"`
}

type schoolName struct {
	LongName     string `bson:"longName"`
	ShortName    string `bson:"shortName"`
	Abbreviation string `bson:"abbreviation"`
}

type conference struct {
	StartDate time.Time  `bson:"startDate"`
	EndDate   *time.Time `bson:"endDate"`
	Name      string     `bson:"name"`
}
