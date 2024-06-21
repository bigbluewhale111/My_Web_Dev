package models

type Task struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type NewTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
