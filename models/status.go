package models

import "time"

type Status struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    int       `gorm:"index"`
	ArticleID int       `gorm:"index"`
	Status    string    // PENDING, PROCESSING, COMPLETED, FAILED
	CreatedAt time.Time
	UpdatedAt time.Time
}