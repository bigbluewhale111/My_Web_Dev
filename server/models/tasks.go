package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ID          uint32 `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      uint8  `json:"status"`
	DueDate     uint64 `json:"due_date"`
	AuthorID    uint32 `json:"author_id"`
}

type NewTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      uint8  `json:"status"`
	DueDate     uint64 `json:"due_date"`
}
