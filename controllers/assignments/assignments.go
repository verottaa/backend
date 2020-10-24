package assignments

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"verottaa/databaser"
	"verottaa/models/assignments"
	"verottaa/models/dto"
	"verottaa/utils"
)

func Router(router *mux.Router) {
	router.HandleFunc("/", createAssignment).Methods("POST")
	router.HandleFunc("/", getAssignments).Methods("GET")
	router.HandleFunc("/", deleteAllAssignments).Methods("DELETE")
	router.HandleFunc("/{id}", getAssignmentById).Methods("GET")
	router.HandleFunc("/{id}", updateAssignment).Methods("PUT")
	router.HandleFunc("/{id}", deleteAssignment).Methods("DELETE")
}

var database = databaser.GetDatabaser()

func createAssignment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var assignment assignments.Assignment
	if err := json.NewDecoder(r.Body).Decode(&assignment); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusBadRequest)
	}
	if id, err := database.CreateAssignment(assignment); err != nil {
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

func getAssignments(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if allAssignments, err := database.ReadAllAssignments(); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(allAssignments)
		if err != nil {
			// TODO: логирование
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getAssignmentById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	if assignment, err := database.ReadAssignmentById(id); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(assignment)
		if err != nil {
			// TODO: логирование
			w.WriteHeader(http.StatusOK)
		}
	}
}

func updateAssignment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusBadRequest)
	}
	var assignment assignments.Assignment
	if err = json.NewDecoder(r.Body).Decode(&assignment); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := database.UpdateAssignment(id, assignment); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAssignment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	if err := database.DeleteAssignmentById(id); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAllAssignments(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := database.DeleteAllAssignments(); err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
