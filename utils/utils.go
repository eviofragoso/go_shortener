package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadDotEnv reads .env file if in development env
func LoadDotEnv() {
	if os.Getenv("GO_ENV") == "" || os.Getenv("GO_ENV") == "development" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

// InitDBFile creates db file if it doesn't existis
func InitDBFile(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, _ := json.MarshalIndent("{}", "", " ")
		_ = ioutil.WriteFile(filename, file, 0644)
	}
}
