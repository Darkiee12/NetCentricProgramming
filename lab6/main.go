package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	scraper := WebScraper{
		Genres:  []string{"action", "romance", "comedy", "fantasy", "drama", "slice_of_life", "sci-fi", "supernatural", "historical", "heartwarming"},
		BaseURL: "https://www.webtoons.com/en/genres/",
	}

	scraper.scrape()

	now := time.Now().UnixNano()
	filename := fmt.Sprintf("webtoons_%d.json", now)

	err := scraper.outputJSON(filename)
	if err != nil {
		log.Fatalf("Failed to write JSON file: %v", err)
	}

	fmt.Printf("Scraping completed and data saved to %s\n", filename)
}
