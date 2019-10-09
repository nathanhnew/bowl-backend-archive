package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type School struct {
	ID         primitive.ObjectID `bson: "_id"`
	Name       schoolName         `bson: "name"`
	Conference []conference       `bson: "conferences"`
	Colors     colorScheme        `bson: "colors"`
	Location   location           `bson: "location"`
	TeamName   string             `bson: "teamName"`
	Logo       string             `bson: "logo"`
}

type schoolName struct {
	LongName     string `bson: "longName"`
	ShortName    string `bson: "shortName"`
	Abbreviation string `bson: "abbreviation"`
}

type conference struct {
	StartDate time.Time `bson: "startDate"`
	EndDate   time.Time `bson: "endDate"`
	Name      string    `bson: "name"`
}
