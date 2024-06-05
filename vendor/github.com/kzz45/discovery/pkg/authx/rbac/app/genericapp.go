package app

import (
	"context"
	"strings"

	rbacv1 "github.com/kzz45/discovery/pkg/apis/rbac/v1"
	"github.com/kzz45/discovery/pkg/client-go/kubernetes"
)

type GenericApp struct {
	ctx            context.Context
	app            *rbacv1.App
	kind           Kind
	role           Role
	serviceAccount ServiceAccount
}

func NewGenericApp(ctx context.Context, cfg kubernetes.Interface, a *rbacv1.App) *GenericApp {
	kind := NewGroupVersionKinds(ctx, cfg, strings.ToLower(a.Spec.Id))
	roles := NewRoles(ctx, cfg, kind)
	asa := NewAppServiceAccounts(ctx, cfg, roles)
	ga := &GenericApp{
		ctx:            ctx,
		app:            a,
		kind:           kind,
		role:           roles,
		serviceAccount: asa,
	}
	return ga
}

func (ga *GenericApp) Sync(a *rbacv1.App) {
	if ga.app.UID == a.UID {
		if ga.app.ResourceVersion >= a.ResourceVersion {
			return
		}
	}
	ga.app = a
}

func (ga *GenericApp) RbacV1App() *rbacv1.App {
	return ga.app
}

func (ga *GenericApp) Kind() Kind {
	return ga.kind
}

func (ga *GenericApp) Role() Role {
	return ga.role
}

func (ga *GenericApp) ServiceAccount() ServiceAccount {
	return ga.serviceAccount
}
