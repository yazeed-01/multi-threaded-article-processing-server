package models

type Article struct {
	UserID  int    `json:"userId"`
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"body"`
}