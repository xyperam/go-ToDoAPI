package models

type Task struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (Task) TableName() string {
	return "tasks" // nama tabel di PostgreSQL
}
