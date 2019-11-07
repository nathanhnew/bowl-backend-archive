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

func createBowl(w http.ResponseWriter, req *http.Request) {
	bowl := models.NewBowl()
	authorization := req.Context().Value("auth").(auth.Claim)
	if authorization.IsAdmin == false {
		http.Error(w, "Cannot update bowls. Please contact admin", http.StatusForbidden)
		fmt.Printf("User %s attempted to create bowl\n", authorization.User)
		return
	}
	payload := make(map[string]interface{})
	decodeErr := json.NewDecoder(req.Body).Decode(&payload)
	if decodeErr != nil {
		http.Error(w, "Unable to parse payload", http.StatusBadRequest)
	}
	bowl.Slug = payload["slug"].(string)
	bowlExists, err := db.BowlExists(bowl.Slug)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to verify bowl %s", bowl.Slug), http.StatusInternalServerError)
		return
	}
	if bowlExists {
		http.Error(w, fmt.Sprintf("Bowl '%s' already exists", bowl.Slug), http.StatusBadRequest)
		return
	}

	bowl.UpdateFromMap(payload)

	err = db.CreateBowl(bowl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to create bowl %s", bowl.Slug), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(bowl)
	return
}

func getAllBowls(w http.ResponseWriter, req *http.Request) {
	var bowls []models.BowlHeadline

	bowls, err := db.GetAllBowlHeadlines()
	if err != nil {
		http.Error(w, "Unable to get bowls", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(bowls)
	return
}

func getBowl(w http.ResponseWriter, req *http.Request) {
	var bowl models.Bowl
	params := mux.Vars(req)
	bowl, err := db.GetBowlBySlug(params["bowl"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to get bowl %s", params["bowl"]), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(bowl)
	return
}

func getBowlSlugFromName(w http.ResponseWriter, req *http.Request) {
	payload := make(map[string]string)
	_ = json.NewDecoder(req.Body).Decode(&payload)
	slug, err := newSlug("bowl", payload["name"])
	if err != nil {
		http.Error(w, "Unable to generate slug from name", http.StatusInternalServerError)
		return
	}
	payload["slug"] = slug

	_ = json.NewEncoder(w).Encode(payload)
	return
}

var GetAllBowlsHandler = http.HandlerFunc(getAllBowls)
var GetBowlHandler = http.HandlerFunc(getBowl)
var CreateBowlHandler = http.HandlerFunc(createBowl)
var GetBowlSlugFromNameHandler = http.HandlerFunc(getBowlSlugFromName)
