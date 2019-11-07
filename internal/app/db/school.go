package db

import (
	"fmt"
	"github.com/nathanhnew/bowl-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

func GetSchool(slug string) (models.School, error) {
	ctx = getContext()
	var school models.School
	Client.Database(database).Collection(schoolCollection).FindOne(ctx, bson.D{{"name.slug", slug}}).Decode(&school)
	if reflect.DeepEqual(school, models.School{}) {
		return school, fmt.Errorf("unable to find school %s", slug)
	}
	return school, nil
}

func GetThemeFromSchool(slug string) (models.ColorScheme, error) {
	ctx = getContext()
	var theme models.ColorScheme
	cursor, err := Client.Database(database).Collection(schoolCollection).Aggregate(ctx, bson.A{
		bson.M{"$match": bson.M{"name.slug": slug}},
		bson.M{"$project": bson.M{"_id": 0, "colors": 1}},
	})
	if err != nil {
		return theme, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		payload := make(map[string]interface{})
		cursor.Decode(&payload)
		theme.UpdateFromMap(payload["colors"].(map[string]interface{}))
	}
	return theme, nil
}
