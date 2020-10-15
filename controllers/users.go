package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"verottaa/databaser"
	"verottaa/models/dto"
	"verottaa/models/users"
	"verottaa/utils"
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

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user users.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusBadRequest)
	}
	if id, err := database.CreateUser(user); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		response := dto.ObjectCreatedDto{Id: utils.IdFromInterfaceToString(id)}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			// TODO: логирование
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func getUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if allUsers, err := database.ReadAllUsers(); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(allUsers)
		if err != nil {
			// TODO: логирование
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	if user, err := database.ReadUserById(id); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			// TODO: логирование
			w.WriteHeader(http.StatusOK)
		}
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	var user users.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := database.UpdateUser(id, user); err != nil {
		// TODO: логирование
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
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAllUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := database.DeleteAllUsers(); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
