package controllers

import (
	"net/http"
	"strconv"

	"mta/repositories"
	"mta/services"

	"github.com/gin-gonic/gin"
)

func Process(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	go services.ProcessUserArticles(id)

	c.JSON(http.StatusAccepted, gin.H{"message": "processing started"})
}

func GetArticles(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	articles, err := repositories.GetProcessedArticles(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func GetStatus(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	statuses, err := repositories.GetProcessingStatus(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, statuses)
}