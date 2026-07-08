package wiki_grpc

import (
	"context"
	"time"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/usecase"
)

type CommentsServer struct {
	UnimplementedCommentsServiceServer
	usecase usecase.CommentUseCase
}

func NewCommentsServer(uc usecase.CommentUseCase) *CommentsServer {
	return &CommentsServer{usecase: uc}
}

func (s *CommentsServer) CreateComment(ctx context.Context, req *CreateCommentRequest) (*CommentResponse, error) {
	comment := &domain.Comment{
		Content:   req.GetContent(),
		ArticleID: req.GetArticleId(),
		UserID:    req.GetUserId(),
	}

	created, err := s.usecase.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &CommentResponse{
		Id:            created.ID,
		Content:       created.Content,
		ArticleId:     created.ArticleID,
		UserId:        created.UserID,
		CreatedAt:     created.CreatedAt.String(),
		LikesCount:    int32(created.LikesCount),
		DislikesCount: int32(created.DislikesCount),
	}, nil
}

func toCommentResponse(c *domain.Comment) *CommentResponse {
	if c == nil {
		return nil
	}

	resp := &CommentResponse{
		Id:            c.ID,
		ArticleId:     c.ArticleID,
		Content:       c.Content,
		UserId:        c.UserID,
		LikesCount:    int32(c.LikesCount),
		DislikesCount: int32(c.DislikesCount),
		CreatedAt:     c.CreatedAt.Format(time.RFC3339),
	}

	if c.ParentID != nil {
		resp.ParentId = c.ParentID
	}

	if len(c.Replies) > 0 {
		resp.Replies = make([]*CommentResponse, len(c.Replies))
		for i, reply := range c.Replies {
			resp.Replies[i] = toCommentResponse(reply) // Вызываем сами себя
		}
	}

	return resp
}

func (s *CommentsServer) ListComments(ctx context.Context, req *ListCommentsRequest) (*ListCommentsResponse, error) {
	comments, err := s.usecase.List(ctx, req.GetArticleId(), int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, err
	}

	tree := make([]*CommentResponse, len(comments))
	for i, c := range comments {
		tree[i] = toCommentResponse(c)
	}

	return &ListCommentsResponse{
		Comments: tree,
	}, nil
}

func (s *CommentsServer) UpdateComment(ctx context.Context, req *UpdateCommentRequest) (*CommentResponse, error) {
	comment := &domain.Comment{
		ID:      req.GetId(),
		Content: req.GetNewContent(),
	}

	updated, err := s.usecase.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &CommentResponse{
		Id:            updated.ID,
		Content:       updated.Content,
		ArticleId:     updated.ArticleID,
		UserId:        updated.UserID,
		CreatedAt:     updated.CreatedAt.String(),
		LikesCount:    int32(updated.LikesCount),
		DislikesCount: int32(updated.DislikesCount),
	}, nil
}

func (s *CommentsServer) DeleteComment(ctx context.Context, req *DeleteCommentRequest) (*DeleteCommentResponse, error) {
	err := s.usecase.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &DeleteCommentResponse{Success: true}, nil
}

func (s *CommentsServer) VoteComment(ctx context.Context, req *VoteCommentRequest) (*VoteCommentResponse, error) {
	err := s.usecase.SetVote(ctx, &domain.CommentVote{
		CommentID: req.GetCommentId(),
		UserID:    req.GetUserId(),
		VoteType:  int32(req.GetVoteType()),
	})
	if err != nil {
		return nil, err
	}

	return &VoteCommentResponse{Success: true}, nil
}
