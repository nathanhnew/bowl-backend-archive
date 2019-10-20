package db

import (
	"github.com/nathanhnew/bowl-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetLeaguesByUser(email string) ([]models.League, error) {
	ctx = getContext()
	var leagues []models.League
	cursor, err := Client.Database(database).Collection(leagueCollection).Find(ctx, bson.M{"seasons.players.user": email})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var league models.League
		cursor.Decode(&league)
		leagues = append(leagues, league)
	}
	return leagues, nil
}
