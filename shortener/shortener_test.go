package shortener

import (
	"fmt"
	"testing"
)

func TestGenerateHash(t *testing.T) {
	want := "uint32"
	got := fmt.Sprintf("%T", generateHash("http://example.com"))

	if want != got {
		t.Fatalf(`generateHash("http://example.com") = %q, is different from want %q`, got, want)
	}
}

// func TestRetrieveValue(t *testing.T) {
// 	// utils.LoadDotEnv()
// 	utils.InitDBFile("database.json")
// 	createToken("https://google.com")

// 	want := fmt.Sprintf("%v", generateHash("https://google.com"))
// 	got := fmt.Sprintf("%v", retrieveValue("https://google.com"))

// 	if want != got {
// 		t.Fatalf(`retrieveValue("https://google.com") = %q, is different from want %q`, got, want)
// 	}
// }
