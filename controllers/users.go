package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"verottaa/databaser"
	"verottaa/models"
	"verottaa/models/dto"
	"verottaa/utils"
	logpack "verottaa/utils/logger"
)

func UserRouter(router *mux.Router) {
	router.HandleFunc("/", createUser).Methods("POST")
	router.HandleFunc("/", getUsers).Methods("GET")
	router.HandleFunc("/", deleteAllUsers).Methods("DELETE")
	router.HandleFunc("/{id}", getUserById).Methods("GET")
	router.HandleFunc("/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/{id}", deleteUser).Methods("DELETE")
}

var database = databaser.GetDatabaser()
var logTag = "USER_CONTROLLER"
var logger *logpack.Logger

func init() {
	logger = logpack.CreateLogger(logTag)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	if id, err := database.CreateUser(user); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		response := dto.ObjectCreatedDto{Id: utils.IdFromInterfaceToString(id)}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func getUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if users, err := database.ReadAllUsers(); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	if users, err := database.ReadUserById(id); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := database.UpdateUser(id, user); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	if err := database.DeleteUserById(id); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := database.DeleteAllUsers(); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
