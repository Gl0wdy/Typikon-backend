package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/repository/models"
)

type PostgresCommentRepo struct {
	db *gorm.DB
}

func NewPostgresCommentRepo(db *gorm.DB) *PostgresCommentRepo {
	return &PostgresCommentRepo{db: db}
}

func (r *PostgresCommentRepo) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	model := toCommentModel(comment)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return toCommentDomain(model), nil
}

func (r *PostgresCommentRepo) List(ctx context.Context, articleID string, limit, offset int) ([]*domain.Comment, error) {
	var commentModels []models.CommentModel

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	err := r.db.WithContext(ctx).
		Where("article_id = ? AND parent_id IS NULL", articleID).
		Limit(limit).
		Offset(offset).
		Preload("Replies").
		Find(&commentModels).
		Error

	if err != nil {
		return nil, err
	}

	comments := make([]*domain.Comment, len(commentModels))
	for i, model := range commentModels {
		comments[i] = toCommentDomain(&model)
	}
	return comments, nil
}

func (r *PostgresCommentRepo) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	var model models.CommentModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCommentNotFound
		}
		return nil, err
	}
	return toCommentDomain(&model), nil
}

func (r *PostgresCommentRepo) Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	model := toCommentModel(comment)

	result := r.db.WithContext(ctx).Model(&models.CommentModel{ID: comment.ID}).Updates(model)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrCommentNotFound
	}

	return toCommentDomain(model), nil
}

func (r *PostgresCommentRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.CommentModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrCommentNotFound
	}
	return nil
}

const (
	voteTypeLike    = 1
	voteTypeDislike = 2
	voteTypeUnvote  = 3
)

func columnForVoteType(voteType int32) string {
	if voteType == voteTypeDislike {
		return "dislikes_count"
	}
	return "likes_count"
}

func incrementColumn(tx *gorm.DB, commentID, column string, delta int) error {
	return tx.Model(&models.CommentModel{}).
		Where("id = ?", commentID).
		Update(column, gorm.Expr("GREATEST("+column+" + ?, 0)", delta)).Error
}

// TODO fix increment issues
// i don't really know where problem is tho...
func (r *PostgresCommentRepo) SetVote(ctx context.Context, vote *domain.CommentVote) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingVote models.CommentVoteModel
		err := tx.Where("comment_id = ? AND user_id = ?", vote.CommentID, vote.UserID).
			First(&existingVote).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		hasExistingVote := !errors.Is(err, gorm.ErrRecordNotFound)

		if vote.VoteType == voteTypeUnvote {
			if !hasExistingVote {
				return nil
			}
			if err := tx.Delete(&existingVote).Error; err != nil {
				return err
			}
			column := columnForVoteType(existingVote.VoteType)
			return incrementColumn(tx, vote.CommentID, column, -1)
		}

		if !hasExistingVote {
			newVote := toCommentVoteModel(vote)
			if err := tx.Create(newVote).Error; err != nil {
				return err
			}
			column := columnForVoteType(vote.VoteType)
			return incrementColumn(tx, vote.CommentID, column, 1)
		}

		if existingVote.VoteType == vote.VoteType {
			if err := tx.Delete(&existingVote).Error; err != nil {
				return err
			}
			column := columnForVoteType(vote.VoteType)
			return incrementColumn(tx, vote.CommentID, column, -1)
		}

		oldColumn := columnForVoteType(existingVote.VoteType)
		newColumn := columnForVoteType(vote.VoteType)

		existingVote.VoteType = vote.VoteType
		if err := tx.Save(&existingVote).Error; err != nil {
			return err
		}

		if err := incrementColumn(tx, vote.CommentID, oldColumn, -1); err != nil {
			return err
		}
		return incrementColumn(tx, vote.CommentID, newColumn, 1)
	})
}
