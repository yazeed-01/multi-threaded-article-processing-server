package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mta/initializers"
	"mta/models"

	"github.com/gin-gonic/gin"
)

func jsonToReader(data interface{}) (*bytes.Reader, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return bytes.NewReader(jsonData), nil
}

func getArticlesFromElasticSearch(query map[string]interface{}) ([]models.Article, error) {
	bodyReader, err := jsonToReader(query)
	if err != nil {
		return nil, err
	}

	res, err := initializers.ES.Search(
		initializers.ES.Search.WithIndex("articles"),
		initializers.ES.Search.WithBody(bodyReader),
		initializers.ES.Search.WithContext(context.Background()),
	)

	if err != nil {
		return nil, fmt.Errorf("elasticsearch search failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("elasticsearch error response: %s", err)
		}
		return nil, fmt.Errorf("elasticsearch error: %v", e)
	}

	var searchResponse struct {
		Hits struct {
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var articles []models.Article
	for _, hit := range searchResponse.Hits.Hits {
		var article models.Article
		if err := json.Unmarshal(hit.Source, &article); err != nil {
			continue // skip invalid articles
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func SearchSubstringArticles(query string) ([]models.Article, error) {
	elasticQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]interface{}{
				"body": map[string]interface{}{
					"value":            "*" + query + "*",
					"case_insensitive": true,
				},
			},
		},
	}
	return getArticlesFromElasticSearch(elasticQuery)
}

func SearchExactArticles(query string) ([]models.Article, error) {
	elasticQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"body": query,
			},
		},
	}
	return getArticlesFromElasticSearch(elasticQuery)
}

func SearchOnArticles(c *gin.Context) {
	query := c.Param("query")
	searchType := c.Param("type")

	// 0 -> exact, 1 -> substring

	if searchType != "0" && searchType != "1" {
		c.JSON(400, gin.H{"error": "invalid search type, use '0' for exact search and  '1' for substring search"})
		return
	}

	var articles []models.Article
	var err error

	switch searchType {
	case "0":
		articles, err = SearchExactArticles(query)
	case "1":
		articles, err = SearchSubstringArticles(query)
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, articles)
}
