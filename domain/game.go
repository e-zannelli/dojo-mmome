package domain

import (
	"fmt"
	"math/rand"
)

const codeLength = 5

type Game struct {
	isRunning     bool
	playableRunes []rune
	code          []rune
}

func NewGame() Game {
	return Game{
		playableRunes: []rune("123456"),
	}
}

func (g *Game) Reset() error {
	if g.isRunning {
		return fmt.Errorf("game is already running")
	}

	g.code = choices(g.playableRunes, codeLength)
	g.isRunning = true

	return nil
}

func (g *Game) Guess(guess []rune) (int, int, error) {
	if !g.isRunning {
		return 0, 0, fmt.Errorf("game is not running")
	}
	if len(guess) != len(g.code) {
		return 0, 0, fmt.Errorf("invalid guess %s", string(guess))
	}

	correct := 0
	misplaced := 0
	checked := make([]bool, len(g.code))

	for i, r := range guess {
		if r == g.code[i] {
			checked[i] = true
			correct++
		}
	}
	for _, r := range guess {
		for j, s := range g.code {
			if checked[j] {
				continue
			}
			if r == s {
				checked[j] = true
				misplaced++
				break
			}
		}
	}

	if correct == len(g.code) {
		g.isRunning = false
	}

	return correct, misplaced, nil
}

func choices(runes []rune, length int) []rune {
	if len(runes) == 0 {
		return []rune{}
	}

	choice := make([]rune, length)
	for i := 0; i < length; i++ {
		choice[i] = runes[rand.Intn(len(runes))]
	}
	return choice
}
