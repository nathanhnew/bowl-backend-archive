package controllers

import (
	"encoding/json"
	"net/http"
)

func GetAllTeams(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode("Success! Very Nice")
}
