package main

import (
	"math/rand"
	"strconv"
)

const (
	rows      = 20
	cols      = 25
	numMines  = 99
	mineMark  = "*"
	emptyMark = "."
	delimeter = "\n"
)

type MineSweeper struct{}

func (p *MineSweeper) Problem() string {
	return "(Minesweeper) Add the mine counts to a completed Minesweeper board."
}

func (p *MineSweeper) Solve(args ...interface{}) interface{} {
	board := make([][]string, rows)
	for i := range board {
		board[i] = make([]string, cols)
	}
	i := 0
	for row := range rows {
		for col := range cols {
			cell := args[i].(string)
			if cell == delimeter {
				i++
			}
			board[row][col] = args[i].(string)
			i++
		}
	}
	for i := range board {
		for j := range board[i] {
			if board[i][j] == emptyMark {
				count := 0
				for x := -1; x <= 1; x++ {
					for y := -1; y <= 1; y++ {
						if i+x >= 0 && i+x < len(board) && j+y >= 0 && j+y < len(board[i]) && board[i+x][j+y] == mineMark {
							count++
						}
					}
				}
				if count > 0 {
					board[i][j] = strconv.Itoa(count)
				}
			}
		}
	}
	return output(board)
}

func output(board [][]string) string {
	var result string
	result += delimeter
	for _, row := range board {
		for _, cell := range row {
			result += cell + " "
		}
		result += delimeter
	}
	return result
}

func (p *MineSweeper) GenerateTestCases(cases int) [][]string {
	testCases := make([][]string, 0, cases)
	for i := 0; i < cases; i++ {
		testCase := p.generateOneCase()
		testCases = append(testCases, testCase)
	}
	return testCases
}

func (p *MineSweeper) generateOneCase() []string {
	board := generateMinefield(rows, cols, numMines)
	var result []string
	result = append(result, delimeter)
	for _, row := range board {
		result = append(result, row...)
		result = append(result, delimeter)
	}
	if len(result) > 0 {
		result = result[:len(result)-1]
	}

	return result
}

func generateMinefield(rows, cols, numMines int) [][]string {
	board := make([][]string, rows)
	for i := range board {
		board[i] = make([]string, cols)
		for j := range board[i] {
			board[i][j] = emptyMark
		}
	}
	minesPlaced := 0
	for minesPlaced < numMines {
		randRow := rand.Intn(rows)
		randCol := rand.Intn(cols)
		if board[randRow][randCol] == emptyMark {
			board[randRow][randCol] = mineMark
			minesPlaced++
		}
	}

	return board
}
