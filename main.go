package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"
	"verottaa/controllers"
	"verottaa/variables"
)

func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	api := mux.NewRouter().PathPrefix(variables.ROOT_ROUTE).Subrouter()
	api.Use(mux.CORSMethodMiddleware(api))
	frontend := mux.NewRouter()

	InitControllers(api)

	staticHandler := http.StripPrefix("/", FileServer(http.Dir("./frontend/")))
	frontend.PathPrefix("/").Handler(staticHandler)
	frontend.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	apiSrv := &http.Server{
		Addr:         variables.ApiPort,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      api,
	}

	frontSrv := &http.Server{
		Addr:         variables.FrontendPort,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      frontend,
	}

	go func() {
		log.Println("Api listening on port ", variables.ApiPort)
		if err := apiSrv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	go func() {
		log.Println("Frontend listening on port ", variables.FrontendPort)
		if err := frontSrv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := apiSrv.Shutdown(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = frontSrv.Shutdown(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer cancel()
	log.Println("Server gracefully stopped!")
}

func InitControllers(router *mux.Router) {
	controllers.UserRouter(router.PathPrefix(variables.USERS_ROUTE).Subrouter())
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/index.html")
}

func FileServer(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			NotFoundHandler(w, r)
			return
		}
		fsh.ServeHTTP(w, r)
	})
}
