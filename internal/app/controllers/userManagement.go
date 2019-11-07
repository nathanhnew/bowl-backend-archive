package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nathanhnew/bowl-backend/internal/app/auth"
	"github.com/nathanhnew/bowl-backend/internal/app/db"
	"github.com/nathanhnew/bowl-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

func ValidateNewUserPayload(payload map[string]interface{}, w http.ResponseWriter) bool {
	if payload["email"].(string) == "" || payload["firstName"].(string) == "" || payload["lastName"].(string) == "" || payload["password"].(string) == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return false
	}
	if !auth.ValidEmail(payload["email"].(string)) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return false
	}
	validEmail, err := db.ValidateNewEmail(payload["email"].(string))
	if err != nil {
		fmt.Printf("email: %s\n", err)
		http.Error(w, "Unable to verify email address", http.StatusInternalServerError)
		return false
	}
	if validEmail == false {
		http.Error(w, "Email already exists.", http.StatusBadRequest)
		fmt.Printf("Email address %s attemped re-register\n", payload["email"].(string))
		return false
	}
	return true
}

func createUser(w http.ResponseWriter, req *http.Request) {
	user := models.NewUser()
	vars := mux.Vars(req)
	var payload map[string]interface{}
	decodeErr := json.NewDecoder(req.Body).Decode(&payload)
	payload["email"] = vars["user"]
	if decodeErr != nil {
		http.Error(w, "Unable to read payload", http.StatusBadRequest)
		fmt.Printf("%s\n", decodeErr)
		return
	}

	// Validation
	isValid := ValidateNewUserPayload(payload, w)
	if !isValid {
		return
	}

	// Set Password
	if pwd, hashErr := bcrypt.GenerateFromPassword([]byte(payload["password"].(string)), bcrypt.DefaultCost); hashErr == nil {
		http.Error(w, "Unable to create password hash", 500)
		fmt.Printf("Unable to create password hash\n%s\n", hashErr)
		return
	} else {
		payload["password"] = string(pwd)
	}

	if _, ok := payload["icon"].(string); !ok {
		files, err := ioutil.ReadDir("./assets/img/icon/default")
		if err == nil {
			i := rand.Intn(len(files))
			payload["icon"] = files[i].Name()
		}
	}

	user.UpdateFromMap(payload)

	id, err := db.CreateUser(user)
	if err != nil {
		http.Error(w, "Cannot add user", 500)
		fmt.Printf("Unable to add user %s\n%s\n", user.Email, err)
		return
	}
	user.ID = id.InsertedID.(primitive.ObjectID)

	fmt.Printf("Created user %s\n", user.Email)

	token, _ := auth.RefreshToken(user.Email, user.Admin)

	user.Token = token

	_ = json.NewEncoder(w).Encode(user)

	return
}

func deleteUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var user string = params["user"]
	authorized := req.Context().Value("auth").(auth.Claim)
	if authorized.User != user && authorized.IsAdmin == false {
		http.Error(w, "Cannot delete this account", http.StatusForbidden)
		fmt.Printf("User %s attempted to delete account %s\n", authorized.User, user)
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
	email := params["user"]
	var payload map[string]interface{}
	decodeErr := json.NewDecoder(req.Body).Decode(&payload)
	if decodeErr != nil {
		http.Error(w, "Unable to read payload", 400)
		return
	}
	// Vefify user can make this change
	authorized := req.Context().Value("auth").(auth.Claim)
	if authorized.User != email && authorized.IsAdmin == false {
		http.Error(w, "Cannot update this account", http.StatusForbidden)
		fmt.Printf("User %s attempted to update account %s\n", authorized.User, email)
		return
	}
	// Verify that email doesn't exist if trying to change
	if email, ok := payload["email"]; ok {
		if validEmail, err := db.ValidateNewEmail(email.(string)); err != nil {
			// Error
			http.Error(w, "Unable to verify email address", http.StatusInternalServerError)
			return
		} else if !validEmail {
			// Email already in system
			http.Error(w, "Email already exists", http.StatusBadRequest)
			return
		}
		// No problems, continue on
	}
	// Need to re-hash password if provided
	if pwd, ok := params["password"]; ok {
		pwd, hashErr := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if hashErr != nil {
			http.Error(w, "Unable to create password hash", 500)
			fmt.Printf("Unable to create password hash\n%s\n", hashErr)
			return
		}
		params["password"] = string(pwd)
	}
	// Update user from the payload
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
var CreateUserHandler = http.HandlerFunc(createUser)
