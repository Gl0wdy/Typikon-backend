package usecase

import (
	"context"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
)

type CommentUseCase interface {
	Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	List(ctx context.Context, articleID string, limit, offset int) ([]*domain.Comment, error)
	GetByID(ctx context.Context, id string) (*domain.Comment, error)
	Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	Delete(ctx context.Context, id string) error
	SetVote(ctx context.Context, vote *domain.CommentVote) error
}

type commentUseCase struct {
	repo        domain.CommentRepository
	articleRepo domain.ArticleRepository
}

func NewCommentUseCase(repo domain.CommentRepository, articleRepo domain.ArticleRepository) CommentUseCase {
	return &commentUseCase{repo: repo, articleRepo: articleRepo}
}

func (u *commentUseCase) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	if err := comment.Validate(); err != nil {
		return nil, err
	}

	if _, err := u.articleRepo.GetByID(ctx, comment.ArticleID); err != nil {
		return nil, err
	}

	createdComment, err := u.repo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}
	return createdComment, nil
}

func (u *commentUseCase) List(ctx context.Context, articleID string, limit, offset int) ([]*domain.Comment, error) {
	comments, err := u.repo.List(ctx, articleID, limit, offset)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (u *commentUseCase) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	comment, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (u *commentUseCase) Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	if err := comment.Validate(); err != nil {
		return nil, err
	}

	if _, err := u.repo.GetByID(ctx, comment.ID); err != nil {
		return nil, err
	}

	return u.repo.Update(ctx, comment)
}

func (u *commentUseCase) Delete(ctx context.Context, id string) error {
	if _, err := u.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return u.repo.Delete(ctx, id)
}

func (u *commentUseCase) SetVote(ctx context.Context, vote *domain.CommentVote) error {
	if _, err := u.repo.GetByID(ctx, vote.CommentID); err != nil {
		return err
	}
	if vote.VoteType < 1 || vote.VoteType > 3 {
		return domain.ErrInvalidVoteType
	}
	return u.repo.SetVote(ctx, vote)
}
