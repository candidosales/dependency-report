package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func  (app *App) routes() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)
	router.Use(mux.CORSMethodMiddleware(router))

	router.HandleFunc("/generate-report", app.GenerateReportHandler).Methods("GET")
	router.HandleFunc("/ping", app.PingHandler).Methods("GET")
	router.HandleFunc("/", app.RootHandler).Methods("GET")
	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
