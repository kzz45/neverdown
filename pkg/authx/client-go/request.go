package authority

import (
	"context"
	"fmt"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	"k8s.io/klog/v2"
)

var client *Client

func Init(ctx context.Context, opt *Option) {
	if client != nil {
		return
	}
	client = New(ctx, opt)
}

type Request struct {
	c *Client

	username         string
	password         string
	namespace        string
	groupVersionKind rbacv1.GroupVersionKind
	verb             string
}

func NewRequest() *Request {
	if client == nil {
		klog.Fatal("please Init before call Request")
	}
	r := &Request{
		c: client,
	}
	return r
}

func (r *Request) Username(username string) *Request {
	r.username = username
	return r
}

func (r *Request) Password(password string) *Request {
	r.password = password
	return r
}

func (r *Request) Namespace(namespace string) *Request {
	r.namespace = namespace
	return r
}

func (r *Request) GroupVersionKind(gvk rbacv1.GroupVersionKind) *Request {
	r.groupVersionKind = gvk
	return r
}

func (r *Request) Verb(verb string) *Request {
	r.verb = verb
	return r
}

func (r *Request) ValidateAccount() (string, error) {
	if r.username == "" || r.password == "" {
		return "", fmt.Errorf("ValidateAccount must specify username and password first")
	}
	return r.c.ValidateAccount(r.username, r.password)
}

func (r *Request) ValidateAccess() (bool, error) {
	if r.username == "" {
		return false, fmt.Errorf("ValidateAccess must specify username")
	}
	if r.groupVersionKind.Kind == "" {
		return false, fmt.Errorf("ValidateAccess must specify GroupVersionKind kind")
	}
	if r.verb == "" {
		return false, fmt.Errorf("ValidateAccess must specify GroupVersionKind verb")
	}
	return r.c.ValidateRules(r.username, Rule{
		Namespace:        r.namespace,
		GroupVersionKind: r.groupVersionKind,
		Verb:             r.verb,
	})
}

func (r *Request) ClusterRole() (*rbacv1.ClusterRole, error) {
	if r.username == "" {
		return nil, fmt.Errorf("ClusterRole must specify username")
	}
	return r.c.ClusterRole(r.username)
}
