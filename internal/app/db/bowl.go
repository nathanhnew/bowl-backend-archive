package db

import (
	"errors"
	"fmt"
	"github.com/nathanhnew/bowl-backend/internal/app/config"
	"github.com/nathanhnew/bowl-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
)

func GetBowlSlugByBase(slug string) (string, error) {
	var returnValue string
	ctx = getContext()
	cursor, err := Client.Database(database).Collection(bowlCollection).Aggregate(ctx, bson.A{
		bson.M{"$match": bson.M{"slug": bson.M{"$regex": slug}}},
		bson.M{"$sort": bson.M{"slug": -1}},
		bson.M{"$limit": 1},
		bson.M{"$project": bson.M{"_id": 0, "slug": 1}},
	})
	if err != nil {
		return "", err
	}
	for cursor.Next(ctx) {
		var queryReturn map[string]interface{}
		cursor.Decode(&queryReturn)
		returnValue = queryReturn["slug"].(string)
	}
	return returnValue, nil
}

func getNextBowlFromList(bowls []string, season string, bowlChan chan map[string]interface{}) {
	ctx = getContext()
	payload := make(map[string]interface{})
	output := make(map[string]interface{})
	if bowls == nil {
		bowlChan <- nil
		return
	}
	cursor, err := Client.Database(database).Collection(bowlCollection).Aggregate(ctx, bson.A{
		bson.M{"$match": bson.M{"slug": bson.M{"$in": bowls}}},
		bson.M{"$unwind": "$games"},
		bson.M{"$match": bson.M{"games.season": season}},
		bson.M{"$match": bson.M{"games.gameDate": bson.M{"$gte": time.Now()}}},
		bson.M{"$sort": bson.M{"games.gameDate": 1}},
		bson.M{"$project": bson.M{"name": 1, "games.gameDate": 1, "_id": 0}},
		bson.M{"$limit": 1},
	})
	if err != nil {
		fmt.Printf("%+v", err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		cursor.Decode(&payload)
		output["nextBowl"] = payload["name"].(string)
		output["nextBowlTime"] = payload["games"].(map[string]interface{})["gameDate"].(primitive.DateTime).Time()
		bowlChan <- output
	}
}

func getBowlsFromLeague(league string, season string) []string {
	ctx = getContext()
	var bowls []string
	payload := make(map[string]interface{})
	cursor, err := Client.Database(database).Collection(leagueCollection).Aggregate(ctx, bson.A{
		bson.M{"$match": bson.M{"slug": league}},
		bson.M{"$unwind": "$seasons"},
		bson.M{"$match": bson.M{"seasons.season": season}},
		bson.M{"$project": bson.M{"_id": 0, "seasons.bowls": 1}},
	})
	if err != nil {
		fmt.Printf("%+v", err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		cursor.Decode(&payload)
		for _, bwl := range payload["seasons"].(map[string]interface{})["bowls"].(bson.A) {
			bowls = append(bowls, bwl.(string))
		}
		return bowls
	}
	return nil
}

func GetAllBowlHeadlines() ([]models.BowlHeadline, error) {
	var bowls []models.BowlHeadline
	ctx = getContext()
	cursor, err := Client.Database(database).Collection(bowlCollection).Aggregate(ctx, bson.A{
		bson.M{"$unwind": "$games"},
		bson.M{"$match": bson.M{"games.season": config.GetCurrentSeason()}},
		bson.M{"$unwind": "$games.spreads"},
		bson.M{"$sort": bson.M{"games.gameDate": -1, "games.spreads.queryDate": -1}},
		bson.M{"$group": bson.M{"_id": "$slug", "info": bson.M{"$first": "$$ROOT"}}},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$info"}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var bowl models.BowlHeadline
		cursor.Decode(&bowl)
		bowls = append(bowls, bowl)
	}
	return bowls, nil
}

func GetBowlBySlug(slug string) (models.Bowl, error) {
	ctx = getContext()
	var bowl models.Bowl
	err := Client.Database(database).Collection(bowlCollection).FindOne(ctx, bson.M{"slug": slug}).Decode(&bowl)
	if reflect.DeepEqual(bowl, models.Bowl{}) {
		// Couldn't find bowl in DB
		return bowl, errors.New(fmt.Sprintf("cannot find bowl %s", slug))
	}
	return bowl, err
}

func CreateBowl(bowl models.Bowl) error {
	ctx = getContext()
	_, err := Client.Database(database).Collection(bowlCollection).InsertOne(ctx, bowl)
	return err
}

func BowlExists(slug string) (bool, error) {
	ctx = getContext()
	cursor, err := Client.Database(database).Collection(bowlCollection).Find(ctx, bson.M{"slug": slug})
	if err != nil {
		return true, err
	}
	for cursor.Next(ctx) {
		// Exists in DB return true
		return true, nil
	}
	return false, nil
}
