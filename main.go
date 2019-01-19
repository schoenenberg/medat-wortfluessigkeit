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
	"context"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

const (
	objectName = "words.csv"
)

var (
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
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	bucketName := os.Getenv("BUCKET_NAME")
	port := os.Getenv("PORT")
	_, debug := os.LookupEnv("DEBUG")

	// Check if environment variables are set.
	if projectID == "" {
		log.Fatalln("GOOGLE_CLOUD_PROJECT environment variable must be set.")
	} else if bucketName == "" {
		log.Fatalln("BUCKET_NAME environment variable must be set.")
	} else if port == "" {
		log.Fatalln("PORT environment variable must be set.")
	}

	// Creates a client.
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Read file from storage bucket
	b, err := readFileFromStorageBucket(ctx, client, bucketName, objectName)
	if err != nil {
		log.Fatalln(err)
	}

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

func readFileFromStorageBucket(ctx context.Context, client *storage.Client, bucket, object string) ([]byte, error) {
	// Create bucket reader
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	// .. and read
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
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
