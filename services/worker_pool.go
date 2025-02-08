package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mta/initializers"
	"mta/models"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var ProcessingQueue chan models.Article
var resultChan chan string // channel to send the time taken for processing

// love worker pool pattern
func InitWorkerPool(numWorkers int) {
	ProcessingQueue = make(chan models.Article, 100)
	resultChan = make(chan string, 100)

	for i := 0; i < numWorkers; i++ {
		go worker(i + 1)
	}
}

func worker(id int) {
	for article := range ProcessingQueue {
		processArticle(article)
	}
}

func processArticle(article models.Article) {
	// start timer
	start := time.Now()

	updateStatus(article.UserID, article.ID, "PROCESSING")

	// insert article into Elasticsearch
	err := indexArticle(article)
	if err != nil {
		log.Printf("Error indexing article %d: %v", article.ID, err)
		updateStatus(article.UserID, article.ID, "FAILED")
		return
	}

	updateStatus(article.UserID, article.ID, "COMPLETED")

	// end timer
	duration := time.Since(start)

	// send result
	resultChan <- fmt.Sprintf("Article %d processed in %s", article.ID, duration)

}

func ProcessTimeTaken() {
	for result := range resultChan {
		fmt.Println(result)
	}
}

func indexArticle(article models.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "articles",
		DocumentID: strconv.Itoa(article.ID),
		Body:       bytes.NewReader(data),
		Refresh:    "true", // update index
	}

	res, err := req.Do(context.Background(), initializers.ES)
	if err != nil {
		return err
	}
	defer res.Body.Close() // close response body when done (memory leak)

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	return nil
}

func updateStatus(userID, articleID int, status string) {
	initializers.DB.Model(&models.Status{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Update("status", status)
}
