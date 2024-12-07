package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"net/http"
	"sync"
)

type Distributor struct {
	name               string
	basePath           string
	newsFeedPath       string
	regex              string
	extractingFunction func(e *colly.HTMLElement, distributor Distributor, newsItems *[]NewsItem)
}

type NewsItem struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /news", func(w http.ResponseWriter, r *http.Request) {
		output, _ := parallelScraper()
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(output)
		if err != nil {
			fmt.Println(err)
		}
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}

func parallelScraper() (map[string][]NewsItem, error) {
	var wg sync.WaitGroup
	result := make(map[string][]NewsItem)

	newsDistributors := getDistributors()
	for _, distributor := range newsDistributors {
		wg.Add(1)
		go func(distributor Distributor) {
			defer wg.Done()
			news, err := scrapeNews(distributor)
			if err != nil {
				log.Fatal(err)
			}
			result[distributor.name] = news
		}(distributor)
	}

	wg.Wait()

	return result, nil
}

func scrapeNews(distributor Distributor) ([]NewsItem, error) {
	var newsItems []NewsItem

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %v failed with response: %v\nError: %v", r.Request.URL, r, err)
	})

	c.OnHTML(distributor.regex, func(e *colly.HTMLElement) {
		distributor.extractingFunction(e, distributor, &newsItems)
	})

	err := c.Visit(distributor.newsFeedPath)
	if err != nil {
		return nil, err
	}

	return newsItems, nil
}
