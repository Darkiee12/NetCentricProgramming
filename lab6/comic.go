package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)
type Comic struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Like   string `json:"like"`
}

func (c *Comic) fromHTML(selection *goquery.Selection) *Comic {
	c.Title = strings.TrimSpace(selection.Find("p.subj").Text())
	c.Author = strings.TrimSpace(selection.Find("p.author").Text())
	c.Like = strings.TrimSpace(selection.Find("em.grade_num").Text())
	return c
}