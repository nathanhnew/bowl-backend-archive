package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nathanhnew/bowl-backend/internal/app/config"
	"github.com/nathanhnew/bowl-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

func GetLeaguesByUser(email string) ([]models.League, error) {
	ctx = getContext()
	var leagues []models.League
	cursor, err := Client.Database(database).Collection(leagueCollection).Find(ctx, bson.M{"seasons.players.user": email})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var league models.League
		cursor.Decode(&league)
		leagues = append(leagues, league)
	}
	return leagues, nil
}

func getWinnerByLeague(league string, winnerChan chan string) {
	ctx = getContext()
	payload := make(map[string]interface{})
	cursor, err := Client.Database(database).Collection(leagueCollection).Aggregate(ctx, bson.A{
		bson.M{"$match": bson.M{"slug": league}},
		bson.M{"$unwind": "$seasons"},
		bson.M{"$match": bson.M{"seasons.winner": bson.M{"$exists": true}}},
		bson.M{"$sort": bson.M{"seasons.season": -1}},
		bson.M{"$limit": 1},
		bson.M{"$project": bson.M{"seasons.winner": 1, "_id": 0}},
	})
	if err != nil {
		fmt.Printf("%+v", err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		cursor.Decode(&payload)
		var lastWinner = payload["seasons"].(map[string]interface{})["winner"].(string)
		winnerChan <- lastWinner
	}
}

func GetLeagueBySlug(slug string) (models.League, error) {
	ctx = getContext()
	var league models.League
	err := Client.Database(database).Collection(leagueCollection).FindOne(ctx, bson.M{"slug": slug}).Decode(&league)
	if reflect.DeepEqual(league, models.League{}) {
		// Couldn't find league in DB
		return league, errors.New(fmt.Sprintf("cannot find league %s", slug))
	}
	return league, err
}

func DeactivateLeague(slug string) error {
	ctx = getContext()
	_, err := Client.Database(database).Collection(leagueCollection).UpdateOne(ctx,
		bson.M{"slug": slug},
		bson.M{"$set": bson.M{"active": false}})
	return err
}

func UpdateLeague(slug string, league models.League) {
	ctx = getContext()
	_ = Client.Database(database).Collection(leagueCollection).FindOneAndReplace(ctx, bson.M{"slug": slug}, league)
}

func GetAllLeagueHeadlines(start int64, limit int64) ([]models.LeagueHeadline, error) {
	var leagues []models.LeagueHeadline
	ctx = getContext()
	var queryOptions options.FindOptions
	queryOptions.SetSkip(start)
	queryOptions.SetLimit(limit)
	cursor, err := Client.Database(database).Collection(leagueCollection).Find(ctx, bson.M{}, &queryOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var league models.LeagueHeadline
		cursor.Decode(&league)

		winnerChan := make(chan string, 1)
		nextGameChan := make(chan map[string]interface{}, 1)

		go getWinnerByLeague(league.Slug, winnerChan)
		go getNextBowlFromList(getBowlsFromLeague(league.Slug, config.GetCurrentSeason()), config.GetCurrentSeason(), nextGameChan)

		league.PreviousWinner = <-winnerChan
		nextGame := <-nextGameChan
		nextGameString, err := json.Marshal(nextGame)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(nextGameString), &league)
		if err != nil {
			return nil, err
		}

		leagues = append(leagues, league)
	}
	return leagues, nil
}

func GetLeagueSlugByBase(slug string) (string, error) {
	var returnValue string
	ctx = getContext()
	cursor, err := Client.Database(database).Collection(leagueCollection).Aggregate(ctx, bson.A{
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

func ValidateNewLeague(slug string) (bool, error) {
	ctx = getContext()
	cursor, err := Client.Database(database).Collection(leagueCollection).Find(ctx, bson.M{"slug": slug})
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		return false, nil
	}
	return true, nil
}

func GetLeagueCommissioner(slug string) (string, error) {
	ret := make(map[string]interface{})
	ctx = getContext()
	cursor, err := Client.Database(database).Collection(leagueCollection).Aggregate(ctx, bson.A{
		bson.M{"$match": bson.M{"slug": slug}},
		bson.M{"$limit": 1},
		bson.M{"$project": bson.M{"commissioner": 1}},
	})
	if err != nil {
		return "", err
	}
	for cursor.Next(ctx) {
		cursor.Decode(&ret)
		return ret["commissioner"].(string), nil
	}
	return "", errors.New("league not found")
}

func CreateLeague(league models.League) error {
	ctx = getContext()
	_, err := Client.Database(database).Collection(leagueCollection).InsertOne(ctx, league)
	return err
}
