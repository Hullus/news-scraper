package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
)

type NewsItem struct {
	Title string
	Url   string
}

func main() {
	baseUrl := "https://www.radiotavisupleba.ge"
	news, err := scrapeNews(baseUrl)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range news {
		fmt.Printf("Title: %s\nURL: %s\n\n", item.Title, item.Url)
	}
}

func scrapeNews(baseUrl string) ([]NewsItem, error) {
	var newsItems []NewsItem

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %v failed with response: %v\nError: %v", r.Request.URL, r, err)
	})

	c.OnHTML(".media-block__title", func(e *colly.HTMLElement) {
		linkEl := e.DOM.ParentFiltered("a")
		Url, exists := linkEl.Attr("href")
		if !exists {
			Url = "#"
		}

		fullUrl := Url
		if !strings.HasPrefix(Url, "http") {
			fullUrl = baseUrl + Url
		}

		newsItem := NewsItem{
			Title: strings.TrimSpace(e.Text),
			Url:   fullUrl,
		}

		newsItems = append(newsItems, newsItem)
	})

	err := c.Visit(baseUrl + "/news")
	if err != nil {
		return nil, err
	}

	return newsItems, nil
}
