package db

import (
	"fmt"
	"github.com/nathanhnew/bowl-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

func GetSchoolBySlug(slug string) (models.School, error) {
	ctx = getContext()
	var school models.School
	Client.Database(database).Collection("School").FindOne(ctx, bson.D{{"name.slug", slug}}).Decode(&school)
	if reflect.DeepEqual(school, models.School{}) {
		return school, fmt.Errorf("unable to find school %s", slug)
	}
	return school, nil
}
