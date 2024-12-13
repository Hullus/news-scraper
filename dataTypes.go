package main

import "github.com/gocolly/colly/v2"

type Distributor struct {
	name               string
	basePath           string
	newsFeedPath       string
	regex              string
	extractingFunction func(e *colly.HTMLElement, distributor Distributor, newsItems *[]NewsItem)
}

type NewsItem struct {
	Title     string `json:"title"`
	Url       string `json:"url"`
	ID        int    `json:"id"`
	Source    string `json:"source"`
	FetchDate string `json:"fetch_date"`
}

type NewsItemDTO struct {
	Title     string `json:"title"`
	Url       string `json:"url"`
	ID        int    `json:"id"`
	Source    string `json:"source"`
	FetchDate string `json:"fetch_date"`
	RN        int    `json:"row_number"`
}
