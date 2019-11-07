package models

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strings"
	"time"
)

type User struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name           userName           `json:"name" bson:"name"`
	Email          string             `json:"email" bson:"email"`
	Password       string             `json:"-" bson:"password"`
	FavoriteSchool string             `json:"favoriteSchool" bson:"favoriteSchool,omitempty"`
	Theme          ColorScheme        `json:"theme" bson:"theme,omitempty"`
	Icon           string             `json:"icon" bson:"icon,omitempty"`
	CreateTime     time.Time          `json:"-" bson:"createTime"`
	UpdateTime     time.Time          `json:"-" bson:"updateTime"`
	Token          json.Token         `json:"token" bson:"-"`
	Active         bool               `json:"active" bson:"active"`
	Admin          bool               `json:"-" bson:"admin"`
}

type userName struct {
	First  string `json:"first" bson:"first"`
	Last   string `json:"last" bson:"last"`
	Suffix string `json:"suffix,omitempty" bson:"suffix,omitempty"`
}

type ColorScheme struct {
	PrimaryColor   string `json:"primaryColor" bson:"primary,omitempty"`
	SecondaryColor string `json:"secondaryColor" bson:"secondary,omitempty"`
	TertiaryColor  string `json:"tertiaryColor" bson:"tertiary,omitempty"`
}

func NewUser() User {
	var user User
	user.CreateTime = time.Now()
	user.UpdateTime = time.Now()
	user.Active = true
	user.Admin = false
	return user
}

func (theme *ColorScheme) UpdateFromMap(payload map[string]interface{}) {
	for k, v := range payload {
		if match, _ := regexp.Match("[P|p]rimary", []byte(k)); match {
			theme.PrimaryColor = v.(string)
		}
		if match, _ := regexp.Match("[S|s]econdary", []byte(k)); match {
			theme.SecondaryColor = v.(string)
		}
		if match, _ := regexp.Match("[T|t]ertiary", []byte(k)); match {
			theme.TertiaryColor = v.(string)
		}
	}
}

func (user *User) UpdateFromMap(payload map[string]interface{}) {
	for k, v := range payload {
		if k == "email" {
			user.Email = strings.ToLower(v.(string))
		}
		if k == "password" {
			user.Password = v.(string)
		}
		if k == "firstName" {
			user.Name.First = strings.Title(v.(string))
		}
		if k == "lastName" {
			user.Name.Last = strings.Title(v.(string))
		}
		if k == "suffix" {
			user.Name.Suffix = strings.Title(v.(string))
		}
		if k == "favoriteSchool" {
			user.FavoriteSchool = v.(string)
		}
		if k == "primaryColor" {
			user.Theme.PrimaryColor = v.(string)
		}
		if k == "secondaryColor" {
			user.Theme.SecondaryColor = v.(string)
		}
		if k == "tertiaryColor" {
			user.Theme.TertiaryColor = v.(string)
		}
		if k == "icon" {
			user.Icon = v.(string)
		}
		if k == "active" {
			user.Active = v.(bool)
		}
		user.UpdateTime = time.Now()
	}
}
