package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Exception for request errors.
type Exception struct {
	Message string `json:"message"`
}

func main() {

	router := mux.NewRouter()

	servicesRouter := router.PathPrefix("/api/services").Subrouter()
	servicesRouter.HandleFunc("/login", loginHandler)
	servicesRouter.HandleFunc("/login/callback", callbackHandler)

	//dataRouter := router.PathPrefix("/api/data").Subrouter()

	router.HandleFunc("/{rest:.*}", clientHandler)

	port := ":" + strconv.Itoa(Settings.Port)

	log.Print("Server application started at ", "http://localhost"+port)

	log.Fatal(http.ListenAndServe(port, router))
}
