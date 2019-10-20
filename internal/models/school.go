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
	Long         string `json:"longName" bson:"longName"`
	Short        string `json:"shortName" bson:"shortName"`
	Abbreviation string `json:"abbreviation" bson:"abbreviation"`
	Slug         string `json:"slug" bson:"slug"`
}

type conference struct {
	StartDate time.Time  `json:"startDate" bson:"startDate"`
	EndDate   *time.Time `json:"endDate" bson:"endDate,omitempty"`
	Name      string     `json:"name" bson:"name"`
	Slug      string     `json:"slug" bson:"slug"`
}

func (school *School) UpdateFromMap(payload map[string]interface{}) {
	for k, v := range payload {
		if k == "longName" {
			school.Name.Long = v.(string)
		}
		if k == "shortName" {
			school.Name.Short = v.(string)
		}
		if k == "abbreviation" {
			school.Name.Abbreviation = v.(string)
		}
		if k == "slug" {
			school.Name.Slug = v.(string)
		}
		if k == "teamName" {
			school.TeamName = v.(string)
		}
		if k == "logo" {

		}
	}
}
