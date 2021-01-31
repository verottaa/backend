package app

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"time"
	"verottaa/config"
	"verottaa/constants"
	"verottaa/controllers"
	_ "verottaa/modules/database"
	"verottaa/modules/plans"
	"verottaa/modules/users"
)

var configuration config.Configurable

func infoLogger(info ...string) {
	fmt.Println(info)
}

func errorLogger(err error) {
	fmt.Println(err.Error())
}

func StartApplication() {
	configuration = initConfiguration()
	router := initRouter()
	stopChan := initStopChan()

	rootRouteSubrouter := router.PathPrefix(constants.RootRoute).Subrouter()
	initControllers(rootRouteSubrouter)

	router = initFileServerAndNotFoundHandlers(router)

	originsOk, headersOk, methodsOk := initCORSOptions()
	corsHandler := handlers.CORS(originsOk, headersOk, methodsOk)(router)

	server := setServer(corsHandler)

	go func() {
		infoLogger("Api listening on port ", configuration.GetPort())
		err := server.ListenAndServe()
		if err != nil {
			errorLogger(err)
		}
	}()

	<-stopChan

	infoLogger("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := server.Shutdown(ctx)
	if err != nil {
		errorLogger(err)
	}
	defer cancel()
	infoLogger("Server gracefully stopped!")
}

func initConfiguration() config.Configurable {
	return config.GetConfiguration()
}

func initRouter() *mux.Router {
	return mux.NewRouter()
}

func initStopChan() chan os.Signal {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	return stopChan
}

func initControllers(router *mux.Router) {
	router.StrictSlash(true).HandleFunc("/", statusApi).Methods("GET")
	users.RegistryControllers(router.PathPrefix(constants.UsersRoute).Subrouter())
	plans.RegistryControllers(router.PathPrefix(constants.PlansRoute).Subrouter())
	controllers.AssignmentsRouter(router.PathPrefix(constants.AssignmentsRoute).Subrouter())
	controllers.AuthRouter(router.PathPrefix(constants.AuthRoute).Subrouter())
}

func statusApi(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Server is running"))
	if err != nil {
		errorLogger(err)
	}
}

func initFileServerAndNotFoundHandlers(router *mux.Router) *mux.Router {
	router = initFileServerHandler(router)
	router = initNotFoundHandler(router)
	return router
}

func initFileServerHandler(router *mux.Router) *mux.Router {
	var fileServer = initFileServer("./frontend/")
	staticHandler := http.StripPrefix("/", fileServer)
	router.PathPrefix("/").Handler(staticHandler)
	return router
}

func initFileServer(path string) http.Handler {
	return fileServer(http.Dir(path))
}

func fileServer(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", getMimeTypeForFile(r.URL.Path))
		fmt.Println(w.Header().Get("Content-Type"))
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			NotFoundHandler(w, r)
			return
		}
		fsh.ServeHTTP(w, r)
	})
}

func getMimeTypeForFile(file string) string {
	ext := filepath.Ext(file)
	switch ext {
	case ".htm", ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	default:
		return mime.TypeByExtension(ext)
	}
}

func initNotFoundHandler(router *mux.Router) *mux.Router {
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	return router
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/index.html")
}

func initCORSOptions() (handlers.CORSOption, handlers.CORSOption, handlers.CORSOption) {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	return headersOk, originsOk, methodsOk
}

func setServer(handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:         configuration.GetPort(),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      handlers.LoggingHandler(os.Stdout, handler),
	}
	return server
}
