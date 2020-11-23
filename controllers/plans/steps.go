package plans

import (
	"github.com/gorilla/mux"
	"net/http"
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
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	var step plans.Step
	err := jsonWorker.Decode(r, step)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusBadRequest)
	}
	createdId, err := database.CreateStepInPlan(id, step)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		response := dto.ObjectCreatedDto{Id: utils.IdFromInterfaceToString(createdId)}
		err = jsonWorker.Encode(w, response)
		if err != nil {
			// TODO: логирование
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func getAllSteps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	steps, err := database.ReadAllStepsInPlan(id)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = jsonWorker.Encode(w, steps)
		if err != nil {
			// TODO: логирование
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getStepById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	stepId := variableReader.GetObjectIdFromQueryByName(r, "stepId")
	step, err := database.ReadStepByIdInPlan(id, stepId)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = jsonWorker.Encode(w, step)
		if err != nil {
			// TODO: логирование
			w.WriteHeader(http.StatusOK)
		}
	}
}

func updateStepById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	stepId := variableReader.GetObjectIdFromQueryByName(r, "stepId")
	var step plans.Step
	err := jsonWorker.Decode(r, step)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusBadRequest)
	}
	err = database.UpdateStepInPlan(id, stepId, step)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteStepById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	stepId := variableReader.GetObjectIdFromQueryByName(r, "stepId")
	err := database.DeleteStepInPlan(id, stepId)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAllSteps(w http.ResponseWriter, r *http.Request) {
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	w.Header().Set("Content-Type", "application/json")
	err := database.DeleteAllStepsInPlan(id)
	if err != nil {
		// TODO: логирование
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
