package main

import (
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"strconv"
	"github.com/patrickmn/go-cache"
	"time"
	"github.com/lukaszromerowicz/dictionary-api/strutil"
)

var Cache *cache.Cache 

type WordsResponse struct {
	Count int
	Words []Word
}

type Word struct {
	Word   string
	Length int
	Meaning []string `json:",omitempty"`
}

type WordList struct {
	Words []Word
}

func mapWords() (WordList, error) {
	var dictionaryWords WordList

	file, err := os.Open("words.json")
	if err != nil {
		return dictionaryWords, err
	}
	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &dictionaryWords.Words)

	return dictionaryWords, nil
}

func (wordList *WordList) findWords(letters string, maxSize int) ([]Word, error) {
	words := make([]Word, 0)
	letters = strings.ToLower(letters)
	letters = strutil.SortString(letters)

	cachedWords, found := Cache.Get(letters)
	if found {
		return cachedWords.([]Word), nil
	}

	for _, word := range wordList.Words {
		wordLetters := strings.Split(word.Word, "")
		if len(wordLetters) < maxSize {
			continue
		}
		otherLetters := make([]string, 0)
		currentSearchLetters := letters
		
		for _, letter := range wordLetters {
			if strings.Index(currentSearchLetters, letter) == -1 {
				otherLetters = append(otherLetters, letter)
				break
			} else {
				currentSearchLetters = strings.Replace(currentSearchLetters, letter, "", 1)
			}
		}

		if len(otherLetters) == 0 {
			word.Length = len(word.Word)
			words = append(words, word)
		}
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].Length > words[j].Length
	})

	Cache.Set(letters, words, cache.NoExpiration)
	
	return words, nil
}

func (wordList *WordList) definitionHandler(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	requestedWord := r.URL.Query().Get("word")
	var foundWord *Word

	for _, word := range wordList.Words {
		if word.Word == requestedWord {
			word.Length = len(word.Word)
			foundWord = &word
			break
		}
	}

	if foundWord == nil {
		json.NewEncoder(w).Encode("Word not found.")
		return
	}
	
	json.NewEncoder(w).Encode(foundWord)
}

func (wordList *WordList) wordsHandler(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	letters := r.URL.Query().Get("letters")	
	meaning := r.URL.Query().Get("meaning") == "true"
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	matched, err := regexp.MatchString("^[a-zA-Z]+$", letters)

	if !matched {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Missing or incorrect letters.")
		return
	}

	words, err := wordList.findWords(letters, size)
	if limit > 0 && len(words) >= limit {
		words = words[0: limit]
	}

	response := WordsResponse{Words: words, Count: len(words)}

	if !meaning {
		for i := range response.Words {
			response.Words[i].Meaning = nil
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	wordList, err := mapWords()
	if err != nil {
		panic(err)
	}

	Cache = cache.New(5*time.Minute, 10*time.Minute)

	http.HandleFunc("/words", wordList.wordsHandler)
	http.HandleFunc("/definition", wordList.definitionHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
