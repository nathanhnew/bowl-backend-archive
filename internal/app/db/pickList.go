package db

import "github.com/nathanhnew/bowl-backend/internal/models"

func CreatePicks(picks models.PickList) error {
	ctx = getContext()
	_, err := Client.Database(database).Collection(pickCollection).InsertOne(ctx, picks)
	return err
}
