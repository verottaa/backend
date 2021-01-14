package users

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"verottaa/controllers/controller_helpers"
	"verottaa/databaser"
	"verottaa/models/dto"
	"verottaa/models/users"
	"verottaa/utils"
)

func Router(router *mux.Router) {
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
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "users",
			"function": "createUser",
			"error":    err,
			"cause":    "decoding json",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusBadRequest)
	}
	id, err := database.CreateUser(user)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "users",
			"function": "createUser",
			"error":    err,
			"cause":    "creating user",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		response := dto.ObjectCreatedDto{Id: utils.IdFromInterfaceToString(id)}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "users",
				"function": "createUser",
				"error":    err,
				"cause":    "encoding results",
			}).Error("Unexpected error")
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func getUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allUsers, err := database.ReadAllUsers()
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "users",
			"function": "getUsers",
			"error":    err,
			"cause":    "reading all users",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(allUsers)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "users",
				"function": "getUsers",
				"error":    err,
				"cause":    "encoding results",
			}).Error("Unexpected error")
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	var variableReader = controller_helpers.GetVariableReader()
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	user, err := database.ReadUserById(id)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "users",
			"function": "getUserById",
			"error":    err,
			"cause":    "read user by id",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "users",
				"function": "getUserById",
				"error":    err,
				"cause":    "encoding results",
			}).Error("Unexpected error")
			w.WriteHeader(http.StatusOK)
		}
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	var variableReader = controller_helpers.GetVariableReader()
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	var user users.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "users",
			"function": "updateUser",
			"error":    err,
			"cause":    "decoding user",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusBadRequest)
	}
	err = database.UpdateUser(id, user)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "users",
			"function": "updateUser",
			"error":    err,
			"cause":    "updating user",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	var variableReader = controller_helpers.GetVariableReader()
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	err := database.DeleteUserById(id)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "users",
			"function": "deleteUser",
			"error":    err,
			"cause":    "deleting user",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAllUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := database.DeleteAllUsers()
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "users",
			"function": "deleteAllUsers",
			"error":    err,
			"cause":    "deleting all users",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
