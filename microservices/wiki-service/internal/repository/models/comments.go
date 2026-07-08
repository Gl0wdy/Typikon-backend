package models

import (
	"time"
)

type CommentModel struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ArticleID     string    `gorm:"type:uuid;not null"`
	UserID        string    `gorm:"type:uuid;not null"`
	Content       string    `gorm:"type:text;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	ParentID      *string   `gorm:"index;type:uuid"`
	LikesCount    int       `gorm:"not null;default:0"`
	DislikesCount int       `gorm:"not null;default:0"`

	Replies []CommentModel `gorm:"foreignKey:ParentID"`
}

func (CommentModel) TableName() string {
	return "comments"
}

type CommentVoteModel struct {
	CommentID string `gorm:"primaryKey;type:uuid"`
	UserID    string `gorm:"primaryKey;type:uuid"`
	VoteType  int32  `gorm:"not null"` // LIKE=1, DISLIKE=2, REMOVE=0
}

func (CommentVoteModel) TableName() string {
	return "comment_votes"
}
