package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eviofragoso/go_shortener/shortener"
	"github.com/eviofragoso/go_shortener/utils"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize Dependencies
	utils.LoadDotEnv()
	utils.InitDBFile(os.Getenv("DB_NAME"))

	// instantiate router
	router := mux.NewRouter()
	// define routes
	router.HandleFunc("/shortener", shortener.ServeShortenedURL).Queries("url", "{url}").Methods("GET")
	router.HandleFunc("/{hash}", shortener.RedirectToURL).Methods("GET")
	// define path
	servePath := fmt.Sprintf(":%v", os.Getenv("PORT"))
	// serve
	log.Fatal(http.ListenAndServe(servePath, router))
}
