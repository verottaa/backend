package assignments

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"verottaa/controllers/controller_helpers"
	"verottaa/databaser"
	"verottaa/models/assignments"
	"verottaa/models/dto"
	"verottaa/utils"
)

func Router(router *mux.Router) {
	router.HandleFunc("/", Assign).Methods("POST")
	router.HandleFunc("/", getAssignments).Methods("GET")
	router.HandleFunc("/", deleteAllAssignments).Methods("DELETE")
	router.HandleFunc("/{id}", getAssignmentById).Methods("GET")
	router.HandleFunc("/{id}", updateAssignment).Methods("PUT")
	router.HandleFunc("/{id}", deleteAssignment).Methods("DELETE")
}

var database = databaser.GetDatabaser()
var variableReader = controller_helpers.GetVariableReader()
var jsonWorker = controller_helpers.GetJsonWorker()

func Assign(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var assignDto dto.AssignCreateDto
	err := jsonWorker.Decode(r, assignDto)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "assignments",
			"function": "Assign",
			"error":    err,
			"cause":    "decoding assignment",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusBadRequest)
	}

	assignment := assignments.NewAssignment(assignDto)
	id, err := database.CreateAssignment(assignment)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "assignments",
			"function": "Assign",
			"error":    err,
			"cause":    "creating assignment",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		response := dto.ObjectCreatedDto{Id: utils.IdFromInterfaceToString(id)}
		err = jsonWorker.Encode(w, response)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "assignments",
				"function": "Assign",
				"error":    err,
				"cause":    "encoding results",
			}).Error("Unexpected error")
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func getAssignments(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allAssignments, err := database.ReadAllAssignments()
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "assignments",
			"function": "getAssignments",
			"error":    err,
			"cause":    "reading all assignments",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = jsonWorker.Encode(w, allAssignments)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "assignments",
				"function": "getAssignments",
				"error":    err,
				"cause":    "encoding results",
			}).Error("Unexpected error")
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getAssignmentById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	assignment, err := database.ReadAssignmentById(id)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "assignments",
			"function": "getAssignmentById",
			"error":    err,
			"cause":    "reading assignment",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = jsonWorker.Encode(w, assignment)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "assignments",
				"function": "getAssignmentById",
				"error":    err,
				"cause":    "encoding assignment",
			}).Error("Unexpected error")
			w.WriteHeader(http.StatusOK)
		}
	}
}

func updateAssignment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	var assignment assignments.Assignment
	err := jsonWorker.Decode(r, assignment)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "assignments",
			"function": "updateAssignment",
			"error":    err,
			"cause":    "decoding assignment",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := database.UpdateAssignment(id, assignment); err != nil {
		log.WithFields(log.Fields{
			"package":  "assignments",
			"function": "updateAssignment",
			"error":    err,
			"cause":    "updating assignment",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAssignment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	err := database.DeleteAssignmentById(id)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "assignments",
			"function": "deleteAssignment",
			"error":    err,
			"cause":    "deleting assignment",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAllAssignments(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := database.DeleteAllAssignments()
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "assignments",
			"function": "deleteAllAssignments",
			"error":    err,
			"cause":    "deleting all assignments",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
