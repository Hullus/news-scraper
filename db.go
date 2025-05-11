package main

import (
	"fmt"
	"time"
)

func getLastArticles() (map[string][]NewsItemDTO, error) {
	getArticlesQuery := `
		SELECT title, link, source, fetch_date, rn
		FROM (
			SELECT title, link, source, fetch_date,
				   ROW_NUMBER() OVER (PARTITION BY source ORDER BY fetch_date DESC) AS rn
			FROM news_articles
		) AS subquery
		WHERE rn <= 10;`

	rows, err := db.Query(getArticlesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for last articles: %w", err)
	}
	defer rows.Close()

	newsItemsMap := make(map[string][]NewsItemDTO)
	for rows.Next() {
		var newsItem NewsItemDTO
		err := rows.Scan(&newsItem.Title, &newsItem.Url, &newsItem.Source, &newsItem.FetchDate, &newsItem.RN)
		if err != nil {
			return nil, fmt.Errorf("failed to scan news item: %w", err)
		}
		newsItemsMap[newsItem.Source] = append(newsItemsMap[newsItem.Source], newsItem)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over news item rows: %w", err)
	}

	return newsItemsMap, nil
}

func insertData(newsItems []NewsItem, provider string) error {
	getLastArticleQuery := `
		SELECT title, link, source, fetch_date 
		FROM news_articles 
		WHERE source = $1 
		ORDER BY fetch_date DESC 
		LIMIT 1;`

	rows, err := db.Query(getLastArticleQuery, provider)
	if err != nil {
		return fmt.Errorf("could not query last article for provider %s: %w", provider, err)
	}
	defer rows.Close()

	var lastArticle NewsItem

	if rows.Next() {
		err = rows.Scan(&lastArticle.Title, &lastArticle.Url, &lastArticle.Source, &lastArticle.FetchDate)
		if err != nil {
			return fmt.Errorf("could not scan last article for provider %s: %w", provider, err)
		}
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("error after checking last article for provider %s: %w", provider, err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin database transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	insertArticleQuery := `INSERT INTO news_articles (title, link, source, fetch_date) VALUES ($1, $2, $3, $4)`
	stmt, err := tx.Prepare(insertArticleQuery)
	if err != nil {
		return fmt.Errorf("could not prepare insert statement: %w", err)
	}
	defer stmt.Close()

	insertTime := time.Now()
	var itemsToInsert []NewsItem

	if lastArticle.Title == "" {
		itemsToInsert = newsItems
	} else {
		for _, newsItem := range newsItems {
			if lastArticle.Title == newsItem.Title && lastArticle.Url == newsItem.Url {
				break
			}
			itemsToInsert = append(itemsToInsert, newsItem)
		}
	}

	if len(itemsToInsert) == 0 {
		fmt.Printf("No new articles to insert for source %s.\n", provider)
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("could not commit transaction (even if empty): %w", err)
		}
		return nil
	}

	fmt.Printf("Inserting %d new articles for source %s.\n", len(itemsToInsert), provider)
	for _, newsItem := range itemsToInsert {
		_, execErr := stmt.Exec(newsItem.Title, newsItem.Url, provider, insertTime)
		if execErr != nil {
			err = fmt.Errorf("failed to execute insert for article '%s': %w", newsItem.Title, execErr)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
