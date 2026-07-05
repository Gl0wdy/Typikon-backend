package models

import "time"

type ArticleModel struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title      string    `gorm:"type:varchar(255);not null"`
	Content    string    `gorm:"type:text;not null"`
	CategoryID string    `gorm:"type:uuid;not null"`
	UserID     string    `gorm:"type:uuid;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (ArticleModel) TableName() string {
	return "articles"
}
