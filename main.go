package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

type WordResult struct {
	Word   string
	Length int
}

func findWords(letters string) ([]WordResult, error) {
	file, err := os.Open("words.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	words := make([]WordResult, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		wordLetters := strings.Split(scanner.Text(), "")
		otherLetters := make([]string, 0)
		currentSearchLetters := letters

		for _, letter := range wordLetters {
			if strings.Index(currentSearchLetters, letter) == -1 {
				otherLetters = append(otherLetters, letter)
			} else {
				currentSearchLetters = strings.Replace(currentSearchLetters, letter, "", 1)
			}
		}

		if len(otherLetters) == 0 {
			words = append(words, WordResult{Word: scanner.Text(), Length: len(scanner.Text())})
		}
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].Length > words[j].Length
	})

	return words, scanner.Err()
}

func wordsHandler(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	letters := r.URL.Query().Get("letters")

	matched, err := regexp.MatchString("^[a-zA-Z]+$", letters)

	if letters == "" || !matched {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Missing or incorrect letters.")
		return
	}

	words, err := findWords(letters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(words)
}

func main() {
	http.HandleFunc("/words", wordsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
