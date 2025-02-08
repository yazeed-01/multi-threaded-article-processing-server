package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mta/initializers"
	"mta/models"
)

func jsonToReader(data interface{}) io.Reader {
	b, _ := json.Marshal(data)
	return bytes.NewReader(b) // stream
}

func GetProcessedArticles(userID int) ([]models.Article, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"userId": userID,
			},
		},
	}

	res, err := initializers.ES.Search(
		initializers.ES.Search.WithIndex("articles"),
		initializers.ES.Search.WithBody(jsonToReader(query)),
		initializers.ES.Search.WithContext(context.Background()), // without timeout
	)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var searchResponse struct {
		Hits struct {
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	var articles []models.Article
	for _, hit := range searchResponse.Hits.Hits {
		var article models.Article
		if err := json.Unmarshal(hit.Source, &article); err != nil {
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}
