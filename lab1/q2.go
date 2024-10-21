package main

import (
	"strings"

	"github.com/tjarratt/babble"
)

type ScrabbleScore struct{}

func (p *ScrabbleScore) Problem() string {
	return "(Scrabble score) Given a word, compute the Scrabble score for that word."
}

func (p *ScrabbleScore) Solve(args ...interface{}) interface{} {
	s := strings.ToUpper(args[0].(string))
	points := map[rune]int{
		'A': 1, 'E': 1, 'I': 1, 'O': 1, 'U': 1,
		'L': 1, 'N': 1, 'R': 1, 'S': 1, 'T': 1,
		'D': 2, 'G': 2,
		'B': 3, 'C': 3, 'M': 3, 'P': 3,
		'F': 4, 'H': 4, 'V': 4, 'W': 4, 'Y': 4,
		'K': 5,
		'J': 8, 'X': 8,
		'Q': 10, 'Z': 10,
	}

	point := 0
	for _, r := range s {
		point += points[r]
	}
	return point
}

func (p *ScrabbleScore) GenerateTestCases(cases int) [][]string {
	testCases := make([][]string, 0, cases)
	for i := 0; i < cases; i++ {
		testCase := p.generateOneCase()
		testCases = append(testCases, testCase)
	}
	return testCases
}

func (p *ScrabbleScore) generateOneCase() []string {
	babbler := babble.NewBabbler()
	babbler.Count = 1
	return []string{babbler.Babble()}
}
