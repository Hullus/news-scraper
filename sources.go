package main

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"strings"
)

func imedi(e *colly.HTMLElement, distributor Distributor, newsItems *[]NewsItem) {
	newsItem := NewsItem{
		Title: e.ChildText(".title"),
		Url:   e.Attr("href"),
	}
	err := fraudCheck(newsItem)
	if err != nil {
		return
	}

	*newsItems = append(*newsItems, newsItem)
}

func radioFreedom(e *colly.HTMLElement, distributor Distributor, newsItems *[]NewsItem) {
	title := e.ChildText(".media-block__title")
	Url := e.ChildAttr("a", "href")

	fullUrl := Url
	if !strings.HasPrefix(Url, "http") {
		fullUrl = distributor.basePath + Url
	}

	newsItem := NewsItem{
		Title: strings.TrimSpace(title),
		Url:   fullUrl,
	}
	err := fraudCheck(newsItem)
	if err != nil {
		return
	}

	*newsItems = append(*newsItems, newsItem)
}

func netGazeti(e *colly.HTMLElement, distributor Distributor, newsItems *[]NewsItem) {
	newsItem := NewsItem{
		Title: strings.TrimSpace(e.ChildText("h5 a")),
		Url:   e.ChildAttr("h5 a", "href"),
	}
	err := fraudCheck(newsItem)
	if err != nil {
		return
	}

	*newsItems = append(*newsItems, newsItem)
}

func formula(e *colly.HTMLElement, distributor Distributor, newsItems *[]NewsItem) {
	title := e.ChildText(".main__new__slider__desc a")
	Url := e.ChildAttr(".main__new__slider__desc a", "href")

	fullUrl := Url
	if !strings.HasPrefix(Url, "http") {
		fullUrl = distributor.basePath + Url
	}

	newsItem := NewsItem{
		Title: strings.TrimSpace(title),
		Url:   fullUrl,
	}
	err := fraudCheck(newsItem)
	if err != nil {
		return
	}

	*newsItems = append(*newsItems, newsItem)
}

func getDistributors() []Distributor {
	newsDistributors := map[string]struct {
		name               string
		newsFeedPath       string
		regex              string
		extractingFunction func(e *colly.HTMLElement, distributor Distributor, newsItems *[]NewsItem)
	}{
		"https://www.radiotavisupleba.ge": {
			name:               "radioFreedom",
			newsFeedPath:       "https://www.radiotavisupleba.ge/news",
			regex:              `.media-block`,
			extractingFunction: radioFreedom,
		},
		"https://netgazeti.ge/": {
			name:               "netgazeti",
			newsFeedPath:       "https://netgazeti.ge/category/news/",
			regex:              ".col-lg-4",
			extractingFunction: netGazeti,
		},
		"https://formulanews.ge": {
			name:               "formula",
			newsFeedPath:       "https://formulanews.ge/Category/All",
			regex:              `.news__box__card`,
			extractingFunction: formula,
		},
		"https://www.imedi.ge/": {
			name:               "imedi",
			newsFeedPath:       "https://imedinews.ge/ge/all-news",
			regex:              `.news-list .single-item`,
			extractingFunction: imedi,
		},
	}

	var distributors []Distributor

	for basePath, info := range newsDistributors {
		distributors = append(distributors, Distributor{
			name:               info.name,
			basePath:           basePath,
			newsFeedPath:       info.newsFeedPath,
			regex:              info.regex,
			extractingFunction: info.extractingFunction,
		})
	}

	return distributors
}

func fraudCheck(item NewsItem) error {
	if item.Title == "" || item.Url == "" || item.Title == "{{n.Title}}" {
		return errors.New("DON'T ADD THIS YOU DUMB FUCK")
	}

	return nil
}
