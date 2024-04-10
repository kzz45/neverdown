package groupversionkindrule

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

// groupversionkindruleStrategy implements behavior for Deployments.
type groupversionkindruleStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Deployment
// objects via the REST API.
var Strategy = groupversionkindruleStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// DefaultGarbageCollectionPolicy returns OrphanDependents for extensions/v1beta1, apps/v1beta1, and apps/v1beta2 for backwards compatibility,
// and DeleteDependents for all other versions.
func (groupversionkindruleStrategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
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
func (groupversionkindruleStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate clears fields that are not allowed to be set by end users on creation.
func (groupversionkindruleStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	groupversionkindrule := obj.(*rbac.GroupVersionKindRule)
	groupversionkindrule.Generation = 1
}

// Validate validates a new service.
func (groupversionkindruleStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (groupversionkindruleStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return []string{}
}

// Canonicalize normalizes the object after validation.
func (groupversionkindruleStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is false for deployments.
func (groupversionkindruleStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (groupversionkindruleStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newGroupVersionKindRule := obj.(*rbac.GroupVersionKindRule)
	oldGroupVersionKindRule := old.(*rbac.GroupVersionKindRule)

	if !equality.Semantic.DeepEqual(newGroupVersionKindRule.Spec, oldGroupVersionKindRule.Spec) ||
		!equality.Semantic.DeepEqual(newGroupVersionKindRule.Annotations, oldGroupVersionKindRule.Annotations) {
		newGroupVersionKindRule.Generation = oldGroupVersionKindRule.Generation + 1
	}
}

// ValidateUpdate is the default update validation for an end user.
func (groupversionkindruleStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (groupversionkindruleStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return []string{}
}

func (groupversionkindruleStrategy) AllowUnconditionalUpdate() bool {
	return true
}
