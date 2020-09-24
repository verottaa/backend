package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"verottaa/constants"
)

func AuthRouter(router *mux.Router) {
	router.HandleFunc("/signIn", signInHandler).Methods("POST")
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":    "username",
		"expired": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(constants.SIGNING_KEY))

	w.Write([]byte(tokenString))
}
