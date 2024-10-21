package main

import (
	"math/rand"
	"strconv"
	"strings"
)

type Luhn struct{}

func (p *Luhn) Problem() string {
	return "(Luhn) Given a number, determine whether or not it is valid per the Luhn formula."
}

func (p *Luhn) Solve(args ...interface{}) interface{} {
	number := args[0].(string)
	var sum int
	var odd bool
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if odd {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		odd = !odd
	}
	return sum%10 == 0
}

func GenerateValidCreditCard() string {
	cardNumber := make([]int, 15)
	for i := 0; i < 15; i++ {
		cardNumber[i] = rand.Intn(10)
	}

	checksum := calculateLuhnChecksum(cardNumber)
	cardNumber = append(cardNumber, checksum)

	return sliceToString(cardNumber)
}

func GenerateInvalidCreditCard() string {
	validCard := GenerateValidCreditCard()
	cardNumber := []rune(validCard)
	index := rand.Intn(len(cardNumber))
	cardNumber[index] = rune(rand.Intn(10) + '0')

	return string(cardNumber)
}

func GenerateRandomCreditCard() string {
	cardNumber := make([]int, 15)
	for i := 0; i < 15; i++ {
		cardNumber[i] = rand.Intn(10)
	}
	return sliceToString(cardNumber)
}

func calculateLuhnChecksum(numbers []int) int {
	sum := 0
	shouldDouble := true

	for i := len(numbers) - 1; i >= 0; i-- {
		digit := numbers[i]
		if shouldDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		shouldDouble = !shouldDouble
	}

	return (10 - (sum % 10)) % 10
}

func sliceToString(numbers []int) string {
	var sb strings.Builder
	for _, num := range numbers {
		sb.WriteString(strconv.Itoa(num))
	}
	return sb.String()
}

func (p *Luhn) GenerateTestCases(cases int) [][]string {
	testCases := make([][]string, 0, cases)
	for i := 0; i < cases; i++ {
		testCase := p.generateOneCase()
		testCases = append(testCases, testCase)
	}
	return testCases
}

func (p *Luhn) generateOneCase() []string {
	valid := rand.Intn(2) == 0
	if valid {
		return []string{GenerateValidCreditCard()}
	}
	return []string{GenerateInvalidCreditCard()}
}
