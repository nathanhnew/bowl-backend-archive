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

func parseTokenClaims(tkn *jwt.Token, w *http.ResponseWriter, req *http.Request) bool {
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		if claims["expiresAt"] == nil {
			http.Error(*w, "Invalid token", http.StatusForbidden)
			return false
		}
		if exp, _ := time.Parse(time.RFC3339, claims["expiresAt"].(string)); ok {
			if exp.Before(time.Now()) {
				http.Error(*w, "Token Expired", http.StatusForbidden)
				return false
			}
		} else {
			http.Error(*w, "Invalid expiration", http.StatusUnauthorized)
			return false
		}
		// If make it here, everything is OK
		req.Header.Set("authUser", claims["email"].(string))
		req.Header.Set("authAdmin", strconv.FormatBool(claims["admin"].(bool)))
		return true
	} else {
		http.Error(*w, "Invalid authorization token", http.StatusForbidden)
		return false
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
		validToken := parseTokenClaims(tkn, &w, req)
		if !validToken {
			return
		}
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
