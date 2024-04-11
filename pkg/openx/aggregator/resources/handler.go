package resources

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kzz45/neverdown/pkg/jwttoken"
	"github.com/kzz45/neverdown/pkg/zaplogger"
	"go.uber.org/zap"

	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	authority "github.com/kzz45/neverdown/pkg/authx/client-go"
	"github.com/kzz45/neverdown/pkg/jingx/registry"

	"github.com/kzz45/neverdown/pkg/openx/aggregator/proto"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func (r *Resources) errorResponse(code int32, err error) ([]byte, error) {
	res := &proto.Response{
		Code:             code,
		GroupVersionKind: rbacv1.GroupVersionKind{},
		Namespace:        "",
		Verb:             "",
		Raw:              []byte(err.Error()),
	}
	return res.Marshal()
}

func (r *Resources) Login(in []byte) ([]byte, error) {
	req := &proto.Request{}
	if err := req.Unmarshal(in); err != nil {
		return r.errorResponse(http.StatusBadRequest, err)
	}
	meta := &rbacv1.AccountMeta{}
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

func (r *Resources) ValidateAccess(username string, namespace string, gvk rbacv1.GroupVersionKind, verb string) (bool, error) {
	return r.Authority.ValidateRules(username, authority.Rule{
		Namespace:        namespace,
		GroupVersionKind: gvk,
		Verb:             verb,
	})
}

func (r *Resources) Handler(username string, req *proto.Request) (int32, []byte, error) {
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
	r.metricsCollector.AccumulateRequest(req.Namespace, gvk.String(), string(req.Verb))
	switch gvk.GroupVersion().String() {
	case schema.GroupVersion{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version}.String():
		return r.jingxV1Handler(gvk, req.Namespace, req.Verb, req.Raw)
	case schema.GroupVersion{Group: metricsv1beta1.GroupName, Version: metricsv1beta1.SchemeGroupVersion.Version}.String():
		return r.metricsV1Beta1Handler(gvk, req.Namespace, req.Verb, req.Raw)
	default:
		return r.k8s1Handler(gvk, req.Namespace, req.Verb, req.Raw)
	}
}

func (r *Resources) k8s1Handler(gvk *schema.GroupVersionKind, namespace string, verb proto.Verb, raw []byte) (code int32, res []byte, err error) {
	switch verb {
	case proto.VerbCreate:
		return r.create(gvk, namespace, raw)
	case proto.VerbDelete:
		return r.delete(gvk, namespace, raw)
	case proto.VerbGet:
	case proto.VerbList:
		return r.list(gvk, namespace, raw)
	case proto.VerbUpdate:
		return r.update(gvk, namespace, raw)
	case proto.VerbWatch:
	}
	return 1, []byte("failed"), fmt.Errorf("no verb:%s exist", verb)
}

func (r *Resources) jingxV1Handler(gvk *schema.GroupVersionKind, namespace string, verb proto.Verb, raw []byte) (code int32, res []byte, err error) {
	switch verb {
	case proto.VerbList:
		return r.jingxV1List(gvk, namespace, raw)
	}
	return 1, []byte("failed"), fmt.Errorf("no verb:%s exist", verb)
}

func (r *Resources) jingxV1List(gvk *schema.GroupVersionKind, namespace string, raw []byte) (code int32, res []byte, err error) {
	ctx, cancel := context.WithCancel(r.ctx)
	defer cancel()
	var timeout int64 = 5
	listOpts := metav1.ListOptions{
		TimeoutSeconds: &timeout,
	}
	var objList Object
	switch *gvk {
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}:
		objList, err = r.discoveryClientSet.JingxV1().Projects(registry.DefaultNamespace).List(ctx, listOpts)
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}:
		objList, err = r.discoveryClientSet.JingxV1().Repositories(registry.DefaultNamespace).List(ctx, listOpts)
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}:
		objList, err = r.discoveryClientSet.JingxV1().Tags(registry.DefaultNamespace).List(ctx, listOpts)
	case schema.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Event"}:
		objList, err = r.discoveryClientSet.JingxV1().Events(registry.DefaultNamespace).List(ctx, listOpts)
	default:
		err = fmt.Errorf(ErrGVKNotExist, *gvk)
	}
	if err != nil {
		return 1, nil, err
	}
	res, err = objList.Marshal()
	if err != nil {
		return 1, nil, err
	}
	return code, res, err
}

func (r *Resources) metricsV1Beta1Handler(gvk *schema.GroupVersionKind, namespace string, verb proto.Verb, raw []byte) (code int32, res []byte, err error) {
	switch verb {
	case proto.VerbList:
		return r.metricsV1Beta1List(gvk, namespace, raw)
	}
	return 1, []byte("failed"), fmt.Errorf("no verb:%s exist", verb)
}

func (r *Resources) metricsV1Beta1List(gvk *schema.GroupVersionKind, namespace string, raw []byte) (code int32, res []byte, err error) {
	ctx, cancel := context.WithCancel(r.ctx)
	defer cancel()
	var timeout int64 = 5
	listOpts := metav1.ListOptions{
		TimeoutSeconds: &timeout,
	}
	var objList Object
	switch *gvk {
	case schema.GroupVersionKind{Group: metricsv1beta1.GroupName, Version: metricsv1beta1.SchemeGroupVersion.Version, Kind: "NodeMetrics"}:
		t1 := time.Now()
		objList, err = r.metricsClientSet.MetricsV1beta1().NodeMetricses().List(ctx, listOpts)
		zaplogger.Sugar().Infow("list NodeMetricses",
			zap.Int("length", len(objList.(*metricsv1beta1.NodeMetricsList).Items)),
			zap.Any("used-ms", time.Now().Sub(t1).Milliseconds()),
		)
	case schema.GroupVersionKind{Group: metricsv1beta1.GroupName, Version: metricsv1beta1.SchemeGroupVersion.Version, Kind: "PodMetrics"}:
		t1 := time.Now()
		objList, err = r.metricsClientSet.MetricsV1beta1().PodMetricses(namespace).List(ctx, listOpts)
		zaplogger.Sugar().Debugw("list PodMetricses",
			zap.String("namespace", namespace),
			zap.Int("length", len(objList.(*metricsv1beta1.PodMetricsList).Items)),
			zap.Any("used-ms", time.Now().Sub(t1).Milliseconds()),
		)
	default:
		err = fmt.Errorf(ErrGVKNotExist, *gvk)
	}
	if err != nil {
		return 1, nil, err
	}
	res, err = objList.Marshal()
	if err != nil {
		return 1, nil, err
	}
	return code, res, err
}
