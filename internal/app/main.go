package main

import (
	//"context"
	//"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nathanhnew/bowl-backend/internal/app/auth"
	"github.com/nathanhnew/bowl-backend/internal/app/config"
	"github.com/nathanhnew/bowl-backend/internal/app/controllers"
	"log"
	"net/http"
	//"os"
)

func main() {
	r := mux.NewRouter()

	var cfg, _ = config.GetConfig(config.DefaultConfigLocation)

	port := cfg.GetListenPort()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("assets/"))))
	r.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	r.Handle("/users", auth.ValidToken(controllers.DeleteUserHandler)).Methods("DELETE")
	r.Handle("/users/{user}", auth.ValidToken(controllers.GetUserHandler)).Methods("GET")
	r.Handle("/users/{user}", auth.ValidToken(controllers.UpdateUserHandler)).Methods("PATCH")
	r.Handle("/users/{user}/leagues", auth.ValidToken(controllers.GetLeaguesByUserHandler)).Methods("GET")
	r.HandleFunc("/login", auth.Login).Methods("POST")

	fmt.Printf("Listening to port %d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
