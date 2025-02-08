package repositories

import (
	"mta/initializers"
	"mta/models"
)

func GetProcessingStatus(userID int) ([]models.Status, error) {
	var statuses []models.Status
	result := initializers.DB.Where("user_id = ?", userID).Find(&statuses)
	return statuses, result.Error
}