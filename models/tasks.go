package models

import "time"

type Task struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	UserID      int       `json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"` // Relasi ke User
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Task) TableName() string {
	return "tasks" // nama tabel di PostgreSQL
}
