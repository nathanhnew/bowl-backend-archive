package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           userName           `bson:"name"`
	Email          string             `bson:"email"`
	Password       string             `bson:"password"`
	FavoriteSchool primitive.ObjectID `bson:"favoriteSchool,omitempty"`
	Theme          colorScheme        `bson:"theme,omitempty"`
	Icon           string             `bson:"icon,omitempty"`
	CreateTime     time.Time          `bson:"createTime"`
	UpdateTime     time.Time          `bson:"updateTime"`
}

type userName struct {
	First  string `bson:"first"`
	Last   string `bson:"last"`
	Suffix string `bson:"suffix,omitempty"`
}

type colorScheme struct {
	PrimaryColor   string `bson:"primary,omitempty"`
	SecondaryColor string `bson:"secondary,omitempty"`
	TertiaryColor  string `bson:"tertiary,omitempty"`
}
