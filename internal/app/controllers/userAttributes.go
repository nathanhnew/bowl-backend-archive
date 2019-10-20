package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nathanhnew/bowl-backend/internal/app/db"
	"github.com/nathanhnew/bowl-backend/internal/models"
	"net/http"
)

func getLeaguesByUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var email = params["user"]
	var leagueList []models.League
	leagueList, err := db.GetLeaguesByUser(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_ = json.NewEncoder(w).Encode(leagueList)
}

var GetLeaguesByUserHandler = http.HandlerFunc(getLeaguesByUser)
