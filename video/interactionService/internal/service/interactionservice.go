package service

import (
	"context"

	pb "west2-video/api/interaction/v1"
)

type InteractionServiceService struct {
	pb.UnimplementedInteractionServiceServer
}

func NewInteractionServiceService() *InteractionServiceService {
	return &InteractionServiceService{}
}

func (s *InteractionServiceService) LikeAction(ctx context.Context, req *pb.LikeActionRequest) (*pb.LikeActionReply, error) {
	return &pb.LikeActionReply{}, nil
}
func (s *InteractionServiceService) LikeList(ctx context.Context, req *pb.LikeListRequest) (*pb.LikeListReply, error) {
	return &pb.LikeListReply{}, nil
}
func (s *InteractionServiceService) CommentAction(ctx context.Context, req *pb.CommentActionRequest) (*pb.CommentActionReply, error) {
	return &pb.CommentActionReply{}, nil
}
func (s *InteractionServiceService) CommentList(ctx context.Context, req *pb.CommentListRequest) (*pb.CommentListReply, error) {
	return &pb.CommentListReply{}, nil
}
func (s *InteractionServiceService) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentReply, error) {
	return &pb.DeleteCommentReply{}, nil
}





