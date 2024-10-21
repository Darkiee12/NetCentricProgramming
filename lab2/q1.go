package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func q1() {
	file, err := os.Open("q1.txt")
	if err != nil {
		fmt.Println("An error encountered:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var line string
	if scanner.Scan() {
		line = strings.ToLower(scanner.Text())
	}

	start := time.Now()
	freqConcurrent := countCharactersConcurrent(line, 100)
	durationConcurrent := time.Since(start).Nanoseconds()
	fmt.Println("Concurrent Frequency Count:", freqConcurrent)
	fmt.Printf("Concurrent Duration: %d nanoseconds\n", durationConcurrent)
}

func countCharactersConcurrent(s string, numSegments int) map[string]int {
	var wg sync.WaitGroup
	freqChannel := make(chan map[string]int)

	segmentLength := (len(s) + numSegments - 1) / numSegments

	for i := 0; i < numSegments; i++ {
		start := i * segmentLength
		end := start + segmentLength
		if end > len(s) {
			end = len(s)
		}

		if start >= len(s) {
			break
		}

		wg.Add(1)
		go func(segment string) {
			defer wg.Done()
			freq := make(map[string]int)
			for _, r := range segment {
				freq[string(r)]++
			}
			freqChannel <- freq
		}(s[start:end])
	}

	go func() {
		wg.Wait()
		close(freqChannel)
	}()

	totalFreq := make(map[string]int)
	for freq := range freqChannel {
		for k, v := range freq {
			totalFreq[k] += v
		}
	}

	return totalFreq
}
