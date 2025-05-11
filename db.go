package main

import (
	"fmt"
	"time"
)

func getLastArticles() (map[string][]NewsItemDTO, error) {
	getArticlesQuery := `SELECT *
FROM (SELECT title, link, source, fetch_date,
             ROW_NUMBER() OVER (PARTITION BY source ORDER BY fetch_date DESC) AS rn
      FROM news_articles) AS subquery
WHERE rn <= 10;`
	rows, _ := db.Query(getArticlesQuery)

	newsItems := make(map[string][]NewsItemDTO)
	for rows.Next() {
		var newsItem NewsItemDTO
		err := rows.Scan(&newsItem.Title,
			&newsItem.Url,
			&newsItem.Source,
			&newsItem.FetchDate,
			&newsItem.RN)

		if err != nil {
			return nil, err
		}

		newsItems[newsItem.Source] = append(newsItems[newsItem.Source], newsItem)
	}

	return newsItems, nil
}

func insertData(newsItems []NewsItem, provider string) error {
	//Get latest article from DB
	getLastArticleQuery := `SELECT * FROM news_articles WHERE source = $1 ORDER BY fetch_date LIMIT 10;`
	rows, _ := db.Query(getLastArticleQuery, provider)
	lastArticle := NewsItem{}

	err := rows.Scan(&lastArticle.Title,
		&lastArticle.Url,
		&lastArticle.Source,
		&lastArticle.ID,
		&lastArticle.FetchDate)
	if err != nil {
		return err
	}

	//Start insert batch preparation
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("could not start DB insert: %s\n", err)
	}

	insertArticleQuery := `INSERT INTO news_articles (title, link, source, fetch_date) VALUES ($1, $2, $3, $4)`
	stmt, err := tx.Prepare(insertArticleQuery)
	if err != nil {
		return fmt.Errorf("could not prepare insert statement: %w", err)
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

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
