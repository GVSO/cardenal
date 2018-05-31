package main

import (
	"log"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	linkedinLogin(&w, r)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	linkedinCallback(&w, r)
}

func clientHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	if Settings.Development {
		log.Println("path:", "client/dist"+path)
	}

	http.ServeFile(w, r, "client/dist"+path)
}
