package auth

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/nathanhnew/bowl-backend/internal/app/config"
	"github.com/nathanhnew/bowl-backend/internal/app/db"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var cfg, _ = config.GetConfig("")

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	user, err := db.GetUser(creds.Email)
	fmt.Printf("%+v\n", user)
	if err != nil {
		http.Error(w, "Unable to verify account", http.StatusBadRequest)
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

	token, err := refreshToken(user.Email, user.Admin)
	if err != nil {
		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
		fmt.Printf("Unable to generate token for %s\n%s\n", user.Email, err)
		return
	}

	user.Token = token

	json.NewEncoder(w).Encode(user)

	return

}

func parseTokenClaims(tkn *jwt.Token, w *http.ResponseWriter, req *http.Request) {
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		if claims["iss"] == nil || claims["expiresAt"] == nil {
			fmt.Println("here")
			http.Error(*w, "Invalid token", http.StatusForbidden)
			return
		}
		if claims["iss"] != cfg.Values["validationKey"].(string) {
			http.Error(*w, "Invalid authorization token", http.StatusForbidden)
		}
		if exp, _ := time.Parse(time.RFC3339, claims["expiresAt"].(string)); ok {
			if exp.Before(time.Now()) {
				http.Error(*w, "Token Expired", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(*w, "Invalid expiration", http.StatusForbidden)
			return
		}
		// If make it here, everything is OK
		req.Header.Set("authUser", claims["email"].(string))
		req.Header.Set("authAdmin", strconv.FormatBool(claims["admin"].(bool)))
	} else {
		http.Error(*w, "Invalid authorization token", http.StatusForbidden)
		return
	}
}

func ValidToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "No authorization provided", http.StatusForbidden)
			return
		} else if len(strings.Split(token, " ")) != 2 {
			http.Error(w, "Invalid authorization", http.StatusBadRequest)
			return
		}
		token = strings.Split(token, " ")[1]
		tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error validating token authenticity")
			}
			return []byte(cfg.Values["validationPassKey"].(string)), nil
		})
		if err != nil {
			http.Error(w, "Unable to verify token", http.StatusForbidden)
			return
		}
		parseTokenClaims(tkn, &w, req)
		next.ServeHTTP(w, req)
	})
}

func refreshToken(email string, admin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":       cfg.Values["validationKey"].(string),
		"email":     email,
		"admin":     admin,
		"expiresAt": time.Now().Add(time.Hour * time.Duration(cfg.Values["tokenTimeout"].(float64))),
	})
	tokenString, err := token.SignedString([]byte(cfg.Values["validationPassKey"].(string)))
	if err != nil {
		return "", err
		fmt.Println("generating token", err)
	}
	return tokenString, nil
}
