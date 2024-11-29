package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type WebScraper struct {
	Genres  []string
	BaseURL string
	Data    map[string][]Comic
}

func (ws *WebScraper) fetchPage(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch page, status code: %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (ws *WebScraper) extractComics(doc *goquery.Document) []Comic {
	var comics []Comic
	doc.Find("ul > li").Each(func(index int, item *goquery.Selection) {
		comic := new(Comic).fromHTML(item)
		comics = append(comics, *comic)
	})
	return comics
}

func (ws *WebScraper) scrape() {
	ws.Data = make(map[string][]Comic)
	for _, genre := range ws.Genres {
		url := ws.BaseURL + genre
		doc, err := ws.fetchPage(url)
		if err != nil {
			log.Printf("Failed to fetch page for genre %s: %v", genre, err)
			continue
		}
		ws.Data[genre] = ws.extractComics(doc)
	}
}

func (ws *WebScraper) outputJSON(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(ws.Data)
	if err != nil {
		return err
	}
	return nil
}
