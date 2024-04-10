package event

import (
	context "context"

	"github.com/kzz45/neverdown/pkg/apis/jingx"
	"github.com/kzz45/neverdown/pkg/kubernetes/api/legacyscheme"

	"k8s.io/apiserver/pkg/registry/rest"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/storage/names"
)

// eventStrategy implements behavior for Deployments.
type eventStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Deployment
// objects via the REST API.
var Strategy = eventStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// DefaultGarbageCollectionPolicy returns OrphanDependents for extensions/v1beta1, apps/v1beta1, and apps/v1beta2 for backwards compatibility,
// and DeleteDependents for all other versions.
func (eventStrategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
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
func (eventStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate clears fields that are not allowed to be set by end users on creation.
func (eventStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	event := obj.(*jingx.Event)
	event.Generation = 1
}

// Validate validates a new service.
func (eventStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (eventStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return []string{}
}

// Canonicalize normalizes the object after validation.
func (eventStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is false for deployments.
func (eventStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (eventStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newEvent := obj.(*jingx.Event)
	oldEvent := old.(*jingx.Event)

	if !equality.Semantic.DeepEqual(newEvent.Spec, oldEvent.Spec) ||
		!equality.Semantic.DeepEqual(newEvent.Annotations, oldEvent.Annotations) {
		newEvent.Generation = oldEvent.Generation + 1
	}
}

// ValidateUpdate is the default update validation for an end user.
func (eventStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (eventStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return []string{}
}

func (eventStrategy) AllowUnconditionalUpdate() bool {
	return true
}
