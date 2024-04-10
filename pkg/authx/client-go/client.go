package authority

import (
	"context"
	"fmt"

	"github.com/kzz45/neverdown/pkg/jwttoken"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	"github.com/kzz45/neverdown/pkg/authx/rbac/admin"
	"github.com/kzz45/neverdown/pkg/authx/rbac/app"
	kubernetes "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/klog/v2"
)

type Option struct {
	AppId     string
	AppSecret string
	ClientSet kubernetes.Interface
}

// GetLabelSelector returns the LabelSelector of the metav1.ListOptions
func GetLabelSelector(in map[string]string) string {
	ls := labels.NewSelector()
	for k, v := range in {
		req, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			klog.Fatal(err)
		}
		ls = ls.Add(*req)
	}
	return ls.String()
}

func (opt *Option) Validate(ctx context.Context) error {
	listOption := metav1.ListOptions{
		LabelSelector: GetLabelSelector(map[string]string{admin.LabelKey: opt.AppId}),
	}
	list, err := opt.ClientSet.RbacV1().Apps(admin.DefaultNamespace).List(ctx, listOption)
	if err != nil {
		return err
	}
	if len(list.Items) != 1 {
		return fmt.Errorf("error app:%s was not available, cause there were %d accounts", opt.AppId, len(list.Items))
	}
	a := list.Items[0]
	if a.Spec.Id == opt.AppId && a.Spec.Secret == opt.AppSecret {
		return nil
	}
	return fmt.Errorf("error invalid appid or appsecret")
}

type Client struct {
	opt        *Option
	GenericApp *app.GenericApp
}

func New(ctx context.Context, opt *Option) *Client {
	if err := opt.Validate(ctx); err != nil {
		klog.Fatal(err)
	}
	c := &Client{
		opt: opt,
		GenericApp: app.NewGenericApp(ctx, opt.ClientSet, &rbacv1.App{
			Spec: rbacv1.AppSpec{
				Id: opt.AppId,
			},
		}),
	}
	return c
}

func (c *Client) ValidateAccount(username string, password string) (string, error) {
	bo := c.GenericApp.ServiceAccount().Validate(username, password)
	if !bo {
		return "", fmt.Errorf("invalid username or password")
	}
	return jwttoken.Generate(username, false)
}

type Rule struct {
	Namespace        string
	GroupVersionKind rbacv1.GroupVersionKind
	Verb             string
}

func (c *Client) ValidateRules(username string, rule Rule) (bool, error) {
	asa, err := c.GenericApp.ServiceAccount().Get(username)
	if err != nil {
		return false, err
	}
	if asa.Spec.RoleRef.ClusterRoleName == "" {
		return false, nil
	}
	role, err := c.GenericApp.Role().Get(asa.Spec.RoleRef.ClusterRoleName)
	if err != nil {
		return false, err
	}
	match := false
	for _, v := range role.Spec.Rules {
		if v.Namespace == rule.Namespace {
			if v.GroupVersionKind.String() == rule.GroupVersionKind.String() {
				for _, verb := range v.Verbs {
					if verb == rule.Verb {
						match = true
					}
				}
			}
		}
	}
	if !match {
		return false, nil
	}
	return true, nil
}

// todo ClusterRole should validate username and password at the sametime here. Maybe this is a Bug
func (c *Client) ClusterRole(username string) (*rbacv1.ClusterRole, error) {
	asa, err := c.GenericApp.ServiceAccount().Get(username)
	if err != nil {
		return nil, err
	}
	if asa.Spec.RoleRef.ClusterRoleName == "" {
		return &rbacv1.ClusterRole{}, nil
	}
	return c.GenericApp.Role().Get(asa.Spec.RoleRef.ClusterRoleName)
}
