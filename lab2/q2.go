package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	maxCapacity   = 30
	totalStudents = 100
)

type Reader struct {
	ID   int
	Stay time.Duration
}

type Library struct {
	Capacity int
	Seats    chan *Reader
}

func NewLibrary(capacity int) *Library {
	return &Library{
		Capacity: capacity,
		Seats:    make(chan *Reader, capacity),
	}
}

func (l *Library) Enter(reader *Reader, logger *log.Logger) {
	select {
	case l.Seats <- reader:
		logger.Printf("Time %.0f: Student %d starts reading at the library for %v hours.\n", float64(time.Now().Unix()-startTime.Unix()), reader.ID, reader.Stay.Seconds())
		time.Sleep(reader.Stay)
		logger.Printf("Time %.0f: Student %d is leaving. Spent %v hours reading.\n", float64(time.Now().Unix()-startTime.Unix()), reader.ID, reader.Stay.Seconds())
		<-l.Seats
	default:
		logger.Printf("Time %.0f: Student %d is waiting to enter the library.\n", float64(time.Now().Unix()-startTime.Unix()), reader.ID)
		time.Sleep(1 * time.Second)
		l.Enter(reader, logger)
	}
}

func GenerateReaders(count int) []*Reader {
	readers := make([]*Reader, count)
	for i := 0; i < count; i++ {
		staySeconds := time.Duration(rand.Intn(4)+1) * time.Second
		readers[i] = &Reader{ID: i + 1, Stay: staySeconds}
	}
	return readers
}

var startTime time.Time

func q2() {
	library := NewLibrary(maxCapacity)
	readers := GenerateReaders(totalStudents)

	startTime = time.Now()

	fileName := "q2.txt"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	logger := log.New(file, "", 0)

	var wg sync.WaitGroup

	for _, reader := range readers {
		wg.Add(1)
		go func(r *Reader) {
			defer wg.Done()
			library.Enter(r, logger)
		}(reader)
	}

	wg.Wait()
	endTime := time.Now()

	logger.Printf("Time %.0f: No more students. Let's call it a day.\n", float64(endTime.Unix()-startTime.Unix()))
}
