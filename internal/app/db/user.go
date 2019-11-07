package db

import (
	"github.com/nathanhnew/bowl-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(user models.User) (*mongo.InsertOneResult, error) {
	ctx = getContext()
	result, err := Client.Database(database).Collection("User").InsertOne(ctx, user)

	return result, err
}

func ValidateNewEmail(email string) (bool, error) {
	ctx = getContext()
	cursor, err := Client.Database(database).Collection(userCollection).Find(ctx, bson.D{{"email", email}, {"active", true}})
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		return false, nil
	}
	return true, nil
}

func DeleteUser(email string) error {
	ctx = getContext()
	_, err := Client.Database(database).Collection(userCollection).DeleteOne(ctx, bson.M{"email": email})
	return err
}

func GetUser(email string) (models.User, error) {
	ctx = getContext()
	var user models.User
	err := Client.Database(database).Collection(userCollection).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return user, err
}

func UpdateUser(email string, payload map[string]interface{}) (models.User, error) {
	ctx = getContext()
	var user models.User
	err := Client.Database(database).Collection(userCollection).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return user, err
	}
	user.UpdateFromMap(payload)
	_ = Client.Database(database).Collection(userCollection).FindOneAndReplace(ctx, bson.M{"email": email}, user)
	return user, err
}
