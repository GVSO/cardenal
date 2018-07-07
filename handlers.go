package main

import (
	"log"
	"net/http"

	"github.com/gvso/cardenal/linkedin"
	"github.com/gvso/cardenal/settings"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	linkedin.Login(w, r)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	linkedin.Callback(w, r)
}

func clientHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	if settings.Development {
		log.Println("path:", "client/dist"+path)
	}

	http.ServeFile(w, r, "client/dist"+path)
}
