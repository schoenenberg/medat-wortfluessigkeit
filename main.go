// Copyright (c) 2019 Maximilian Schoenenberg
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	//go:embed etc/words.csv
	b             []byte
	filteredWords []string
)

type Word struct {
	Shuffled string
	Solution string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Get environment variables
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	_, debug := os.LookupEnv("DEBUG")

	// Parse file
	words, err := readWords(b)
	if err != nil {
		log.Fatalln(err)
	}

	// filter words
	if debug {
		log.Println("Number of elements in list:", len(words))
		filteredWords = filterWords(words)
		log.Println("Number of elements in filtered list:", len(filteredWords))
	} else {
		filteredWords = filterWords(words)
	}

	// initialise handles
	http.Handle("/", http.RedirectHandler("/index.html", http.StatusPermanentRedirect))
	http.HandleFunc("/word/new", handleWordRequest)

	// start http server
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleWordRequest(w http.ResponseWriter, r *http.Request) {
	// get a random word
	word := filteredWords[rand.Intn(len(filteredWords))]

	// create result struct
	wordResponse := Word{Solution: strings.ToUpper(word), Shuffled: strings.ToUpper(string(shuffleWord(word)))}
	log.Println(wordResponse)

	// Marshal response
	b, err := json.Marshal(wordResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func readWords(b []byte) ([]string, error) {
	// Read csv list
	reader := csv.NewReader(bytes.NewReader(b))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var recordsList []string

	for _, i := range records {
		recordsList = append(recordsList, i[0])
	}

	return recordsList, nil
}

func filterWords(wordList []string) []string {
	filteredWordList := make([]string, 0)

	// filter word by length and not allowed characters
	for _, word := range wordList {
		if len(word) < 7 || len(word) > 9 {
			continue
		} else if strings.ContainsAny(word, "äöüß") {
			continue
		} else {
			filteredWordList = append(filteredWordList, word)
		}
	}

	return filteredWordList
}

func shuffleWord(word string) string {
	runes := []rune(word)

	rand.Shuffle(len(runes), func(i int, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	return string(runes)
}
