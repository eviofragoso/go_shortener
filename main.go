package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func token(original string) string {
	token := retrieveToken(original)
	if token == "" {
		token = createToken(original)
	}

	return token
}

func retrieveToken(original string) string {
	// Open our jsonFile
	database, err := os.Open("database.json")
	if err != nil {
		fmt.Println(err)
	}
	defer database.Close()

	byteValue, _ := ioutil.ReadAll(database)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	if result[original] != nil {
		return fmt.Sprintf("%v", result[original])
	}

	return ""
}

func createToken(original string) string {

	urlHash := hash(original)
	// Open :our jsonFile
	database, err := os.Open("database.json")
	if err != nil {
		fmt.Println(err)
	}
	defer database.Close()

}

func shortenedURL(original string) string {
	domain := os.Getenv("DOMAIN")

}

// GetShortenedURL retrieve shortened url based on original
func serveShortenedURL(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	original := params["url"]

	json.NewEncoder(w).Encode(shortenedURL(original))
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", serveShortenedURL).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
