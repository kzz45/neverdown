package resources

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kzz45/neverdown/pkg/jingx/registry"

	"github.com/kzz45/neverdown/pkg/jwttoken"
	"github.com/kzz45/neverdown/pkg/zaplogger"

	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"

	authority "github.com/kzz45/neverdown/pkg/authx/client-go"
	"github.com/kzz45/neverdown/pkg/authx/rbac/admin"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	"github.com/kzz45/neverdown/pkg/client-go/clientset/versioned/scheme"
	"github.com/kzz45/neverdown/pkg/jingx/aggregator"
	"github.com/kzz45/neverdown/pkg/jingx/proto"
	"github.com/kzz45/neverdown/pkg/kubernetes/providers"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

const (
	AppName = "jingx"
)

func New(ctx context.Context, authorityClientSet kubernetes.Interface, api aggregator.Api) *Resource {
	sub, cancel := context.WithCancel(ctx)
	r := &Resource{
		ClientSet:          api.ClientSet(),
		authorityClientSet: authorityClientSet,
		Codec:              providers.LegacyCodec(),
		api:                api,
		ctx:                sub,
		cancel:             cancel,
	}
	r.initAuthorityApp()
	r.registerKinds()
	r.registerRoles()
	r.registerServiceAccount()
	return r
}

type Resource struct {
	ClientSet          kubernetes.Interface
	authorityClientSet kubernetes.Interface

	Codec runtime.Codec

	Authority *authority.Client
	AppId     string

	api aggregator.Api

	ctx    context.Context
	cancel context.CancelFunc
}

func errNotNil(err error) bool {
	// if err != nil {
	// 	return true
	// }
	// return false
	return err != nil
}

func (r *Resource) initAuthorityApp() {
	defaultBackoff := wait.Backoff{
		Steps:    6,
		Duration: 200 * time.Millisecond,
		Factor:   5.0,
		Jitter:   0.1,
	}
	err := retry.OnError(defaultBackoff, errNotNil, func() error {
		ctx, cancel := context.WithTimeout(r.ctx, time.Second*10)
		app, err := admin.CreateRbacV1App(ctx, r.authorityClientSet, &rbacv1.App{
			ObjectMeta: metav1.ObjectMeta{
				Name: AppName,
			},
		})
		cancel()
		if err != nil {
			if !errors.IsAlreadyExists(err) {
				zaplogger.Sugar().Error(err)
				return err
			} else {
				ctx, cancel = context.WithTimeout(r.ctx, time.Second*5)
				app, err = admin.GetRbacV1App(ctx, r.authorityClientSet, AppName)
				cancel()
				if err != nil {
					zaplogger.Sugar().Error(err)
					return err
				}
			}
		}
		r.AppId = app.Spec.Id
		r.Authority = authority.New(r.ctx, &authority.Option{
			AppId:     app.Spec.Id,
			AppSecret: app.Spec.Secret,
			ClientSet: r.authorityClientSet,
		})
		return nil
	})
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
}

func (r *Resource) registerKinds() {
	<-time.After(time.Millisecond * 1200)
	gvks := r.AllSupportGVKs()
	for _, gvk := range gvks {
		rule := &rbacv1.GroupVersionKindRule{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "",
				Namespace: strings.ToLower(r.AppId),
			},
			Spec: rbacv1.RuleSpec{
				GroupVersionKind: gvk,
				Verbs:            r.AllSupportVerbs(),
			},
		}
		if err := r.Authority.GenericApp.Kind().Update(r.ctx, rule); err != nil {
			if errors.IsNotFound(err) {
				if err := r.Authority.GenericApp.Kind().Create(r.ctx, rule); err != nil {
					if !errors.IsAlreadyExists(err) {
						zaplogger.Sugar().Error(err)
					}
				}
			} else {
				zaplogger.Sugar().Error(err)
			}
		}
	}
}

func (r *Resource) AllSupportVerbs() []string {
	return []string{
		string(proto.VerbCreate),
		//string(proto.VerbGet),
		string(proto.VerbList),
		string(proto.VerbUpdate),
		string(proto.VerbDelete),
		string(proto.VerbWatch),
	}
}

func (r *Resource) AllSupportGVKs() []rbacv1.GroupVersionKind {
	filters := map[string]bool{
		"ListOptions":   true,
		"GetOptions":    true,
		"DeleteOptions": true,
		"CreateOptions": true,
		"UpdateOptions": true,
		"PatchOptions":  true,
		"ExportOptions": true,
	}
	gvks := make([]rbacv1.GroupVersionKind, 0)
	for gvk := range scheme.Scheme.AllKnownTypes() {
		if gvk.Group == "" {
			continue
		}
		if gvk.Group != jingxv1.GroupName {
			continue
		}
		if gvk.Version == runtime.APIVersionInternal {
			continue
		}
		if _, ok := filters[gvk.Kind]; ok {
			continue
		}
		if gvk.Kind == metav1.WatchEventKind {
			continue
		}
		reg := regexp.MustCompile(`List`)
		if isList := reg.Match([]byte(gvk.Kind)); isList {
			continue
		}
		gvks = append(gvks, rbacv1.GroupVersionKind{
			Group:   gvk.Group,
			Version: gvk.Version,
			Kind:    gvk.Kind,
		})
	}
	return gvks
}

func (r *Resource) getRolesWith(namespace string) []rbacv1.PolicyRule {
	rules := make([]rbacv1.PolicyRule, 0)
	for _, gvk := range r.AllSupportGVKs() {
		rules = append(rules, rbacv1.PolicyRule{
			Namespace:        namespace,
			GroupVersionKind: gvk,
			Verbs:            r.AllSupportVerbs(),
		})
	}
	return rules
}

func (r *Resource) registerRoles() {
	rules := make([]rbacv1.PolicyRule, 0)
	for _, ns := range []string{registry.DefaultNamespace} {
		rules = append(rules, r.getRolesWith(ns)...)
	}
	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "admin",
			Namespace: "",
		},
		Spec: rbacv1.ClusterRoleSpec{
			Desc:  "",
			Rules: rules,
		},
	}
	if err := r.Authority.GenericApp.Role().Create(r.ctx, role); err != nil {
		if !errors.IsAlreadyExists(err) {
			zaplogger.Sugar().Fatal(err)
			return
		}
		zaplogger.Sugar().Infof("Already exist Role:%s appid:%s", role.Name, r.AppId)
	} else {
		zaplogger.Sugar().Infof("Successful create Role:%s appid:%s", role.Name, r.AppId)
	}
}

func (r *Resource) registerServiceAccount() {
	sa := &rbacv1.AppServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "",
			Namespace: "",
		},
		Spec: rbacv1.AppServiceAccountSpec{
			Desc: "",
			RoleRef: rbacv1.RoleRef{
				ClusterRoleName: "admin",
			},
			AccountMeta: rbacv1.AccountMeta{
				Username: "admin",
				Password: "",
			},
		},
	}
	if err := r.Authority.GenericApp.ServiceAccount().Create(r.ctx, sa); err != nil {
		if !errors.IsAlreadyExists(err) {
			zaplogger.Sugar().Fatal(err)
		}
		zaplogger.Sugar().Infof("Already exist ServiceAccount:%s appid:%s", sa.Name, r.AppId)
	} else {
		zaplogger.Sugar().Infof("Successful create ServiceAccount:%s appid:%s", sa.Name, r.AppId)
	}
}

func (r *Resource) errorResponse(code int32, err error) ([]byte, error) {
	res := &proto.Response{
		Code:             code,
		GroupVersionKind: rbacv1.GroupVersionKind{},
		Namespace:        "",
		Verb:             "",
		Raw:              []byte(err.Error()),
	}
	return res.Marshal()
}

func (r *Resource) Login(in []byte) ([]byte, error) {
	req := &proto.Request{}
	if err := req.Unmarshal(in); err != nil {
		return r.errorResponse(http.StatusBadRequest, err)
	}
	meta := &rbacv1.AccountMeta{}
	// zaplogger.Sugar().Infof("Login request:%s", string(req.Raw))
	// zaplogger.Sugar().Infof("Login request:%v", meta)
	if err := meta.Unmarshal(req.Raw); err != nil {
		return r.errorResponse(http.StatusBadRequest, err)
	}
	if len(meta.Username) == 0 || len(meta.Password) == 0 {
		return r.errorResponse(http.StatusBadRequest, fmt.Errorf("nil username or password"))
	}
	token, err := r.Authority.ValidateAccount(meta.Username, meta.Password)
	if err != nil {
		return r.errorResponse(http.StatusForbidden, err)
	}
	role, err := r.Authority.ClusterRole(meta.Username)
	if err != nil {
		return r.errorResponse(http.StatusForbidden, err)
	}
	ctx := &proto.Context{
		Token:       token,
		IsAdmin:     false,
		ExpireAt:    int32(time.Now().Unix() + jwttoken.GetTokenExpirationFromEnv()),
		ClusterRole: *role,
	}
	raw, err := ctx.Marshal()
	if err != nil {
		return r.errorResponse(http.StatusInternalServerError, err)
	}
	res := &proto.Response{
		Code:             0,
		GroupVersionKind: rbacv1.GroupVersionKind{},
		Namespace:        "",
		Verb:             "",
		Raw:              raw,
	}
	return res.Marshal()
}

func (r *Resource) ValidateAccess(username string, namespace string, gvk rbacv1.GroupVersionKind, verb string) (bool, error) {
	return r.Authority.ValidateRules(username, authority.Rule{
		Namespace:        namespace,
		GroupVersionKind: gvk,
		Verb:             verb,
	})
}

//VerbCreate Verb = "create"
//VerbDelete Verb = "delete"
//VerbGet    Verb = "get"
//VerbList   Verb = "list"
//VerbPatch  Verb = "patch"
//VerbUpdate Verb = "update"
//VerbWatch  Verb = "watch"

func (r *Resource) Handler(username string, req *proto.Request) (int32, []byte, error) {
	bo, err := r.ValidateAccess(username, req.Namespace, req.GroupVersionKind, string(req.Verb))
	if err != nil {
		return http.StatusBadRequest, []byte(""), err
	}
	if !bo {
		return http.StatusUnauthorized, []byte(""), fmt.Errorf("401 StatusUnauthorized")
	}
	gvk := &schema.GroupVersionKind{
		Group:   req.GroupVersionKind.Group,
		Version: req.GroupVersionKind.Version,
		Kind:    req.GroupVersionKind.Kind,
	}
	switch req.Verb {
	case proto.VerbCreate:
		return r.api.Create(gvk, req.Namespace, req.Raw, username)
	case proto.VerbDelete:
		return r.api.Delete(gvk, req.Namespace, req.Raw, username)
	case proto.VerbGet:
	case proto.VerbList:
		return r.api.List(gvk, req.Namespace, req.Raw)
	case proto.VerbUpdate:
		return r.api.Update(gvk, req.Namespace, req.Raw, username)
	case proto.VerbWatch:
	}
	return 1, []byte("failed"), fmt.Errorf("no verb:%s exist", req.Verb)
}
