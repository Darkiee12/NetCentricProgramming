package main

import "fmt"

type Solution interface {
	Problem() string
	Solve(args ...interface{}) interface{}
	generateOneCase() []string
	GenerateTestCases(count int) [][]string
}

type Problem struct {
	solution Solution
}

func splitArgument(args []string) []interface{} {
	result := make([]interface{}, len(args))
	for i, arg := range args {
		result[i] = arg
	}
	return result
}

func main() {
	problems := []Problem{
		{solution: &HammingDistance{}},
		{solution: &ScrabbleScore{}},
		{solution: &Luhn{}},
		{solution: &MineSweeper{}},
		{solution: &MatchingBracket{}},
	}

	cases := 10
	for _, problem := range problems {
		fmt.Printf("Problem: %s\n", problem.solution.Problem())
		testCases := problem.solution.GenerateTestCases(cases)
		for i, testCase := range testCases {
			args := splitArgument(testCase)
			result := problem.solution.Solve(args...)
			fmt.Printf("Test Case %d:%s\nResult: %v\n", i+1, testCase, result)
			fmt.Println("---")
		}
		fmt.Println("--------------------")
	}
}
