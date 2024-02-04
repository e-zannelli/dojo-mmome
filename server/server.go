package main

import (
	"encoding/json"
	"log"
	"math/rand"
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
	playableRunes := []rune("🏴💀🔥🎉🚀🤡")
	gameIsRunning := true
	code := choices(playableRunes, 5)
	log.Println("Chosen code: ", string(code))

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		if gameIsRunning {
			w.WriteHeader(http.StatusBadRequest)

			resp, err := json.Marshal(map[string]string{"error": "game is already running"})
			if err != nil {
				log.Println("Error: ", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write(resp)
			return
		}

		code = choices(playableRunes, 5)
		log.Println("New code: ", string(code))
		gameIsRunning = true

		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/guess/", func(w http.ResponseWriter, r *http.Request) {
		correct := 0
		misplaced := 0

		if !gameIsRunning {
			w.WriteHeader(http.StatusBadRequest)
			resp, err := json.Marshal(map[string]string{"error": "game is not running, POST to /new to start a new game"})
			if err != nil {
				log.Println("Error: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(resp)
			return
		}

		guess := []rune(r.URL.Path[7:])
		log.Printf("%s guessed %s", r.RemoteAddr, string(guess))

		w.Header().Set("Content-Type", "application/json")
		if len(guess) != len(code) {
			w.WriteHeader(http.StatusBadRequest)
			resp, err := json.Marshal(map[string]string{"error": "Invalid guess" + string(guess)})
			if err != nil {
				log.Println("Error: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(resp)
			return
		}

		for i, r := range guess {
			if r == code[i] {
				correct++
			} else {
				for _, s := range code {
					if r == s {
						misplaced++
						break
					}
				}
			}
		}

		resp, err := json.Marshal(map[string]int{
			"correct":   correct,
			"misplaced": misplaced,
		})
		if err != nil {
			log.Println("Error: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(resp)

		if correct == len(code) {
			log.Printf("%s won, game is done.", r.RemoteAddr)
			gameIsRunning = false
		}
	})

	log.Println("Starting server on port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Error: ", err)
	}
}
