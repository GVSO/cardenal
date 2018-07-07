package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gvso/cardenal/settings"
)

func main() {

	router := mux.NewRouter()

	servicesRouter := router.PathPrefix("/api/services").Subrouter()
	servicesRouter.HandleFunc("/login", loginHandler)
	servicesRouter.HandleFunc("/login/callback", callbackHandler)

	//dataRouter := router.PathPrefix("/api/data").Subrouter()

	router.HandleFunc("/{rest:.*}", clientHandler)

	port := ":" + settings.Port

	log.Print("Server application started at ", "http://localhost"+port)

	log.Fatal(http.ListenAndServe(port, router))
}
