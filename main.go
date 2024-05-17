package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	won := 0
	totalGames := 10000
	results := make(chan bool, totalGames)

	var wg sync.WaitGroup

	for i := 0; i < totalGames; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := startGame(AllWords)
			results <- result
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result {
			won++
		}
	}

	fmt.Printf("won: %v\n", won)

	duration := time.Since(start)
	fmt.Printf("Execution time: %s\n", duration)
}

func startGame(words []string) bool {
	remainingWords := words
	guess := "crane"
	attempts := 1

	randomIndex := rand.Intn(len(remainingWords))
	answer := words[randomIndex]

	for attempts <= 6 {

		checkedGuess := checkGuess(&guess, &answer)
		filterWords(&remainingWords, guess, checkedGuess)

		if guess == answer {
			return true
		}
		attempts++
		guess = remainingWords[0]
	}

	return false
}

func checkGuess(guess *string, answer *string) *[]rune {
	feedback := make([]rune, 5)
	answerRunes := []rune(*answer)
	guessRunes := []rune(*guess)

	for i := range guessRunes {
		if guessRunes[i] == answerRunes[i] {
			feedback[i] = 'G'
		}
	}

	for i := range guessRunes {
		if feedback[i] == 'G' {
			continue
		}

		for j := range answerRunes {
			if guessRunes[i] == answerRunes[j] {
				feedback[i] = 'Y'
				break
			}
		}

		if feedback[i] == 0 {
			feedback[i] = 'B'
		}
	}

	return &feedback
}

func filterWords(remainingWords *[]string, guess string, filter *[]rune) {

	filterWords := make([]string, 0)

	for _, word := range *remainingWords {
		filterRunes := []rune(*filter)
		match := true
		for i := range filterRunes {
			if (filterRunes[i] == 'G' && word[i] != guess[i]) ||
				(filterRunes[i] == 'Y' && (word[i] == guess[i] || !strings.Contains(word, string(guess[i])))) ||
				(filterRunes[i] == 'B' && strings.Contains(word, string(guess[i]))) {
				match = false
				break
			}
		}

		if match {
			filterWords = append(filterWords, word)
		}
	}

	*remainingWords = filterWords
}
