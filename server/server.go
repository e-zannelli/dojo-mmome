package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"mmome/domain"
	"net/http"
)

func choices(runes []rune, length int) []rune {
	choice := make([]rune, length)
	for i := 0; i < length; i++ {
		choice[i] = runes[rand.Intn(len(runes))]
	}
	return choice
}

func main() {
	game := domain.NewGame()
	game.Reset()

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		err := game.Reset()
		if err != nil {
			serverError(err, w, http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/guess/", func(w http.ResponseWriter, r *http.Request) {
		guess := []rune(r.URL.Path[7:])
		log.Printf("%s guessed %s", r.RemoteAddr, string(guess))
		w.Header().Set("Content-Type", "application/json")

		correct, misplaced, err := game.Guess(guess)
		if err != nil {
			serverError(err, w, http.StatusBadRequest)
			return
		}

		resp, err := json.Marshal(map[string]int{
			"correct":   correct,
			"misplaced": misplaced,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(resp)
	})

	log.Println("Starting server on port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Error: ", err)
	}
}

func serverError(err error, w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	resp, err := json.Marshal(map[string]string{"error": err.Error()})
	if err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}
