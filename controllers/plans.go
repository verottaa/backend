package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"verottaa/models"
	"verottaa/models/dto"
	"verottaa/utils"
)

func PlansRouter(router *mux.Router) {
	router.HandleFunc("/", createPlan).Methods("POST")
	router.HandleFunc("/", getPlans).Methods("GET")
	router.HandleFunc("/", deleteAllPlans).Methods("DELETE")
	router.HandleFunc("/{id}", getPlanById).Methods("GET")
	router.HandleFunc("/{id}", updatePlan).Methods("PUT")
	router.HandleFunc("/{id}", deletePlan).Methods("DELETE")
	/*router.HandleFunc("/{id}/steps/", createStep).Methods("POST")
	router.HandleFunc("/{id}/steps/", getAllSteps).Methods("GET")
	router.HandleFunc("/{id}/steps/", deleteAllSteps).Methods("DELETE")
	router.HandleFunc("/{id}/steps/{stepId}", getStepById).Methods("GET")
	router.HandleFunc("/{id}/steps/{stepId}", updateStepById).Methods("PUT")
	router.HandleFunc("/{id}/steps/{stepId}", deleteStepById).Methods("DELETE")*/
}

func createPlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var plan models.Plan
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	if id, err := database.CreatePlan(plan); err != nil {
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

func getPlans(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if plans, err := database.ReadAllPlans(); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(plans)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getPlanById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	if plans, err := database.ReadPlanById(id); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(plans)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func updatePlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	var plan models.Plan
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := database.UpdatePlan(id, plan); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deletePlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	if err := database.DeletePlanById(id); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAllPlans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := database.DeleteAllPlans(); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// STEPS:

/*func createStep(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	var step models.Step
	if err := json.NewDecoder(r.Body).Decode(&step); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	if id, err := database.CreateStepInPlan(id, step); err != nil {
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

func getAllSteps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	if steps, err := database.ReadAllStepsInPlan(id); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(steps)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getStepById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	stepId, _ := primitive.ObjectIDFromHex(vars["stepId"])
	if step, err := database.ReadStepByIdInPlan(id, stepId); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = json.NewEncoder(w).Encode(step)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func updateStepById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	stepId, _ := primitive.ObjectIDFromHex(vars["stepId"])
	var step models.Step
	if err := json.NewDecoder(r.Body).Decode(&step); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := database.UpdateStepInPlan(id, stepId, step); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteStepById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	stepId, _ := primitive.ObjectIDFromHex(vars["stepId"])
	if err := database.DeleteStepInPlan(id, stepId); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAllSteps(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	if err := database.DeleteAllStepsInPlan(id); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
*/