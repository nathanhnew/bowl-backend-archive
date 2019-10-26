package db

import (
	"encoding/json"
	"fmt"
	"github.com/nathanhnew/bowl-backend/internal/app/config"
	"github.com/nathanhnew/bowl-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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

	for cursor.Next(ctx) {
		cursor.Decode(&payload)
		var lastWinner = payload["seasons"].(map[string]interface{})["winner"].(string)
		winnerChan <- lastWinner
	}
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
	for cursor.Next(ctx) {
		cursor.Decode(&payload)
		for _, bwl := range payload["seasons"].(map[string]interface{})["bowls"].(bson.A) {
			bowls = append(bowls, bwl.(string))
		}
		return bowls
	}
	return nil
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
