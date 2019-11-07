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

type BowlHeadline struct {
	Name     string       `json:"name" bson:"name"`
	Slug     string       `json:"slug" bson:"slug"`
	Location Location     `json:"location" bson:"location"`
	Game     gameHeadline `json:"game" bson:"games"`
}

type Game struct {
	Season   string    `json:"season" bson:"season"`
	Teams    [2]string `json:"teams" bson:"teams"`
	Spreads  []spread  `json:"spreads" bson:"spreads"`
	GameDate time.Time `json:"gameDate" bson:"gameDate"`
	EspnID   string    `json:"espnID" bson:"espnID"`
}

type gameHeadline struct {
	Season   string    `json:"season" bson:"season"`
	Teams    [2]string `json:"teams" bson:"teams"`
	Spreads  spread    `json:"spreads" bson:"spreads"`
	GameDate time.Time `json:"gameDate" bson:"gameDate"`
	EspnID   string    `json:"espnID" bson:"espnID"`
}

type spread struct {
	QueryDate time.Time `json:"queryDate" bson:"queryDate"`
	Favorite  string    `json:"favorite" bson:"favorite"`
	Spread    float64   `json:"spread" bson:"spread"`
}

func NewBowl() Bowl {
	bowl := Bowl{}
	bowl.CreateTime = time.Now()
	bowl.UpdateTime = time.Now()
	return bowl
}

func (bowl *Bowl) UpdateFromMap(payload map[string]interface{}) {
	for top, topValue := range payload {
		if top == "name" {
			bowl.Name = topValue.(string)
		} else if top == "location" {
			for sub, subValue := range topValue.(map[string]interface{}) {
				if sub == "city" {
					bowl.Location.City = subValue.(string)
				} else if sub == "state" {
					bowl.Location.State = subValue.(string)
				} else if sub == "geometry" {
					geo := subValue.(map[string]interface{})
					if t, ok := geo["type"].(string); ok {
						bowl.Location.Geometry.Type = t
					} else if coords, ok := geo["coordinates"].([]float64); ok {
						bowl.Location.Geometry.Coordinates = coords
					}
				}
			}
		} else if top == "games" {
			for ind, game := range topValue.([]map[string]interface{}) {
				for sub, subValue := range game {
					if sub == "season" {
						bowl.Games[ind].Season = subValue.(string)
					} else if sub == "teams" {
						bowl.Games[ind].Teams = subValue.([2]string)
					} else if sub == "gameDate" {
						gameTime, err := time.Parse(time.RFC3339, subValue.(string))
						if err != nil {
							continue
						}
						bowl.Games[ind].GameDate = gameTime
					} else if sub == "espnID" {
						bowl.Games[ind].EspnID = subValue.(string)
					}
				}
			}
		}
	}
	bowl.UpdateTime = time.Now()
}
