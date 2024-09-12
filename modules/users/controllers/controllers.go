package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Controller interface {
	initial
	createUser
	reader
	updater
	deletion
}

type initial interface {
	InitController(router *mux.Router)
}

type createUser interface {
	createOne(w http.ResponseWriter, r *http.Request)
}

type reader interface {
	getOne(w http.ResponseWriter, r *http.Request)
	getAll(w http.ResponseWriter, r *http.Request)
}

type updater interface {
	updateOne(w http.ResponseWriter, r *http.Request)
}

type deletion interface {
	deleteOne(w http.ResponseWriter, r *http.Request)
	deleteAll(w http.ResponseWriter, r *http.Request)
	deleteMany(w http.ResponseWriter, r *http.Request)
}
