package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/nathanhnew/bowl-backend/internal/app/auth"
	"github.com/nathanhnew/bowl-backend/internal/app/db"
	"github.com/nathanhnew/bowl-backend/internal/models"
	"net/http"
)

func createPicks(w http.ResponseWriter, req *http.Request) {
	picks := models.NewPickList()
	decodeErr := json.NewDecoder(req.Body).Decode(&picks)
	if decodeErr != nil {
		http.Error(w, "Unable to read payload", http.StatusBadRequest)
		fmt.Printf("%s\n", decodeErr)
		return
	}
	leagueAdmin, err := db.GetLeagueCommissioner(picks.League)
	if err != nil {
		http.Error(w, "Unalbe to validate league", http.StatusBadRequest)
		fmt.Printf("%s\n", err)
	}
	authorized := req.Context().Value("auth").(auth.Claim)
	if authorized.User != picks.User && authorized.IsAdmin == false && authorized.User != leagueAdmin {
		http.Error(w, "Cannot create picks for this user", http.StatusForbidden)
		fmt.Printf("User %s attempted to create picks for uesr %s\n", authorized.User, picks.User)
		return
	}

	err = db.CreatePicks(picks)
	if err != nil {
		http.Error(w, "Unable to create picks", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.WriteHeader(200)
	return
}

var CreatePicksHandler = http.HandlerFunc(createPicks)
