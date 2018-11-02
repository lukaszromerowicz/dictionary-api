package main

import (
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"fmt"
	"time"
	"strings"
)

type DictionaryEntry interface {
	keyword() string
	length() int
}

type WordsResponse struct {
	Count int
	Words []DictionaryEntry
}

type Word struct {
	Word   string
	Length int
	Meaning []string
}

type WordWithoutMeaning struct {
	Word string
	Length int
}

type WordList struct {
	Words []Word
}

func (word Word) keyword() string {
	return word.Word
}

func (word Word) length() int {
	return word.Length
}

func (word WordWithoutMeaning) keyword() string {
	return word.Word
}

func (word WordWithoutMeaning) length() int {
	return word.Length
}

func mapWords() (WordList, error) {
	var dictionaryWords WordList

	file, err := os.Open("words_new_parsed.json")
	if err != nil {
		return dictionaryWords, err
	}
	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &dictionaryWords.Words)

	return dictionaryWords, nil
}

func (wordList *WordList) findWords(letters string) ([]DictionaryEntry, error) {
	words := make([]DictionaryEntry, 0)
	start := time.Now()
	letters = strings.ToLower(letters)

	for _, word := range wordList.Words {
		wordLetters := strings.Split(word.Word, "")
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
		return words[i].length() > words[j].length()
	})
	
	fmt.Println(time.Since(start))

	return words, nil
}

func (wordList *WordList) wordsHandler(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	letters := r.URL.Query().Get("letters")	
	meaning := r.URL.Query().Get("meaning") == "true"
	matched, err := regexp.MatchString("^[a-zA-Z]+$", letters)

	if !matched {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Missing or incorrect letters.")
		return
	}

	words, err := wordList.findWords(letters)
	response := WordsResponse{Words: words, Count: len(words)}

	if !meaning {
		wordsWithoutMeaning := make([]DictionaryEntry, 0)

		for _, word := range words {
			wordsWithoutMeaning = append(wordsWithoutMeaning, WordWithoutMeaning{ Word: word.keyword(), Length: word.length()})
		}
		response = WordsResponse{Words: wordsWithoutMeaning, Count: len(wordsWithoutMeaning)}
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

	http.HandleFunc("/words", wordList.wordsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
