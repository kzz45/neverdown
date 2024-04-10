package grpc

import (
	"context"

	pb "github.com/kzz45/neverdown/pkg/authx/grpc/proto"
)

type Api interface {
	RegisterRule(ctx context.Context, in *pb.RegisterRuleRequest) (*pb.RegisterRuleReply, error)
	ClusterRole(ctx context.Context, in *pb.ClusterRoleRequest) (*pb.ClusterRoleReply, error)
	ValidateRule(ctx context.Context, in *pb.ValidateRuleRequest) (*pb.ValidateRuleReply, error)
}
