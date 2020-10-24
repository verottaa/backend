package controllers

import (
	"github.com/gorilla/mux"
	"verottaa/controllers/assignments"
	"verottaa/controllers/auth"
	"verottaa/controllers/plans"
	"verottaa/controllers/users"
)

func AssignmentsRouter(router *mux.Router) {
	assignments.Router(router)
}

func AuthRouter(router *mux.Router) {
	auth.Router(router)
}

func UserRouter(router *mux.Router) {
	users.Router(router)
}

func PlansRouter(router *mux.Router) {
	plans.Router(router)
}
