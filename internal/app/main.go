package main

import (
	//"context"
	//"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nathanhnew/bowl-backend/internal/app/auth"
	"github.com/nathanhnew/bowl-backend/internal/app/config"
	"log"
	"net/http"
	//"os"
)

func main() {
	r := mux.NewRouter()

	var cfg, _ = config.GetConfig(config.DefaultConfigLocation)

	port := int(cfg.Values["port"].(float64))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("../assets/"))))
	r.HandleFunc("/users", auth.CreateUser).Methods("POST")
	r.Handle("/users", auth.ValidToken(auth.DeleteUserHandler)).Methods("DELETE")
	r.Handle("/users/{user}", auth.ValidToken(auth.GetUserHandler)).Methods("GET")
	r.Handle("/users/{user}", auth.ValidToken(auth.UpdateUserHandler)).Methods("PATCH")
	r.HandleFunc("/login", auth.Login).Methods("POST")

	fmt.Printf("Listening to port %d\n", port)

	//log.Fatal(http.ListenAndServe("localhost:3000", r))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
