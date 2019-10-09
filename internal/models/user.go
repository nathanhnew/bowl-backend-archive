package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	Name           userName           `bson: "name"`
	Email          string             `bson: "email"`
	Password       string             `bson: "password"`
	FavoriteSchool primitive.ObjectID `bson: "favoriteSchool"`
	Theme          colorScheme        `bson: "theme"`
	Icon           string             `bson: "icon"`
}

type userName struct {
	First  string `bson: "first"`
	Last   string `bson: "last"`
	Suffix string `bson: "suffix"`
}

type colorScheme struct {
	PrimaryColor   string `bson: "primary"`
	SecondaryColor string `bson: "secondary"`
	TertiaryColor  string `bson: "tertiary"`
}
