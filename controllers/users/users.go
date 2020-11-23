package users

import (
	"github.com/gorilla/mux"
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
var variableReader = controller_helpers.GetVariableReader()
var jsonWorker = controller_helpers.GetJsonWorker()

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user users.User
	err := jsonWorker.Decode(r, user)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusBadRequest)
	}
	id, err := database.CreateUser(user)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		response := dto.ObjectCreatedDto{Id: utils.IdFromInterfaceToString(id)}
		err = jsonWorker.Encode(w, response)
		if err != nil {
			// TODO: логирование
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func getUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allUsers, err := database.ReadAllUsers()
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = jsonWorker.Encode(w, allUsers)
		if err != nil {
			// TODO: логирование
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	user, err := database.ReadUserById(id)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = jsonWorker.Encode(w, user)
		if err != nil {
			// TODO: логирование
			w.WriteHeader(http.StatusOK)
		}
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	var user users.User
	err := jsonWorker.Decode(r, user)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusBadRequest)
	}
	err = database.UpdateUser(id, user)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	err := database.DeleteUserById(id)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAllUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := database.DeleteAllUsers()
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
