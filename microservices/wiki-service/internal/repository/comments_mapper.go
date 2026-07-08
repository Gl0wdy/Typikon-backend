package repository

import (
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/repository/models"
)

func toCommentModel(c *domain.Comment) *models.CommentModel {
	return &models.CommentModel{
		ID:            c.ID,
		ArticleID:     c.ArticleID,
		UserID:        c.UserID,
		Content:       c.Content,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
		LikesCount:    c.LikesCount,
		DislikesCount: c.DislikesCount,
		ParentID:      c.ParentID,
	}
}

func toCommentDomain(model *models.CommentModel) *domain.Comment {
	if model == nil {
		return nil
	}

	res := &domain.Comment{
		ID:            model.ID,
		ArticleID:     model.ArticleID,
		Content:       model.Content,
		UserID:        model.UserID,
		ParentID:      model.ParentID,
		CreatedAt:     model.CreatedAt,
		LikesCount:    model.LikesCount,
		DislikesCount: model.DislikesCount,
	}

	if len(model.Replies) > 0 {
		res.Replies = make([]*domain.Comment, len(model.Replies))
		for i, r := range model.Replies {
			res.Replies[i] = toCommentDomain(&r)
		}
	}

	return res
}

func toCommentVoteModel(v *domain.CommentVote) *models.CommentVoteModel {
	return &models.CommentVoteModel{
		CommentID: v.CommentID,
		UserID:    v.UserID,
		VoteType:  v.VoteType,
	}
}

func toCommentVoteDomain(m *models.CommentVoteModel) *domain.CommentVote {
	return &domain.CommentVote{
		CommentID: m.CommentID,
		UserID:    m.UserID,
		VoteType:  m.VoteType,
	}
}
