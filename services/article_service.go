package services

import (
	"encoding/json"
	"fmt"
	"io"
	"mta/initializers"
	"mta/models"
	"net/http"
)

func ProcessUserArticles(userID int) {
	articles, err := fetchArticles(userID)
	if err != nil {
		return
	}

	// Add articles to processing queue
	for _, article := range articles {
		createPendingStatus(userID, article.ID)
		ProcessingQueue <- article
	}
	// timer
	go ProcessTimeTaken()
}

func fetchArticles(userID int) ([]models.Article, error) {
	resp, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/posts?userId=%d", userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var articles []models.Article
	err = json.Unmarshal(body, &articles)
	return articles, err
}

func createPendingStatus(userID, articleID int) {
	status := models.Status{
		UserID:    userID,
		ArticleID: articleID,
		Status:    "PENDING",
	}
	initializers.DB.Create(&status)
}
