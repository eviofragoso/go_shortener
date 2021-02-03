package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

func url(token string) string {
	result := readFileData("database.json")

	for k, v := range result {
		if v.(string) == token {
			return k
		}
	}

	return ""
}

func readFileData(fileName string) map[string]interface{} {
	database, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer database.Close()

	byteValue, _ := ioutil.ReadAll(database)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

func retrieveToken(original string) string {
	// Open our jsonFile
	result := readFileData("database.json")

	if result[original] != nil {
		return fmt.Sprintf("%v", result[original])
	}

	return ""
}

func createToken(original string) string {
	result := readFileData("database.json")
	re := regexp.MustCompile(`http(?:s)?://(?:\w+\.)?\w+\.\w+(?:\.\w+)?`)

	if re.MatchString(original) == false {
		return ""
	}

	urlHash := hash(original)
	result[original] = fmt.Sprintf("%v", urlHash)
	//
	// Open :our jsonFile
	file, _ := json.MarshalIndent(result, "", " ")
	_ = ioutil.WriteFile("database.json", file, 0644)

	return fmt.Sprintf("%v", urlHash)
}

func shortenedURL(original string) string {
	domain := os.Getenv("SHORTENER_DOMAIN")
	token := token(original)

	if token == "" {
		return token
	}

	return fmt.Sprintf("%v/%v", domain, token)
}

// GetShortenedURL retrieve shortened url based on original url
func serveShortenedURL(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	original := params["url"]

	if original == "" {
		json.NewEncoder(w).Encode("")

		return
	}

	json.NewEncoder(w).Encode(shortenedURL(original))
}

func redirectToURL(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hash := params["hash"]
	url := url(hash)

	if url != "" {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func main() {
	if os.Getenv("GO_ENV") == "" || os.Getenv("GO_ENV") == "development" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	router := mux.NewRouter()

	router.HandleFunc("/shortener", serveShortenedURL).Queries("url", "{url}").Methods("GET")
	router.HandleFunc("/{hash}", redirectToURL).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
