package main

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"
	"verottaa/config"
	"verottaa/constants"
	"verottaa/controllers"
	"verottaa/utils/logger"
)

func main() {
	const logTag = "MAIN"
	logger := logger.CreateLogger(logTag)
	configuration := config.GetConfiguration()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	api := mux.NewRouter().PathPrefix(constants.ROOT_ROUTE).Subrouter()
	api.Use(mux.CORSMethodMiddleware(api))
	frontend := mux.NewRouter()

	InitControllers(api)

	staticHandler := http.StripPrefix("/", FileServer(http.Dir("./frontend/")))
	frontend.PathPrefix("/").Handler(staticHandler)
	frontend.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})

	apiSrv := &http.Server{
		Addr:         configuration.GetApiPort(),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      handlers.LoggingHandler(os.Stdout, handlers.CORS(originsOk, headersOk, methodsOk)(api)),
	}

	frontSrv := &http.Server{
		Addr:         configuration.GetStaticPort(),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      frontend,
	}

	go func() {
		logger.Info("Api listening on port ", configuration.GetApiPort())
		if err := apiSrv.ListenAndServe(); err != nil {
			logger.Error(err)
		}
	}()

	go func() {
		logger.Info("Frontend listening on port ", configuration.GetStaticPort())
		if err := frontSrv.ListenAndServe(); err != nil {
			logger.Error(err)
		}
	}()

	<-stopChan

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := apiSrv.Shutdown(ctx)
	if err != nil {
		logger.Error(err)
	}
	err = frontSrv.Shutdown(ctx)
	if err != nil {
		logger.Error(err)
	}
	defer cancel()
	logger.Info("Server gracefully stopped!")
}

func InitControllers(router *mux.Router) {
	router.StrictSlash(true).HandleFunc("/", StatusApi).Methods("GET")
	controllers.UserRouter(router.PathPrefix(constants.USERS_ROUTE).Subrouter())
	controllers.AuthRouter(router.PathPrefix(constants.AUTH_ROUTE).Subrouter())
}

func StatusApi(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server is running"))
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
