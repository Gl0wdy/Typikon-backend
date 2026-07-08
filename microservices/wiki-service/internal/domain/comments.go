package domain

import (
	"context"
	"errors"
	"time"
)

type Comment struct {
	ID            string
	Content       string
	ArticleID     string
	ParentID      *string
	UserID        string
	LikesCount    int
	DislikesCount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Replies       []*Comment
}

var (
	ErrInvalidCommentContent = errors.New("invalid comment content: cannot be empty")
	ErrCommentNotFound       = errors.New("comment not found")
	ErrCommentContentTooLong = errors.New("invalid comment content: too long")
	ErrInvalidVoteType       = errors.New("invalid vote type: must be 0 (remove), 1 (like), or 2 (dislike)")
)

func (c *Comment) Validate() error {
	if c.Content == "" {
		return ErrInvalidCommentContent
	}
	if len(c.Content) > 500 {
		return ErrCommentContentTooLong
	}
	return nil
}

type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) (*Comment, error)
	GetByID(ctx context.Context, id string) (*Comment, error)
	List(ctx context.Context, articleID string, limit, offset int) ([]*Comment, error)
	Update(ctx context.Context, comment *Comment) (*Comment, error)
	Delete(ctx context.Context, id string) error
	SetVote(ctx context.Context, vote *CommentVote) error
	DeleteVote(ctx context.Context, commentID string, userID string) error
}

type CommentVote struct {
	CommentID string
	UserID    string
	VoteType  int32 // LIKE=1, DISLIKE=2, REMOVE=0
}

type CommentVoteRepository interface {
	GetVote(ctx context.Context, commentID string, userID string) (*CommentVote, error)
	SetVote(ctx context.Context, vote *CommentVote) error
	DeleteVote(ctx context.Context, commentID string, userID string) error
}
