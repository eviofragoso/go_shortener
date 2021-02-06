package shortener

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
)

func generateHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func retrieveValue(original string) string {
	token := retrieveToken(original)
	if token == "" {
		token = createToken(original)
	}

	return token
}

func retrieveKey(token string) string {
	result := readFileData()

	for k, v := range result {
		if v.(string) == token {
			return k
		}
	}

	return ""
}

func readFileData() map[string]interface{} {
	fileName := os.Getenv("DB_NAME")
	fmt.Printf("%v", fileName)
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
	result := readFileData()

	if result[original] != nil {
		return fmt.Sprintf("%v", result[original])
	}

	return ""
}

func createToken(original string) string {
	result := readFileData()
	re := regexp.MustCompile(`http(?:s)?://(?:\w+\.)?\w+\.\w+(?:\.\w+)?`)

	if re.MatchString(original) == false {
		return ""
	}

	urlHash := generateHash(original)
	result[original] = fmt.Sprintf("%v", urlHash)
	//
	// Open :our jsonFile
	file, _ := json.MarshalIndent(result, "", " ")
	_ = ioutil.WriteFile("database.json", file, 0644)

	return fmt.Sprintf("%v", urlHash)
}

func shortenedURL(original string) string {
	domain := os.Getenv("SHORTENER_DOMAIN")
	token := retrieveValue(original)

	if token == "" {
		return token
	}

	return fmt.Sprintf("%v/%v", domain, token)
}

// ServeShortenedURL retrieve shortened url based on original url
func ServeShortenedURL(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	original := params["url"]

	w.Header().Set("Content-Type", "application/json")

	if original == "" {
		json.NewEncoder(w).Encode("")

		return
	}

	json.NewEncoder(w).Encode(shortenedURL(original))
}

// RedirectToURL redirects to url
func RedirectToURL(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hash := params["hash"]
	// retrieve url from database
	url := retrieveKey(hash)

	if url != "" {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}
