package main

import "math/rand"

type MatchingBracket struct{}

func (p *MatchingBracket) Problem() string {
	return "(Matching Brackets) Given a string containing brackets[ ], braces { }, parentheses ( ). Verify that all pairs are matched and nested correctly."
}

func (p *MatchingBracket) Solve(args ...interface{}) interface{} {
	s := args[0].(string)
	brackets := map[rune]rune{
		'[': ']',
		'{': '}',
		'(': ')',
	}
	stack := make([]rune, 0, len(s))
	for _, r := range s {
		if _, ok := brackets[r]; ok {
			stack = append(stack, r)
		} else if len(stack) == 0 || brackets[stack[len(stack)-1]] != r {
			return false
		} else {
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func (p *MatchingBracket) GenerateTestCases(cases int) [][]string {
	testCases := make([][]string, 0, cases)
	for i := 0; i < cases; i++ {
		testCase := p.generateOneCase()
		testCases = append(testCases, testCase)
	}
	return testCases
}

func (p *MatchingBracket) generateOneCase() []string {
	valid := rand.Intn(2) == 0
	numBrackets := rand.Intn(10) + 1
	if valid {
		return []string{generateValidBrackets(numBrackets)[rand.Intn(numBrackets)]}
	}
	return []string{generateInvalidBrackets(numBrackets)}
}

func generateInvalidBrackets(length int) string {
	brackets := []rune{'(', ')', '[', ']', '{', '}'}
	result := ""
	for i := 0; i < length; i++ {
		result += string(brackets[rand.Intn(len(brackets))])
	}
	return result
}

func generateValidBrackets(n int) []string {
	res := make([]string, 0)
	stack := make([]rune, 0)
	backtrack(0, 0, n, stack, &res)
	return res
}

func backtrack(open, close, n int, stack []rune, res *[]string) {
	if open == n && close == n {
		*res = append(*res, string(stack))
		return
	}

	if open < n {
		stack = append(stack, '(')
		backtrack(open+1, close, n, stack, res)
		stack = stack[:len(stack)-1]
	}

	if close < open {
		stack = append(stack, ')')
		backtrack(open, close+1, n, stack, res)
		stack = stack[:len(stack)-1]
	}
}
