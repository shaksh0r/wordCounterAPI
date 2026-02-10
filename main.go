package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func handleQuery(q string) []string {
	wordList := make([]string, 0)

	for _, word := range strings.Fields(q) {
		wordList = append(wordList, word)
	}

	return wordList
}

func countWord(wordList []string, mapping map[string]int) {

	for _, word := range wordList {
		mapping[word] += 1
	}
}

type Mapping struct {
	Words map[string]int `json:"words"`
}

func run(q string) (Mapping, error) {
	mapping := make(map[string]int)
	wordList := handleQuery(q)
	countWord(wordList, mapping)

	wordMapping := Mapping{Words: mapping}

	return wordMapping, nil
}

func main() {

	address := ":8080"
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wordMapping, err := run(r.URL.Query().Get("text"))

		if err != nil {
			log.Fatal(err)
		}

		w.Header().Add("header", "some-value")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(wordMapping)

		if err != nil {
			log.Fatal(err)
		}
	})

	server := &http.Server{
		Addr:           address,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 2 << 10,
	}

	log.Printf("The server is running on %v \n", address)
	log.Fatal(server.ListenAndServe())

}
