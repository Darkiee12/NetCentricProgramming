package main

import (
	"fmt"
	"math/rand"
)

type HammingDistance struct{}

func (p *HammingDistance) Problem() string {
	return "(Hamming) Calculate the Hamming Distance between two DNA strands."
}

func (p *HammingDistance) Solve(args ...interface{}) interface{} {
	s1 := args[0].(string)
	s2 := args[1].(string)
	if len(s1) != len(s2) {
		fmt.Println("Error: Strings must be of equal length")
		return nil
	}
	var distance int
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			distance++
		}
	}
	return distance
}

func (p *HammingDistance) GenerateTestCases(cases int) [][]string {
	testCases := make([][]string, 0, cases) // Initialize with length 0 and capacity 'cases'
	for i := 0; i < cases; i++ {
		testCase := p.generateOneCase()         // Assume this returns []string
		testCases = append(testCases, testCase) // Append the entire slice
	}
	return testCases
}

func (p *HammingDistance) generateOneCase() []string {
	bases := "ATCG"
	dna1 := make([]byte, 10)
	dna2 := make([]byte, 10)
	for j := range dna1 {
		dna1[j] = bases[rand.Intn(len(bases))]
	}
	for j := range dna2 {
		dna2[j] = bases[rand.Intn(len(bases))]
	}
	return []string{string(dna1), string(dna2)}
}
