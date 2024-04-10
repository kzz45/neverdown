package clusterrole

import (
	context "context"

	rbac "github.com/kzz45/neverdown/pkg/apis/rbac"
	rest "github.com/kzz45/neverdown/pkg/apiserver/registry/rest"
	legacyscheme "github.com/kzz45/neverdown/pkg/kubernetes/api/legacyscheme"

	equality "k8s.io/apimachinery/pkg/api/equality"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	field "k8s.io/apimachinery/pkg/util/validation/field"
	request "k8s.io/apiserver/pkg/endpoints/request"
	names "k8s.io/apiserver/pkg/storage/names"
)

// clusterroleStrategy implements behavior for Deployments.
type clusterroleStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Deployment
// objects via the REST API.
var Strategy = clusterroleStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// DefaultGarbageCollectionPolicy returns OrphanDependents for extensions/v1beta1, apps/v1beta1, and apps/v1beta2 for backwards compatibility,
// and DeleteDependents for all other versions.
func (clusterroleStrategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	var groupVersion schema.GroupVersion
	if requestInfo, found := request.RequestInfoFrom(ctx); found {
		groupVersion = schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}
	}
	switch groupVersion {
	//case extensionsv1beta1.SchemeGroupVersion, appsv1beta1.SchemeGroupVersion, appsv1beta2.SchemeGroupVersion:
	// for back compatibility
	//return rest.OrphanDependents
	default:
		return rest.DeleteDependents
	}
}

// NamespaceScoped is true for service.
func (clusterroleStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate clears fields that are not allowed to be set by end users on creation.
func (clusterroleStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	clusterrole := obj.(*rbac.ClusterRole)
	clusterrole.Generation = 1
}

// Validate validates a new service.
func (clusterroleStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (clusterroleStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return []string{}
}

// Canonicalize normalizes the object after validation.
func (clusterroleStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is false for deployments.
func (clusterroleStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (clusterroleStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newClusterRole := obj.(*rbac.ClusterRole)
	oldClusterRole := old.(*rbac.ClusterRole)

	if !equality.Semantic.DeepEqual(newClusterRole.Spec, oldClusterRole.Spec) ||
		!equality.Semantic.DeepEqual(newClusterRole.Annotations, oldClusterRole.Annotations) {
		newClusterRole.Generation = oldClusterRole.Generation + 1
	}
}

// ValidateUpdate is the default update validation for an end user.
func (clusterroleStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (clusterroleStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return []string{}
}

func (clusterroleStrategy) AllowUnconditionalUpdate() bool {
	return true
}
