package server

import (
	"context"

	pb "github.com/kzz45/neverdown/pkg/authx/grpc/proto"
)

func (s *Server) KeepAlive(ctx context.Context, in *pb.Ping) (*pb.Pong, error) {
	return &pb.Pong{Hostname: &s.hostname}, nil
}

func (s *Server) RegisterRule(ctx context.Context, in *pb.RegisterRuleRequest) (*pb.RegisterRuleReply, error) {
	rep, err := s.aggregator.RegisterRule(ctx, in)
	if err != nil {
		failed := false
		reason := err.Error()
		return &pb.RegisterRuleReply{Result: &pb.StandardReply{
			Success: &failed,
			Reason:  &reason,
		}}, nil
	}
	return rep, nil
}

func (s *Server) ClusterRole(ctx context.Context, in *pb.ClusterRoleRequest) (*pb.ClusterRoleReply, error) {
	rep, err := s.aggregator.ClusterRole(ctx, in)
	if err != nil {
		failed := false
		reason := err.Error()
		return &pb.ClusterRoleReply{Result: &pb.StandardReply{
			Success: &failed,
			Reason:  &reason,
		}}, nil
	}
	return rep, nil
}

func (s *Server) ValidateRule(ctx context.Context, in *pb.ValidateRuleRequest) (*pb.ValidateRuleReply, error) {
	rep, err := s.aggregator.ValidateRule(ctx, in)
	if err != nil {
		failed := false
		reason := err.Error()
		return &pb.ValidateRuleReply{Result: &pb.StandardReply{
			Success: &failed,
			Reason:  &reason,
		}}, nil
	}
	return rep, nil
}
