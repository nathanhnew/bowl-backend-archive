package main

import (
	//"github.com/nathanhnew/bowl-backend/internal/app/controllers"
	//"encoding/json"
	"fmt"
	"github.com/nathanhnew/bowl-backend/internal/app/config"
	"github.com/nathanhnew/bowl-backend/internal/app/db"
	"github.com/nathanhnew/bowl-backend/internal/models"
	//"github.com/gorilla/mux"
	//"io/ioutil"
	//"log"
	//"net/http"
	"go.mongodb.org/mongo-driver/bson"
	//"os"
)

func main() {
	//r := mux.NewRouter()
	cfg, err := config.GetConfig(config.DefaultConfigLocation)
	if err != nil {
		fmt.Println(err)
	}
	//port := int(cfg["port"].(float64))
	mongoUri := cfg.GetMongoUri()
	ctx := db.GetContext(30)
	client := db.Connect(mongoUri)
	defer db.Disconnect(client)
	var user models.User
	collection := client.Database("application").Collection("User")
	_ = collection.FindOne(ctx, bson.D{}).Decode(&user)
	fmt.Printf("%+v\n", user)
	fmt.Printf("Router initialized on port %d\n", cfg.GetListenPort())
	//r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
	//	http.FileServer(http.Dir("../assets/"))))
	//r.HandleFunc("/teams", controllers.GetAllTeams).Methods("GET")
	//
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
