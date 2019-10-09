package main

import (
	//"github.com/nathanhnew/bowl-backend/internal/app/controllers"
	"encoding/json"
	"fmt"
	"github.com/nathanhnew/bowl-backend/internal/app/db"
	"github.com/nathanhnew/bowl-backend/internal/models"
	//"github.com/gorilla/mux"
	"io/ioutil"
	//"log"
	//"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

func getConfig() (map[string]interface{}, error) {
	cfgFile, err := os.Open("config/main.conf.json")
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()
	byteCfg, _ := ioutil.ReadAll(cfgFile)

	var cfg map[string]interface{}
	json.Unmarshal([]byte(byteCfg), &cfg)

	return cfg, nil
}

func main() {
	//r := mux.NewRouter()
	cfg, _ := getConfig()
	fmt.Println(cfg)
	//port := int(cfg["port"].(float64))
	mongoUri := cfg["mongoUri"].(string)
	ctx := db.GetContext(30)
	client := db.Connect(mongoUri)
	var user models.User
	collection := client.Database("application").Collection("User")
	_ = collection.FindOne(ctx, bson.D{}).Decode(&user)
	fmt.Printf("%+v\n", user)
	fmt.Printf("Router initialized on port %d\n", int(cfg["port"].(float64)))
	//r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
	//	http.FileServer(http.Dir("../assets/"))))
	//r.HandleFunc("/teams", controllers.GetAllTeams).Methods("GET")
	//
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
