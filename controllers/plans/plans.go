package plans

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"verottaa/controllers/controller_helpers"
	"verottaa/databaser"
	"verottaa/models/dto"
	"verottaa/models/plans"
	"verottaa/utils"
)

func Router(router *mux.Router) {
	PlansRouter(router)
	StepsRouter(router)
}

func PlansRouter(router *mux.Router) {
	router.HandleFunc("/", createPlan).Methods("POST")
	router.HandleFunc("/", getPlans).Methods("GET")
	router.HandleFunc("/", deleteAllPlans).Methods("DELETE")
	router.HandleFunc("/{id}", getPlanById).Methods("GET")
	router.HandleFunc("/{id}", updatePlan).Methods("PUT")
	router.HandleFunc("/{id}", deletePlan).Methods("DELETE")
}

var database = databaser.GetDatabaser()
var variableReader = controller_helpers.GetVariableReader()
var jsonWorker = controller_helpers.GetJsonWorker()

func createPlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var plan plans.Plan
	err := jsonWorker.Decode(r, plan)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "createPlan",
			"error":    err,
			"cause":    "decoding plan",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusBadRequest)
	}
	id, err := database.CreatePlan(plan)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "createPlan",
			"error":    err,
			"cause":    "creating plan",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response := dto.ObjectCreatedDto{Id: utils.IdFromInterfaceToString(id)}
		err = jsonWorker.Encode(w, response)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "plans",
				"function": "createPlan",
				"error":    err,
				"cause":    "encoding plan",
			}).Error("Unexpected error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func getPlans(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allPlans, err := database.ReadAllPlans()
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "getPlans",
			"error":    err,
			"cause":    "read all plans",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = jsonWorker.Encode(w, allPlans)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "plans",
				"function": "getPlans",
				"error":    err,
				"cause":    "encoding all plans",
			}).Error("Unexpected error")
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getPlanById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	plan, err := database.ReadPlanById(id)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "getPlanById",
			"error":    err,
			"cause":    "reading plan by id",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusNotFound)
	} else {
		err = jsonWorker.Encode(w, plan)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "plans",
				"function": "getPlanById",
				"error":    err,
				"cause":    "encoding plan",
			}).Error("Unexpected error")
			w.WriteHeader(http.StatusOK)
		}
	}
}

func updatePlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	var plan plans.Plan
	err := jsonWorker.Decode(r, plan)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "updatePlan",
			"error":    err,
			"cause":    "decoding plan",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusBadRequest)
	}
	err = database.UpdatePlan(id, plan)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "updatePlan",
			"error":    err,
			"cause":    "updating plan",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deletePlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := variableReader.GetObjectIdFromQueryByName(r, "id")
	err := database.DeletePlanById(id)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "deletePlan",
			"error":    err,
			"cause":    "deleting plan",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteAllPlans(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := database.DeleteAllPlans()
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "plans",
			"function": "deleteAllPlans",
			"error":    err,
			"cause":    "deleting all plans",
		}).Error("Unexpected error")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
