package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nathanhnew/bowl-backend/internal/app/auth"
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

func createLeague(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	league := models.NewLeague()
	decodeErr := json.NewDecoder(req.Body).Decode(&league)
	if decodeErr != nil {
		http.Error(w, "Unable to read payload", http.StatusBadRequest)
		fmt.Printf("%s\n", decodeErr)
		return
	}
	league.Slug = params["league"]
	// Verify new league
	if validLeague, err := db.ValidateNewLeague(league.Slug); err != nil {
		http.Error(w, "Unable to verify new league", http.StatusInternalServerError)
		return
	} else if !validLeague {
		http.Error(w, "League already exists", http.StatusBadRequest)
		return
	}
	// Verify commissioner is making call
	authorized := req.Context().Value("auth").(auth.Claim)
	if authorized.User != league.Commissioner && authorized.IsAdmin == false {
		http.Error(w, "Cannot assign this commissioner", http.StatusForbidden)
		fmt.Printf("User %s attempted to create league with commissioner %s\n", authorized.User, league.Commissioner)
		return
	}
	err := db.CreateLeague(league)
	if err != nil {
		http.Error(w, "Unable to create league", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(league)
	return
}

func updateLeague(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	league, err := db.GetLeagueBySlug(params["league"])
	leagueName := league.Name
	if err != nil {
		http.Error(w, "Unable to get league from database", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	// Verify user can update this league
	authorized := req.Context().Value("auth").(auth.Claim)
	if authorized.User != league.Commissioner && authorized.IsAdmin == false {
		http.Error(w, "Cannot assign this commissioner", http.StatusForbidden)
		fmt.Printf("User %s attempted to create league with commissioner %s\n", authorized.User, league.Commissioner)
		return
	}
	err = json.NewDecoder(req.Body).Decode(&league)
	if err != nil {
		http.Error(w, "Unable to parse payload", http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}
	// If name changed, generate new slug, otherwise force slug to remain constant
	if league.Name != leagueName {
		newSlug, err := newSlug("league", league.Name)
		if err != nil {
			fmt.Println(err.Error())
		}
		if newSlug != params["league"] {
			league.Slug = newSlug
		} else {
			league.Slug = params["league"]
		}
	} else {
		league.Slug = params["league"]
	}
	db.UpdateLeague(params["league"], league)
	_ = json.NewEncoder(w).Encode(league)
	return
}

func deactivateLeague(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	league, err := db.GetLeagueBySlug(params["league"])
	if err != nil {
		http.Error(w, "Unable to verify league", http.StatusBadRequest)
		return
	}
	// Verify user can deactivate this league
	authorized := req.Context().Value("auth").(auth.Claim)
	if authorized.User != league.Commissioner && authorized.IsAdmin == false {
		http.Error(w, "Cannot assign this commissioner", http.StatusForbidden)
		fmt.Printf("User %s attempted to create league with commissioner %s\n", authorized.User, league.Commissioner)
		return
	}
	err = db.DeactivateLeague(params["league"])
	if err != nil {
		http.Error(w, "Unable to deactivate league", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getLeagueSlugFromName(w http.ResponseWriter, req *http.Request) {
	payload := make(map[string]string)
	_ = json.NewDecoder(req.Body).Decode(&payload)
	slug, err := newSlug("league", payload["name"])
	if err != nil {
		http.Error(w, "Unable to generate slug from name", http.StatusInternalServerError)
		return
	}
	payload["slug"] = slug

	_ = json.NewEncoder(w).Encode(payload)
	return
}

var GetLeaguesHandler = http.HandlerFunc(allLeagueHeadlines)
var CreateLeagueHandler = http.HandlerFunc(createLeague)
var DeactivateLeagueHandler = http.HandlerFunc(deactivateLeague)
var UpdateLeagueHandler = http.HandlerFunc(updateLeague)
var GetLeagueSlugFromNameHandler = http.HandlerFunc(getLeagueSlugFromName)
