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
	r.HandleFunc("/login", auth.Login).Methods("POST")
	r.Handle("/users", auth.ValidToken(controllers.DeleteUserHandler)).Methods("DELETE")
	r.Handle("/users/{user}", auth.ValidToken(controllers.GetUserHandler)).Methods("GET")
	r.Handle("/users/{user}", controllers.CreateUserHandler).Methods("POST")
	r.Handle("/users/{user}", auth.ValidToken(controllers.UpdateUserHandler)).Methods("PATCH")
	r.Handle("/users/{user}/leagues", auth.ValidToken(controllers.GetLeaguesByUserHandler)).Methods("GET")

	r.Handle("/leagues", auth.ValidToken(controllers.GetLeaguesHandler)).Methods("GET")
	r.Handle("/leagues/slug", auth.ValidToken(controllers.GetLeagueSlugFromNameHandler)).Methods("POST")
	r.Handle("/leagues/{league}", auth.ValidToken(controllers.CreateLeagueHandler)).Methods("POST")
	r.Handle("/leagues/{league}", auth.ValidToken(controllers.DeactivateLeagueHandler)).Methods("DELETE")
	r.Handle("/leagues/{league}", auth.ValidToken(controllers.UpdateLeagueHandler)).Methods("PATCH")

	r.Handle("/bowls", auth.ValidToken(controllers.GetAllBowlsHandler)).Methods("GET")
	r.Handle("/bowls/slug", auth.ValidToken(controllers.GetBowlSlugFromNameHandler)).Methods("POST")
	r.Handle("/bowls/{bowl}", auth.ValidToken(controllers.GetBowlHandler)).Methods("GET")
	r.Handle("/bowls/{bowl}", auth.ValidToken(controllers.CreateBowlHandler)).Methods("POST")

	r.Handle("/schools/{school}/theme", controllers.GetThemeFromSchoolHandler).Methods("GET")
	fmt.Printf("Listening to port %d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
