package db

import (
	"fmt"
	"github.com/nathanhnew/bowl-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(user models.User) (*mongo.InsertOneResult, error) {
	ctx = getContext()
	result, err := Client.Database(database).Collection("User").InsertOne(ctx, user)

	return result, err
}

func NonExistentEmail(email string) (bool, error) {
	ctx = getContext()
	c, err := Client.Database(database).Collection(userCollection).Find(ctx, bson.D{{"email", email}, {"active", true}})
	if err != nil {
		return false, err
	}
	for c.Next(ctx) {
		return true, nil
	}
	return false, nil
}

func DeactivateUser(email string) error {
	ctx = getContext()
	_, err := Client.Database(database).Collection(userCollection).UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"active": false}})
	return err
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
	fmt.Printf("+%v\n\n\n", user)
	_ = Client.Database(database).Collection(userCollection).FindOneAndReplace(ctx, bson.M{"email": email}, user)
	return user, err
}