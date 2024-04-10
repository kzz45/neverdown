package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	"github.com/kzz45/neverdown/pkg/authx/http/proto"
	"github.com/kzz45/neverdown/pkg/authx/rbac/admin"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	"github.com/kzz45/neverdown/pkg/jwttoken"

	"k8s.io/client-go/rest"
)

const (
	StaticCertFile = "TLS_OPTION_CERT_FILE"
	StaticKeyFile  = "TLS_OPTION_KEY_FILE"
)

type Handler struct {
	ctx                     context.Context
	cfg                     kubernetes.Interface
	adminApp                admin.App
	adminRbacServiceAccount admin.RbacServiceAccount
}

func NewResource(ctx context.Context) kubernetes.Interface {
	// out_of_cluster
	// certFile, keyFile := os.Getenv(StaticCertFile), os.Getenv(StaticKeyFile)
	// kubeconfig := &rest.Config{
	// 	Host: "127.0.0.1:9443",
	// 	TLSClientConfig: rest.TLSClientConfig{
	// 		Insecure: true,
	// 		CertFile: certFile,
	// 		KeyFile:  keyFile,
	// 	},
	// }
	// in_cluster
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	cfg, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return cfg
}

func NewHandler(ctx context.Context) *Handler {
	cfg := NewResource(ctx)
	adminApp := admin.NewApps(ctx, cfg, admin.DefaultNamespace)
	adminRbacServiceAccount := admin.NewServiceAccount(ctx, cfg, adminApp)
	h := &Handler{
		ctx:                     ctx,
		cfg:                     cfg,
		adminApp:                adminApp,
		adminRbacServiceAccount: adminRbacServiceAccount,
	}
	return h
}

func (h *Handler) AdminApp() admin.App {
	return h.adminApp
}

func (h *Handler) Handle(token string, in []byte) (res []byte, err error) {
	var errCode int32 = 1
	req := &proto.Request{}
	if err = req.Unmarshal(in); err != nil {
		zaplogger.Sugar().Error(err)
		goto end
	}
	switch req.ServiceRoute {
	case proto.Login:
		res, err = h.login(req.Data)
	default:
		var claims *jwttoken.Claims
		claims, err = h.validate(token)
		if err != nil {
			zaplogger.Sugar().Error(err)
			errCode = http.StatusUnauthorized
			goto end
		}
		res, err = h.switchRoute(claims, req)
	}
end:
	response := &proto.Response{
		ServiceRoute: req.ServiceRoute,
		Code:         0,
		Message:      "",
		Data:         nil,
	}
	if err != nil {
		response.Code = errCode
		response.Message = err.Error()
	} else {
		response.Data = res
	}
	return response.Marshal()
}

const (
	MethodList   = "list"
	MethodCreate = "create"
	MethodGet    = "get"
	MethodUpdate = "update"
	MethodDelete = "delete"
)

const (
	NoRoot = "is not admin"
)

func (h *Handler) switchRoute(claims *jwttoken.Claims, req *proto.Request) (res []byte, err error) {
	args := strings.Split(string(req.ServiceRoute), "/")
	if len(args) != 3 {
		return nil, fmt.Errorf("invalid ServiceRoute:%s", req.ServiceRoute)
	}
	switch args[1] {
	case "apps":
		res, err = h.adminAppHandler(claims, args[2], req.Data)
	case "rbacaccount":
		if !claims.IsAdmin {
			return nil, fmt.Errorf(NoRoot)
		}
		res, err = h.adminRbacServiceAccountHandler(args[2], req.Data)
	case "appaccount":
		res, err = h.appServiceAccountHandler(claims, args[2], req.Data)
	case "clusterrole":
		res, err = h.appRoleHandler(claims, args[2], req.Data)
	case "kind":
		res, err = h.appKindHandler(claims, args[2], req.Data)
	default:
		err = fmt.Errorf("invalid ServiceRoute:%s Object:%s", req.ServiceRoute, args[1])
	}
	return res, err
}

func (h *Handler) accountApps(username string) (map[string]bool, error) {
	rsa, err := h.adminRbacServiceAccount.Get(username)
	if err != nil {
		return nil, err
	}
	t := make(map[string]bool, 0)
	for _, v := range rsa.Spec.Apps {
		t[v] = true
	}
	return t, nil
}

func (h *Handler) validateAccountAccess(claims *jwttoken.Claims, appId string) error {
	if claims.IsAdmin {
		return nil
	}
	rsa, err := h.adminRbacServiceAccount.Get(claims.Username)
	if err != nil {
		return err
	}
	ga, err := h.adminApp.GenericApp(appId)
	if err != nil {
		return err
	}
	for _, v := range rsa.Spec.Apps {
		if ga.RbacV1App().Name == v {
			return nil
		}
	}
	return fmt.Errorf(NoRoot)
}

func (h *Handler) adminAppHandler(claims *jwttoken.Claims, method string, data []byte) (res []byte, err error) {
	if method != MethodList && !claims.IsAdmin {
		return nil, fmt.Errorf(NoRoot)
	}
	switch method {
	case MethodList:
		l, err := h.adminApp.List()
		if err != nil {
			return nil, err
		}
		if !claims.IsAdmin {
			items, err := h.accountApps(claims.Username)
			if err != nil {
				return nil, err
			}
			t := make([]rbacv1.App, 0)
			for _, v := range l.Items {
				if _, ok := items[v.Name]; ok {
					t = append(t, v)
				}
			}
			l.Items = t
		}
		return l.Marshal()
	case MethodCreate:
		app := &rbacv1.App{}
		if err = app.Unmarshal(data); err != nil {
			return nil, err
		}
		if err = h.adminApp.Create(h.ctx, app); err != nil {
			return nil, err
		}
	case MethodUpdate:
		app := &rbacv1.App{}
		if err = app.Unmarshal(data); err != nil {
			return nil, err
		}
		if err = h.adminApp.Update(h.ctx, app); err != nil {
			return nil, err
		}
	case MethodDelete:
		app := &rbacv1.App{}
		if err = app.Unmarshal(data); err != nil {
			return nil, err
		}
		if err = h.adminApp.Delete(h.ctx, app.Name); err != nil {
			return nil, err
		}
	default:
		err = fmt.Errorf("adminAppHandler invalid method:%s", method)
	}
	return res, err
}

func (h *Handler) adminRbacServiceAccountHandler(method string, data []byte) (res []byte, err error) {
	switch method {
	case MethodList:
		l, err := h.adminRbacServiceAccount.List()
		if err != nil {
			return nil, err
		}
		return l.Marshal()
	case MethodCreate:
		rsa := &rbacv1.RbacServiceAccount{}
		if err = rsa.Unmarshal(data); err != nil {
			return nil, err
		}
		if err = h.adminRbacServiceAccount.Create(h.ctx, rsa); err != nil {
			return nil, err
		}
	case MethodUpdate:
		rsa := &rbacv1.RbacServiceAccount{}
		if err = rsa.Unmarshal(data); err != nil {
			return nil, err
		}
		if err = h.adminRbacServiceAccount.Update(h.ctx, rsa); err != nil {
			return nil, err
		}
	case MethodDelete:
		rsa := &rbacv1.RbacServiceAccount{}
		if err = rsa.Unmarshal(data); err != nil {
			return nil, err
		}
		if err = h.adminRbacServiceAccount.Delete(h.ctx, rsa.Name); err != nil {
			return nil, err
		}
	default:
		err = fmt.Errorf("adminRbacServiceAccountHandler invalid method:%s", method)
	}
	return res, err
}

func (h *Handler) appKindHandler(claims *jwttoken.Claims, method string, data []byte) (res []byte, err error) {
	kind := &rbacv1.GroupVersionKindRule{}
	if err = kind.Unmarshal(data); err != nil {
		return nil, err
	}
	appId := kind.Namespace
	if err = h.validateAccountAccess(claims, appId); err != nil {
		return nil, err
	}
	genericApp, err := h.adminApp.GenericApp(strings.ToLower(appId))
	if err != nil {
		return nil, err
	}
	switch method {
	case MethodList:
		l, err := genericApp.Kind().List()
		if err != nil {
			return nil, err
		}
		return l.Marshal()
	case MethodCreate:
		if err = genericApp.Kind().Create(h.ctx, kind); err != nil {
			return nil, err
		}
	case MethodUpdate:
		if err = genericApp.Kind().Update(h.ctx, kind); err != nil {
			return nil, err
		}
	case MethodDelete:
		if err = genericApp.Kind().Delete(h.ctx, kind.Name); err != nil {
			return nil, err
		}
	default:
		err = fmt.Errorf("appKindHandler invalid method:%s", method)
	}
	return res, err
}

func (h *Handler) appRoleHandler(claims *jwttoken.Claims, method string, data []byte) (res []byte, err error) {
	role := &rbacv1.ClusterRole{}
	if err = role.Unmarshal(data); err != nil {
		return nil, err
	}
	appId := role.Namespace
	if err = h.validateAccountAccess(claims, appId); err != nil {
		return nil, err
	}
	genericApp, err := h.adminApp.GenericApp(strings.ToLower(appId))
	if err != nil {
		return nil, err
	}
	switch method {
	case MethodList:
		l, err := genericApp.Role().List()
		if err != nil {
			return nil, err
		}
		return l.Marshal()
	case MethodCreate:
		if err = genericApp.Role().Create(h.ctx, role); err != nil {
			return nil, err
		}
	case MethodUpdate:
		if err = genericApp.Role().Update(h.ctx, role); err != nil {
			return nil, err
		}
	case MethodDelete:
		if err = genericApp.Role().Delete(h.ctx, role.Name); err != nil {
			return nil, err
		}
	default:
		err = fmt.Errorf("appRoleHandler invalid method:%s", method)
	}
	return res, err
}

func (h *Handler) appServiceAccountHandler(claims *jwttoken.Claims, method string, data []byte) (res []byte, err error) {
	asa := &rbacv1.AppServiceAccount{}
	if err = asa.Unmarshal(data); err != nil {
		return nil, err
	}
	appId := asa.Namespace
	if err = h.validateAccountAccess(claims, appId); err != nil {
		return nil, err
	}
	genericApp, err := h.adminApp.GenericApp(strings.ToLower(appId))
	if err != nil {
		return nil, err
	}
	switch method {
	case MethodList:
		l, err := genericApp.ServiceAccount().List()
		if err != nil {
			return nil, err
		}
		return l.Marshal()
	case MethodCreate:
		if err = genericApp.ServiceAccount().Create(h.ctx, asa); err != nil {
			return nil, err
		}
	case MethodUpdate:
		if err = genericApp.ServiceAccount().Update(h.ctx, asa); err != nil {
			return nil, err
		}
	case MethodDelete:
		if err = genericApp.ServiceAccount().Delete(h.ctx, asa.Name); err != nil {
			return nil, err
		}
	default:
		err = fmt.Errorf("appServiceAccountHandler invalid method:%s", method)
	}
	return res, err
}
