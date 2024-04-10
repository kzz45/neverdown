package aggregator

import (
	"context"
	"fmt"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	pb "github.com/kzz45/neverdown/pkg/authx/grpc/proto"
	"github.com/kzz45/neverdown/pkg/authx/rbac/admin"

	"k8s.io/apimachinery/pkg/api/errors"
)

type Aggregator struct {
	apps admin.App
}

func New(apps admin.App) *Aggregator {
	return &Aggregator{
		apps: apps,
	}
}

func (ag *Aggregator) RegisterRule(ctx context.Context, in *pb.RegisterRuleRequest) (*pb.RegisterRuleReply, error) {
	if err := ag.apps.Validate(*in.AppMeta.AppId, *in.AppMeta.AppSecret); err != nil {
		return nil, err
	}
	ga, err := ag.apps.GenericApp(*in.AppMeta.AppId)
	if err != nil {
		return nil, err
	}
	rules := &rbacv1.GroupVersionKindRuleList{}
	if err = rules.Unmarshal(in.GroupVersionKindRuleList); err != nil {
		return nil, err
	}
	for _, v := range rules.Items {
		if err = ga.Kind().Update(ctx, &v); err != nil {
			if errors.IsNotFound(err) {
				if err = ga.Kind().Create(ctx, &v); err != nil {
					if !errors.IsAlreadyExists(err) {
						return nil, err
					}
				}
			} else {
				return nil, err
			}
		}
	}
	success := true
	return &pb.RegisterRuleReply{Result: &pb.StandardReply{
		Success: &success,
		Reason:  nil,
	}}, nil
}

func (ag *Aggregator) ClusterRole(ctx context.Context, in *pb.ClusterRoleRequest) (*pb.ClusterRoleReply, error) {
	if err := ag.apps.Validate(*in.AppMeta.AppId, *in.AppMeta.AppSecret); err != nil {
		return nil, err
	}
	ga, err := ag.apps.GenericApp(*in.AppMeta.AppId)
	if err != nil {
		return nil, err
	}
	asa, err := ga.ServiceAccount().Get(*in.AccountMeta.Username)
	if err != nil {
		return nil, err
	}
	if asa.Spec.RoleRef.ClusterRoleName == "" {
		return nil, fmt.Errorf("username:%s empty ClusterRoleName", *in.AccountMeta.Username)
	}
	if asa.Spec.AccountMeta.Password != *in.AccountMeta.Password {
		return nil, fmt.Errorf("accountmeta invalid password")
	}
	role, err := ga.Role().Get(asa.Spec.RoleRef.ClusterRoleName)
	if err != nil {
		return nil, err
	}
	data, err := role.Marshal()
	if err != nil {
		return nil, err
	}
	success := true
	return &pb.ClusterRoleReply{
		Result: &pb.StandardReply{
			Success: &success,
			Reason:  nil,
		},
		Role: data,
	}, nil
}

func (ag *Aggregator) ValidateRule(ctx context.Context, in *pb.ValidateRuleRequest) (*pb.ValidateRuleReply, error) {
	if err := ag.apps.Validate(*in.AppMeta.AppId, *in.AppMeta.AppSecret); err != nil {
		return nil, err
	}
	ga, err := ag.apps.GenericApp(*in.AppMeta.AppId)
	if err != nil {
		return nil, err
	}
	asa, err := ga.ServiceAccount().Get(*in.AccountMeta.Username)
	if err != nil {
		return nil, err
	}
	if asa.Spec.RoleRef.ClusterRoleName == "" {
		return nil, fmt.Errorf("username:%s empty ClusterRoleName", *in.AccountMeta.Username)
	}
	role, err := ga.Role().Get(asa.Spec.RoleRef.ClusterRoleName)
	if err != nil {
		return nil, err
	}
	match := false
	for _, v := range role.Spec.Rules {
		if v.Namespace == *in.Rule.Namespace {
			if v.GroupVersionKind.Group == *in.Rule.GroupVersionKind.Group &&
				v.GroupVersionKind.Version == *in.Rule.GroupVersionKind.Version &&
				v.GroupVersionKind.Kind == *in.Rule.GroupVersionKind.Kind {
				for _, verb := range v.Verbs {
					if verb == *in.Rule.Verb {
						match = true
					}
				}
			}
		}
	}
	success := true
	return &pb.ValidateRuleReply{
		Result: &pb.StandardReply{
			Success: &success,
			Reason:  nil,
		},
		Root: &match,
	}, nil
}
