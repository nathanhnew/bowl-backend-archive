package controllers

import (
	"encoding/json"
	"github.com/nathanhnew/bowl-backend/internal/app/db"
	"github.com/nathanhnew/bowl-backend/internal/models"
	"net/http"
)

func allLeagueHeadlines(w http.ResponseWriter, req *http.Request) {
	var leagues []models.LeagueHeadline
	var payload map[string]interface{}
	decodeErr := json.NewDecoder(req.Body).Decode(&payload)
	if decodeErr != nil {
		http.Error(w, "Unable to parse payload", http.StatusBadRequest)
		return
	}

	leagues, err := db.GetAllLeagueHeadlines(int64(payload["start"].(float64)), int64(payload["limit"].(float64)))
	if err != nil {
		http.Error(w, "Unable to get leagues", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(leagues)
}

var GetLeaguesHandler = http.HandlerFunc(allLeagueHeadlines)
