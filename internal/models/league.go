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
	CreateTime   time.Time          `json:"createTime" bson:"createTime"`
	UpdateTime   time.Time          `json:"updateTime" bson:"updateTime"`
	Active       bool               `json:"active" bson:"active"`
}

type LeagueHeadline struct {
	Name           string    `json:"name" bson:"name"`
	Slug           string    `json:"slug" bson:"slug"`
	Commissioner   string    `json:"commissioner" bson:"commissioner"`
	PreviousWinner string    `json:"previousWinner" bson:"-"`
	NextBowl       string    `json:"nextBowl" bson:"-"`
	NextBowlTime   time.Time `json:"nextBowlTime" bson:"-"`
}

type season struct {
	Season             string       `json:"season" bson:"season"`
	Players            []player     `json:"players" bson:"players,omitempty"`
	Bowls              []string     `json:"bowls" bson:"bowls,omitempty"`
	SubmissionLockDate time.Time    `json:"submissionLockDate" bson:"submissionLockDate"`
	Spread             spreadConfig `json:"spread" bson:"spread"`
	Dues               duesConfig   `json:"dues" bson:"dues"`
	Winner             string       `json:"winner" bson:"winner"`
}

type spreadConfig struct {
	UseSpread bool      `json:"useSpread" bson:"useSpread"`
	Simplify  string    `json:"simplify,omitempty" bson:"simplify,omitempty"`
	LockDate  time.Time `json:"lockDate,omitempty" bson:"lockDate,omitempty"`
}

type duesConfig struct {
	HasDues   bool    `json:"hasDues" bson:"hasDues"`
	Amount    float64 `json:"amount,omitempty" bson:"amount,omitempty"`
	UseOnline bool    `json:"useOnline,omitempty" bson:"useOnline,omitempty"`
	Payout    payout  `json:"payout,omitempty" bson:"payout,omitempty,omitempty"`
}

type payout struct {
	Type    string `json:"type" bson:"type"` // fixed - fixed dollar value, percentage - pct of total money
	Amounts []struct {
		Rank   float64     `json:"rank" bson:"rank"`
		Amount interface{} `json:"amount" bson:"amount"`
	} `json:"amounts" bson:"amounts"`
}

type player struct {
	User           string    `json:"user" bson:"user"`
	Status         string    `json:"status" bson:"status"`
	HasSubmitted   bool      `json:"hasSubmitted" bson:"hasSubmitted"`
	SubmissionDate time.Time `json:"submissionDate,omitempty" bson:"submissionDate,omitempty"`
	HasPaid        bool      `json:"hasPaid,omitempty" bson:"hasPaid,omitempty"`
}

func NewLeague() League {
	league := League{}
	league.ID = primitive.NewObjectID()
	league.CreateTime = time.Now()
	league.UpdateTime = time.Now()
	league.Active = true
	return league
}

func (league *League) UpdateFromMap(payload map[string]interface{}) {
	for k, v := range payload {
		if k == "name" {
			league.Name = v.(string)
		}
		if k == "slug" {
			league.Slug = v.(string)
		}
		if k == "commissioner" {
			league.Commissioner = v.(string)
		}
		if k == "active" {
			league.Active = v.(bool)
		}
	}
}

func (season *season) UpdateFromMap(payload map[string]interface{}) {
	for k, v := range payload {
		if k == "season" {
			season.Season = v.(string)
		}
		if k == "submissionLockDate" {
			season.SubmissionLockDate = v.(time.Time)
		}
		if k == "spread" {

		}
	}
}
