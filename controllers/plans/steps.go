package plans

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"verottaa/controllers/controller_helpers"
	"verottaa/models/dto"
	"verottaa/models/plans"
	"verottaa/utils"
)

func StepsRouter(router *mux.Router) {
	router.HandleFunc("/{id}/steps/", getAllSteps).Methods("GET")
	router.HandleFunc("/{id}/steps/", createStep).Methods("POST")
	router.HandleFunc("/{id}/steps/", deleteAllSteps).Methods("DELETE")
	router.HandleFunc("/{id}/steps/{stepId}", getStepById).Methods("GET")
	router.HandleFunc("/{id}/steps/{stepId}", updateStepById).Methods("PUT")
	router.HandleFunc("/{id}/steps/{stepId}", deleteStepById).Methods("DELETE")
}

func createStep(w http.ResponseWriter, r *http.Request) {
	var variableReader = controller_helpers.GetVariableReader()
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	var step plans.Step
	err := json.NewDecoder(r.Body).Decode(&step)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "createStep",
			"error":    err,
			"cause":    "decoding step",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusBadRequest)
	}
	createdId, err := database.CreateStepInPlan(id, step)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "createStep",
			"error":    err,
			"cause":    "creating step and saving it to plan",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		response := dto.ObjectCreatedDto{Id: utils.IdFromInterfaceToString(createdId)}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "plans",
				"function": "createStep",
				"error":    err,
				"cause":    "encoding results",
			}).Error("Unexpected error")
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func getAllSteps(w http.ResponseWriter, r *http.Request) {
	var variableReader = controller_helpers.GetVariableReader()
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	steps, err := database.ReadAllStepsInPlan(id)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "getAllSteps",
			"error":    err,
			"cause":    "reading all steps from plan",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(steps)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "plans",
				"function": "getAllSteps",
				"error":    err,
				"cause":    "encoding steps",
			}).Error("Unexpected error")
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getStepById(w http.ResponseWriter, r *http.Request) {
	var variableReader = controller_helpers.GetVariableReader()
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	stepId := variableReader.GetObjectIdFromQueryByName(r, "stepId")
	step, err := database.ReadStepByIdInPlan(id, stepId)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "getStepById",
			"error":    err,
			"cause":    "reading step by id in plan",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(step)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "plans",
				"function": "getStepById",
				"error":    err,
				"cause":    "encoding step",
			}).Error("Unexpected error")
			w.WriteHeader(http.StatusOK)
		}
	}
}

func updateStepById(w http.ResponseWriter, r *http.Request) {
	var variableReader = controller_helpers.GetVariableReader()
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	stepId := variableReader.GetObjectIdFromQueryByName(r, "stepId")
	var step plans.Step
	err := json.NewDecoder(r.Body).Decode(&step)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "updateStepById",
			"error":    err,
			"cause":    "decoding step",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusBadRequest)
	}
	err = database.UpdateStepInPlan(id, stepId, step)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "updateStepById",
			"error":    err,
			"cause":    "updating step in plan",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteStepById(w http.ResponseWriter, r *http.Request) {
	var variableReader = controller_helpers.GetVariableReader()
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	stepId := variableReader.GetObjectIdFromQueryByName(r, "stepId")
	err := database.DeleteStepInPlan(id, stepId)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "deleteStepById",
			"error":    err,
			"cause":    "deleting step in plan",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAllSteps(w http.ResponseWriter, r *http.Request) {
	var variableReader = controller_helpers.GetVariableReader()
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	w.Header().Set("Content-Type", "application/json")
	err := database.DeleteAllStepsInPlan(id)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "deleteAllSteps",
			"error":    err,
			"cause":    "deleting all steps in plan",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
