package auth

import (
	"encoding/json"
	"fmt"
	"github.com/nathanhnew/bowl-backend/internal/app/db"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
)

func ValidEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.Match([]byte(email)) {
		return true
	} else {
		return false
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Unable to parse payload", http.StatusBadRequest)
		fmt.Printf("Unable to parse payload %s\n", err)
		return
	}
	if creds.Email == "" || creds.Password == "" {
		http.Error(w, "Bad payload", http.StatusBadRequest)
		return
	}
	if !ValidEmail(creds.Email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}
	user, err := db.GetUser(creds.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusForbidden)
		fmt.Printf("Unable to verify user: %s\n%s\n", creds.Email, err)
		return
	}
	if user.Active == false {
		http.Error(w, "Account inactive, please reactivate", http.StatusForbidden)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Passwords don't match", http.StatusUnauthorized)
		return
	}

	token, err := RefreshToken(user.Email, user.Admin)
	if err != nil {
		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
		fmt.Printf("Unable to generate token for %s\n%s\n", user.Email, err)
		return
	}

	user.Token = token

	json.NewEncoder(w).Encode(user)

	return

}
