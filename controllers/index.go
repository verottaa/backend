package controllers

import (
	"github.com/gorilla/mux"
	"verottaa/controllers/assignments"
	"verottaa/controllers/auth"
)

func AssignmentsRouter(router *mux.Router) {
	assignments.Router(router)
}

func AuthRouter(router *mux.Router) {
	auth.Router(router)
}
