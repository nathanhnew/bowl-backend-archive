package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nathanhnew/bowl-backend/internal/app/db"
	"net/http"
)

func getThemeFromSchool(w http.ResponseWriter, req *http.Request) {
	slug := mux.Vars(req)["school"]
	if theme, err := db.GetThemeFromSchool(slug); err != nil {
		// Have an error, break out
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		// No error, return theme
		_ = json.NewEncoder(w).Encode(theme)
	}
	return
}

var GetThemeFromSchoolHandler = http.HandlerFunc(getThemeFromSchool)
