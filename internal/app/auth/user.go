package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nathanhnew/bowl-backend/internal/app/db"
	"github.com/nathanhnew/bowl-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
)

func CreateUser(w http.ResponseWriter, req *http.Request) {
	user := models.NewUser()
	var payload map[string]interface{}
	decodeErr := json.NewDecoder(req.Body).Decode(&payload)
	if decodeErr != nil {
		http.Error(w, "Unable to read payload", 400)
		fmt.Printf("%s\n", decodeErr)
		return
	}
	if payload["email"].(string) == "" || payload["firstName"].(string) == "" || payload["lastName"].(string) == "" || payload["password"].(string) == "" {
		http.Error(w, "Missing required fields", 400)
		return
	}
	emailExists, err := db.NonExistentEmail(payload["email"].(string))
	if err != nil {
		fmt.Printf("email: %s\n", err)
		http.Error(w, "Unable to verify email address", 400)
		return
	}
	if emailExists == true {
		http.Error(w, "Email already exists.", 400)
		fmt.Printf("Email address %s attemped re-register\n", payload["email"].(string))
		return
	}

	// Set Required fields
	user.Email = strings.ToLower(payload["email"].(string))
	user.Name.First = strings.Title(payload["firstName"].(string))
	user.Name.Last = strings.Title(payload["lastName"].(string))
	if suffix, ok := payload["suffix"].(string); ok {
		user.Name.Suffix = strings.Title(suffix)
	}

	// Set Password
	pwd, hashErr := bcrypt.GenerateFromPassword([]byte(payload["password"].(string)), bcrypt.DefaultCost)
	if hashErr != nil {
		http.Error(w, "Unable to create password hash", 500)
		fmt.Printf("Unable to create password hash\n%s\n", hashErr)
		return
	}
	user.Password = string(pwd)

	favoriteSchool, err := db.GetSchoolBySlug(payload["favoriteSchool"].(string))
	if err != nil {
		http.Error(w, "School not found", 400)
		fmt.Printf("Unable to find favorite school %s: %s\n", user.Email, payload["favoriteSchool"].(string))
		return
	}

	user.FavoriteSchool = favoriteSchool.ID

	if theme, ok := payload["theme"]; !ok {
		user.Theme.PrimaryColor = favoriteSchool.Colors.PrimaryColor
		user.Theme.SecondaryColor = favoriteSchool.Colors.SecondaryColor
		user.Theme.TertiaryColor = favoriteSchool.Colors.TertiaryColor
	} else {
		user.Theme.PrimaryColor = theme.(map[string]string)["primary"]
		user.Theme.PrimaryColor = theme.(map[string]string)["secondary"]
		if tertiary, hasThree := theme.(map[string]string)["tertiary"]; hasThree {
			user.Theme.TertiaryColor = tertiary
		}
	}

	if icon, ok := payload["icon"].(string); ok {
		user.Icon = icon
	} else {
		files, err := ioutil.ReadDir("./assets/img/icon/default")
		if err == nil {
			i := rand.Intn(len(files))
			user.Icon = files[i].Name()
		}
	}

	id, err := db.CreateUser(user)
	if err != nil {
		http.Error(w, "Cannot add user", 500)
		fmt.Printf("Unable to add user %s\n%s\n", user.Email, err)
		return
	}
	user.ID = id.InsertedID.(primitive.ObjectID)

	fmt.Printf("Created user %s\n", user.Email)

	token, _ := refreshToken(user.Email, user.Admin)

	user.Token = token

	_ = json.NewEncoder(w).Encode(user)

	return
}

func deactivateUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var user string = params["user"]
	reqEmail := req.Header.Get("authUser")
	reqAdmin := req.Header.Get("authAdmin")
	if reqEmail != user && reqAdmin == "false" {
		http.Error(w, "Cannot deactivate this account", http.StatusForbidden)
		fmt.Printf("User %s attempted to deactivate account %s\n", reqEmail, user)
		return
	}
	err := db.DeactivateUser(user)
	if err != nil {
		http.Error(w, "Cannot deactivate user", http.StatusInternalServerError)
		fmt.Printf("Unable to deactivate user %s\n%s\n", user, err)
		return
	}
	fmt.Printf("Deactivated user %s\n", user)
	return
}

func deleteUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var user string = params["user"]
	tokenUser := req.Header.Get("authUser")
	tokenAdmin := req.Header.Get("authAdmin")
	if tokenUser != user && tokenAdmin == "false" {
		http.Error(w, "Cannot delete this account", http.StatusForbidden)
		fmt.Printf("User %s attempted to delete account %s\n", tokenUser, user)
		return
	}
	err := db.DeleteUser(user)
	if err != nil {
		http.Error(w, "Cannot delete user", http.StatusInternalServerError)
		fmt.Printf("Unable to delete user %s\n%s\n", user, err)
		return
	}
	fmt.Printf("Deleted user %s\n", user)
	w.WriteHeader(http.StatusOK)
	return
}

func updateUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var email = params["user"]
	var payload map[string]interface{}
	decodeErr := json.NewDecoder(req.Body).Decode(&payload)
	if decodeErr != nil {
		http.Error(w, "Unable to read payload", 400)
		return
	}
	reqEmail := req.Header.Get("authUser")
	reqAdmin := req.Header.Get("authAdmin")
	if reqEmail != email && reqAdmin == "false" {
		http.Error(w, "Cannot update this account", http.StatusForbidden)
		fmt.Printf("User %s attempted to update account %s\n", reqEmail, email)
		return
	}
	if pwd, ok := params["password"]; ok {
		pwd, hashErr := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if hashErr != nil {
			http.Error(w, "Unable to create password hash", 500)
			fmt.Printf("Unable to create password hash\n%s\n", hashErr)
			return
		}
		params["password"] = string(pwd)
	}
	if schoolSlug, ok := params["favoriteSchool"]; ok {
		school, err := db.GetSchoolBySlug(schoolSlug)
		if err != nil {
			http.Error(w, "Unable to get school", http.StatusInternalServerError)
			fmt.Printf("Unable to get school %s\n%s\n", schoolSlug, err)
			return
		}
		payload["favoriteSchool"] = school.ID
	}
	user, err := db.UpdateUser(email, payload)
	if err != nil {
		http.Error(w, "Cannot update user", http.StatusInternalServerError)
		fmt.Printf("Unable to update user %s\n%s\n", email, err)
		return
	}
	_ = json.NewEncoder(w).Encode(user)
	fmt.Printf("Updated user %s\n", email)
	return
}

func getUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var email string = params["user"]
	user, err := db.GetUser(email)
	if err != nil {
		http.Error(w, "Cannot get user", http.StatusBadRequest)
		fmt.Printf("Unable to get user %s\n%s\n", email, err)
		return
	}
	_ = json.NewEncoder(w).Encode(user)
	return
}

var DeleteUserHandler = http.HandlerFunc(deleteUser)
var UpdateUserHandler = http.HandlerFunc(updateUser)
var GetUserHandler = http.HandlerFunc(getUser)
