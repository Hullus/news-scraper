package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"sync"
)

type Distributor struct {
	basePath           string
	newsFeedPath       string
	regex              string
	extractingFunction func(e *colly.HTMLElement, distributor Distributor, newsItems *[]NewsItem)
}

type NewsItem struct {
	Title string
	Url   string
}

func main() {
	var wg sync.WaitGroup

	newsDistributors := getDistributors()
	for _, distributor := range newsDistributors {
		wg.Add(1)

		go func(distributor Distributor) {
			defer wg.Done()
			news, err := scrapeNews(distributor)
			fmt.Println(distributor)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(distributor.basePath)
			fmt.Println(news)
		}(distributor)
	}

	wg.Wait()

	fmt.Println("Done")
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

func getDistributors() []Distributor {
	newsDistributors := map[string]struct {
		newsFeedPath       string
		regex              string
		extractingFunction func(e *colly.HTMLElement, distributor Distributor, newsItems *[]NewsItem)
	}{
		"https://www.radiotavisupleba.ge": {
			newsFeedPath:       "https://www.radiotavisupleba.ge/news",
			regex:              `.media-block`,
			extractingFunction: radioFreedom,
		},
		"https://netgazeti.ge/": {
			newsFeedPath:       "https://netgazeti.ge/category/news/",
			regex:              ".col-lg-4",
			extractingFunction: netGazeti,
		},
		"https://formula.ge/": {
			newsFeedPath:       "https://formulanews.ge/Category/All",
			regex:              `.news__box__card`,
			extractingFunction: formula,
		},
		"https://www.imedi.ge/": {
			newsFeedPath:       "https://imedinews.ge/ge/all-news",
			regex:              `.news-list .single-item`,
			extractingFunction: imedi,
		},
	}

	var distributors []Distributor

	for basePath, info := range newsDistributors {
		distributors = append(distributors, Distributor{
			basePath:           basePath,
			newsFeedPath:       info.newsFeedPath,
			regex:              info.regex,
			extractingFunction: info.extractingFunction,
		})
	}

	return distributors
}
