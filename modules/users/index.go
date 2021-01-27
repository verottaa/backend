package users

import (
	"github.com/gorilla/mux"
	"verottaa/modules/users/controllers"
)

func RegistryControllers(router *mux.Router) {
	var controller = controllers.CreateController()
	controller.InitController(router)
}
