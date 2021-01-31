package controllers

import (
	"github.com/gorilla/mux"
	"verottaa/controllers/assignments"
)

func AssignmentsRouter(router *mux.Router) {
	assignments.Router(router)
}
