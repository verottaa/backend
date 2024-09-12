package plans

import (
	"github.com/gorilla/mux"
	"verottaa/modules/plans/controllers"
)

func RegistryControllers(router *mux.Router) {
	var controller = controllers.CreateController()
	controller.InitController(router)
}
