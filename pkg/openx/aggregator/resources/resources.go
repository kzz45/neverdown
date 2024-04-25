package resources

import (
	"context"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	openxv1scheme "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned/scheme"
	openxinformers "github.com/kzz45/neverdown/pkg/client-go/informers/externalversions"
	"github.com/kzz45/neverdown/pkg/openx/aggregator/proto"
	"github.com/kzz45/neverdown/pkg/openx/clientbuilder"
	"github.com/kzz45/neverdown/pkg/openx/metrics"

	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	authority "github.com/kzz45/neverdown/pkg/authx/client-go"
	"github.com/kzz45/neverdown/pkg/authx/rbac/admin"
	discoveryclientset "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"

	// networkingv1 "k8s.io/api/networking/v1"
	nativerbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	nativescheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/util/retry"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

const (
	AppName = "openx-apiserver"
)

var PodsGroupVersionKind = schema.GroupVersionKind{Group: "native", Version: "v1", Kind: "Pod"}

type Resources struct {
	ctx context.Context

	// nativa k8s interface
	nativeClientSet  kubernetes.Interface
	openxClientSet   discoveryclientset.Interface
	metricsClientSet metricsclientset.Interface
	// authx interface
	authorityClientSet discoveryclientset.Interface
	// discovery interface
	discoveryClientSet discoveryclientset.Interface
	// ClientBuilder will provide a client for this controller to use
	clientBuilder clientbuilder.ControllerClientBuilder

	Authority *authority.Client
	AppId     string

	// InformerFactory gives access to informers for the controller.
	InformerFactory informers.SharedInformerFactory
	// OpenxInformerFactory gives access to openx informers for the controller.
	OpenXInformerFactory openxinformers.SharedInformerFactory

	watchers chan *proto.Response

	metricsCollector metrics.Collector
}

func New(ctx context.Context, clientBuilder clientbuilder.ControllerClientBuilder, informerFactory informers.SharedInformerFactory, openxInformerFactory openxinformers.SharedInformerFactory) *Resources {
	r := &Resources{
		ctx:                  ctx,
		nativeClientSet:      clientBuilder.ClientOrDie("openx-apiserver"),
		openxClientSet:       clientBuilder.OpenxClientOrDie("openx-apiserver"),
		metricsClientSet:     clientBuilder.MetricsClientOrDie("openx-apiserver"),
		authorityClientSet:   clientBuilder.AuthxClientOrDie("openx-apiserver"),
		discoveryClientSet:   clientBuilder.DiscoveryClientOrDie("openx-apiserver"),
		clientBuilder:        clientBuilder,
		InformerFactory:      informerFactory,
		OpenXInformerFactory: openxInformerFactory,
		metricsCollector:     metrics.NewPrometheusCollector(),
	}
	r.initAuthorityApp()
	r.registerKinds()
	r.registerRoles()
	r.registerServiceAccount()
	r.registerInformer()
	return r
}

func (r *Resources) ClientBuilder() clientbuilder.ControllerClientBuilder {
	return r.clientBuilder
}

func (r *Resources) MetricsCollector() metrics.Collector {
	return r.metricsCollector
}

func (r *Resources) initAuthorityApp() {
	defaultBackoff := wait.Backoff{
		Steps:    20,
		Duration: 200 * time.Millisecond,
		Factor:   1.0,
		Jitter:   0.1,
	}
	err := retry.OnError(defaultBackoff, errNotNil, func() error {
		ctx, cancel := context.WithTimeout(r.ctx, time.Second*5)
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

func nativeKindFilters(gvk schema.GroupVersionKind) bool {
	filters := map[schema.GroupVersion][]string{
		corev1.SchemeGroupVersion: {
			"Pod",
			"Service",
			"Endpoints",
			"Node",
			"Event",
			"Secret",
			"Namespace",
			"ConfigMap",
			"ServiceAccount",
			"PersistentVolume",
			"PersistentVolumeClaim",
		},
		appsv1.SchemeGroupVersion: {
			"Deployment",
			"StatefulSet",
			"DaemonSet",
		},
		// networkingv1.SchemeGroupVersion: {
		// 	"Ingress",
		// 	"IngressClass",
		// 	"NetworkPolicy",
		// },
		autoscalingv2.SchemeGroupVersion: {
			"HorizontalPodAutoscaler",
		},
		nativerbacv1.SchemeGroupVersion: {
			"ClusterRole",
			"ClusterRoleBinding",
		},
		metricsv1beta1.SchemeGroupVersion: {
			"NodeMetrics",
			"PodMetrics",
		},
	}
	t, ok := filters[gvk.GroupVersion()]
	if !ok {
		return false
	}
	for _, v := range t {
		if gvk.Kind == v {
			return true
		}
	}
	return false
}

func podsNativeKinds() (rbacv1.GroupVersionKind, []string) {
	return rbacv1.GroupVersionKind{
			Group:   PodsGroupVersionKind.Group,
			Version: PodsGroupVersionKind.Version,
			Kind:    PodsGroupVersionKind.Kind,
		}, []string{
			string(proto.VerbPodsSSH),
			string(proto.VerbPodsLogDownload),
			string(proto.VerbPodsLogStreaming),
		}
}

func metricsKinds() ([]rbacv1.GroupVersionKind, []string) {
	return []rbacv1.GroupVersionKind{
			{
				Group:   metricsv1beta1.GroupName,
				Version: metricsv1beta1.SchemeGroupVersion.Version,
				Kind:    "NodeMetrics",
			},
			{
				Group:   metricsv1beta1.GroupName,
				Version: metricsv1beta1.SchemeGroupVersion.Version,
				Kind:    "PodMetrics",
			},
		}, []string{
			string(proto.VerbGet),
			string(proto.VerbList),
			string(proto.VerbWatch),
		}
}

func GetRegisterKinds() []rbacv1.GroupVersionKind {
	res := make([]rbacv1.GroupVersionKind, 0)
	filters := map[string]bool{
		"ListOptions":   true,
		"GetOptions":    true,
		"DeleteOptions": true,
		"CreateOptions": true,
		"UpdateOptions": true,
		"PatchOptions":  true,
		"ExportOptions": true,
	}
	// k8s native gvk
	for gvk, _ := range nativescheme.Scheme.AllKnownTypes() {
		if ok := nativeKindFilters(gvk); !ok {
			continue
		}
		res = append(res, rbacv1.GroupVersionKind{
			Group:   gvk.Group,
			Version: gvk.Version,
			Kind:    gvk.Kind,
		})
	}
	// openxv1
	for gvk, _ := range openxv1scheme.Scheme.AllKnownTypes() {
		if gvk.Version == runtime.APIVersionInternal {
			continue
		}
		if _, ok := filters[gvk.Kind]; ok {
			continue
		}
		if reflect.DeepEqual(gvk.GroupVersion(), metav1.Unversioned) {
			continue
		}
		if gvk.Kind == metav1.WatchEventKind {
			continue
		}
		reg := regexp.MustCompile(`(.*?)List`)
		if isList := reg.Match([]byte(gvk.Kind)); isList {
			continue
		}
		res = append(res, rbacv1.GroupVersionKind{
			Group:   gvk.Group,
			Version: gvk.Version,
			Kind:    gvk.Kind,
		})
	}
	// discovery jingx
	// for gvk, _ := range discoveryscheme.Scheme.AllKnownTypes() {
	// 	if !reflect.DeepEqual(gvk.GroupVersion(), jingxv1.SchemeGroupVersion) {
	// 		continue
	// 	}
	// 	if _, ok := filters[gvk.Kind]; ok {
	// 		continue
	// 	}
	// 	if gvk.Kind == metav1.WatchEventKind {
	// 		continue
	// 	}
	// 	reg := regexp.MustCompile(`(.*?)List`)
	// 	if isList := reg.Match([]byte(gvk.Kind)); isList {
	// 		continue
	// 	}
	// 	res = append(res, rbacv1.GroupVersionKind{
	// 		Group:   gvk.Group,
	// 		Version: gvk.Version,
	// 		Kind:    gvk.Kind,
	// 	})
	// }
	return res
}

func allSupportVerbs() []string {
	return []string{
		string(proto.VerbCreate),
		string(proto.VerbGet),
		string(proto.VerbList),
		string(proto.VerbUpdate),
		string(proto.VerbDelete),
		string(proto.VerbWatch),
	}
}

func (r *Resources) registerKinds() {
	<-time.After(time.Millisecond * 1200)
	rules := make([]*rbacv1.GroupVersionKindRule, 0)
	for _, gvk := range GetRegisterKinds() {
		rule := &rbacv1.GroupVersionKindRule{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "",
				Namespace: strings.ToLower(r.AppId),
			},
			Spec: rbacv1.RuleSpec{
				GroupVersionKind: gvk,
				Verbs:            allSupportVerbs(),
			},
		}
		rules = append(rules, rule)
	}
	// pods/logs
	gvk, verbs := podsNativeKinds()
	rule := &rbacv1.GroupVersionKindRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "",
			Namespace: strings.ToLower(r.AppId),
		},
		Spec: rbacv1.RuleSpec{
			GroupVersionKind: gvk,
			Verbs:            verbs,
		},
	}
	rules = append(rules, rule)
	// metrics
	gvks, verbs := metricsKinds()
	for _, gvk := range gvks {
		rule = &rbacv1.GroupVersionKindRule{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "",
				Namespace: strings.ToLower(r.AppId),
			},
			Spec: rbacv1.RuleSpec{
				GroupVersionKind: gvk,
				Verbs:            verbs,
			},
		}
		rules = append(rules, rule)
	}
	for _, rule := range rules {
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

func (r *Resources) getRolesWith(namespace string) []rbacv1.PolicyRule {
	rules := make([]rbacv1.PolicyRule, 0)
	for _, gvk := range GetRegisterKinds() {
		rules = append(rules, rbacv1.PolicyRule{
			Namespace:        namespace,
			GroupVersionKind: gvk,
			Verbs:            allSupportVerbs(),
		})
	}
	// pods/logs
	gvk, verbs := podsNativeKinds()
	rules = append(rules, rbacv1.PolicyRule{
		Namespace:        namespace,
		GroupVersionKind: gvk,
		Verbs:            verbs,
	})
	// metrics
	gvks, verbs := metricsKinds()
	for _, gvk := range gvks {
		rules = append(rules, rbacv1.PolicyRule{
			Namespace:        namespace,
			GroupVersionKind: gvk,
			Verbs:            verbs,
		})
	}
	return rules
}

func (r *Resources) registerRoles() {
	rules := make([]rbacv1.PolicyRule, 0)
	for _, ns := range []string{"kube-neverdown", "kube-authx", "kube-discovery", "kube-system"} {
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
	zaplogger.Sugar().Infof("Creating ClusterRole:%s appid:%s", role.Name, r.AppId)

	if err := r.Authority.GenericApp.Role().Create(r.ctx, role); err != nil {
		if !errors.IsAlreadyExists(err) {
			zaplogger.Sugar().Fatal(err)
			return
		}
		zaplogger.Sugar().Infof("Already exist Role:%s appid:%s", role.Name, r.AppId)
	} else {
		zaplogger.Sugar().Infof("Successful create Role:%s appid:%s", role.Name, r.AppId)
	}
	//if err := r.Authority.GenericApp.Role().Update(r.ctx, role); err != nil {
	//	if errors.IsNotFound(err) {
	//		if err := r.Authority.GenericApp.Role().Create(r.ctx, role); err != nil {
	//			if !errors.IsAlreadyExists(err) {
	//				zaplogger.Sugar().Error(err)
	//			}
	//		}
	//	} else {
	//		zaplogger.Sugar().Error(err)
	//	}
	//}
}

func (r *Resources) registerServiceAccount() {
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
	defaultBackoff := wait.Backoff{
		Steps:    6,
		Duration: 200 * time.Millisecond,
		Factor:   5.0,
		Jitter:   0.1,
	}
	err := retry.OnError(defaultBackoff, errNotNil, func() error {
		sa, err := r.Authority.GenericApp.ServiceAccount().Get(sa.Spec.AccountMeta.Username)
		if err != nil {
			zaplogger.Sugar().Error(err)
			return err
		}
		r.savePasswordToSecrets(sa.Spec.AccountMeta.Password)
		return nil
	})
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
}

func errNotNil(err error) bool {
	return err != nil
}

func (r *Resources) savePasswordToSecrets(pwd string) {
	namespace := "kube-neverdown"
	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      "openx-password",
		},
		StringData: map[string]string{
			"password": pwd,
		},
		Type: corev1.SecretTypeOpaque,
	}
	if _, err := r.nativeClientSet.CoreV1().Secrets(namespace).Create(r.ctx, sec, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			zaplogger.Sugar().Fatal(err)
		}
		zaplogger.Sugar().Infof("Already exist Secrets:%s namespace:%s", sec.Name, namespace)
	} else {
		zaplogger.Sugar().Infof("Successful create Secrets:%s namespace:%s", sec.Name, namespace)
	}
}
