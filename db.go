package main

import (
	"fmt"
	"time"
)

func getLastArticles() ([]NewsItemDTO, error) {
	getArticlesQuery := `SELECT *
FROM (SELECT title, link, source, fetch_date,
             ROW_NUMBER() OVER (PARTITION BY source ORDER BY fetch_date DESC) AS rn
      FROM news_articles) AS subquery
WHERE rn <= 10;`
	rows, _ := db.Query(getArticlesQuery)

	var newsItems []NewsItemDTO

	for rows.Next() {
		var newsItem NewsItemDTO
		err := rows.Scan(&newsItem.Title, &newsItem.Url, &newsItem.Source, &newsItem.FetchDate, &newsItem.RN)
		if err != nil {
			return nil, err
		}
		newsItems = append(newsItems, newsItem)
	}

	for i, item := range newsItems {
		fmt.Print(i)
		fmt.Print("   ")
		fmt.Println(item.Title)
	}

	return newsItems, nil
}

func insertData(newsItems []NewsItem, provider string) error {
	//Get latest article from DB
	getLastArticleQuery := `SELECT * FROM news_articles WHERE source = $1 ORDER BY fetch_date LIMIT 1;`
	row := db.QueryRow(getLastArticleQuery, provider)
	lastArticle := NewsItem{}

	err := row.Scan(&lastArticle.Title, &lastArticle.Url, &lastArticle.Source, &lastArticle.ID, &lastArticle.FetchDate)
	if err != nil {
		return err
	}

	//Start insert batch preparation
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	insertArticleQuery := `INSERT INTO news_articles (title, link, source, fetch_date) VALUES ($1, $2, $3, $4)`
	stmt, err := tx.Prepare(insertArticleQuery)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, newsItem := range newsItems {
		if lastArticle.Title != "" && lastArticle.Title == newsItem.Title {
			break
		}
		_, err := stmt.Exec(newsItem.Title, newsItem.Url, provider, time.Now())
		if err != nil {
			fmt.Println(err)
		}
	}

	err = tx.Commit()

	return nil
}
