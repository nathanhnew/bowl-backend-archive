package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/nathanhnew/bowl-backend/internal/app/config"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var cfg, _ = config.GetConfig(config.DefaultConfigLocation)

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func parseTokenClaims(tkn *jwt.Token, w *http.ResponseWriter, req *http.Request) {
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		if claims["iss"] == nil || claims["expiresAt"] == nil {
			fmt.Println("here")
			http.Error(*w, "Invalid token", http.StatusForbidden)
			return
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
			return cfg.GetValidationKey(), nil
		})
		if err != nil {
			http.Error(w, "Unable to verify token", http.StatusForbidden)
			return
		}
		parseTokenClaims(tkn, &w, req)
		next.ServeHTTP(w, req)
	})
}

func RefreshToken(email string, admin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     email,
		"admin":     admin,
		"expiresAt": time.Now().Add(time.Hour * time.Duration(cfg.GetTokenTimeout())),
	})
	tokenString, err := token.SignedString(cfg.GetValidationKey())
	if err != nil {
		return "", err
		fmt.Println("generating token", err)
	}
	return tokenString, nil
}
