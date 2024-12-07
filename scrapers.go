package main

import (
	"github.com/gocolly/colly/v2"
	"strings"
)

func imedi(e *colly.HTMLElement, distributor Distributor, newsItems *[]NewsItem) {
	newsItem := NewsItem{
		Title: e.ChildText(".title"),
		Url:   e.Attr("href"),
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

	*newsItems = append(*newsItems, newsItem)
}

func netGazeti(e *colly.HTMLElement, distributor Distributor, newsItems *[]NewsItem) {
	newsItem := NewsItem{
		Title: strings.TrimSpace(e.ChildText("h5 a")),
		Url:   e.ChildAttr("h5 a", "href"),
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

	*newsItems = append(*newsItems, newsItem)
}
