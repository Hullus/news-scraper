package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
)

type Distributor struct {
	basePath     string
	newsFeedPath string
	regex        string
}

type NewsItem struct {
	Title string
	Url   string
}

func main() {
	newsDistributors := getDistributors()
	for _, distributor := range newsDistributors {
		news, err := scrapeNews(distributor)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(news)
	}
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
		linkEl := e.DOM.ParentFiltered("a")
		Url, exists := linkEl.Attr("href")
		if !exists {
			Url = "#"
		}

		fullUrl := Url
		if !strings.HasPrefix(Url, "http") {
			fullUrl = distributor.basePath + Url
		}

		newsItem := NewsItem{
			Title: strings.TrimSpace(e.Text),
			Url:   fullUrl,
		}

		newsItems = append(newsItems, newsItem)
	})

	err := c.Visit(distributor.newsFeedPath)
	if err != nil {
		return nil, err
	}

	return newsItems, nil
}

func getDistributors() []Distributor {
	newsDistributors := map[string]struct {
		newsFeedPath string
		regex        string
	}{
		"https://www.radiotavisupleba.ge": {
			newsFeedPath: "https://www.radiotavisupleba.ge/news",
			regex:        `.media-block__title`,
		},
		"https://netgazeti.ge/": {
			newsFeedPath: "https://netgazeti.ge/category/news/",
			regex:        `another-pattern`,
		},
		"https://formula.ge/": {
			newsFeedPath: "https://formulanews.ge/Category/All",
			regex:        `some-other-pattern`,
		},
		"https://www.imedi.ge/": {
			newsFeedPath: "https://imedinews.ge/ge/all-news",
			regex:        `pattern-for-imedi`,
		},
		"https://1tv.ge/": {
			newsFeedPath: "https://1tv.ge/akhali-ambebi/politika/",
			regex:        `pattern-for-1tv`,
		},
	}

	var distributors []Distributor

	for basePath, info := range newsDistributors {
		distributors = append(distributors, Distributor{
			basePath:     basePath,
			newsFeedPath: info.newsFeedPath,
			regex:        info.regex,
		})
	}

	return distributors
}
