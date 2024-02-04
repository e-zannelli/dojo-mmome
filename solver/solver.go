package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Clues struct {
	Correct   int `json:"correct"`
	Misplaced int `json:"misplaced"`
}

func main() {
	var playable = []rune("🏴💀🔥🎉🚀🤡")
	var guess = make([]rune, 5)
	var clues Clues

	// Build all possible combinations
	for _, r1 := range playable {
		for _, r2 := range playable {
			for _, r3 := range playable {
				for _, r4 := range playable {
					for _, r5 := range playable {
						guess[0] = r1
						guess[1] = r2
						guess[2] = r3
						guess[3] = r4
						guess[4] = r5
						resp, err := http.Get("http://localhost:8081/guess/" + string(guess))
						if err != nil {
							fmt.Println("Error: ", err)
						}

						content, err := io.ReadAll(resp.Body)
						if err != nil {
							panic(err)
						}
						resp.Body.Close()
						fmt.Printf("Guess: %s, resp: %s\n", string(guess), string(content))

						json.Unmarshal(content, &clues)
						if clues.Correct == 5 {
							fmt.Println("Won! Code was: ", string(guess))
							return
						}
					}
				}
			}
		}
	}
}
