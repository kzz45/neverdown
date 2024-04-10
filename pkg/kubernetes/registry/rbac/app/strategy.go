package app

import (
	"context"

	"github.com/kzz45/neverdown/pkg/apis/rbac"
	"github.com/kzz45/neverdown/pkg/kubernetes/api/legacyscheme"

	"k8s.io/apiserver/pkg/registry/rest"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
)

// coffeeStrategy implements behavior for Deployments.
type appStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Deployment
// objects via the REST API.
var Strategy = appStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// Make sure we correctly implement the interface.
var _ = rest.GarbageCollectionDeleteStrategy(Strategy)

// DefaultGarbageCollectionPolicy returns DeleteDependents for all currently served versions.
func (appStrategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.DeleteDependents
}

// NamespaceScoped is true for deployment.
func (appStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate clears fields that are not allowed to be set by end users on creation.
func (appStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	app := obj.(*rbac.App)
	app.Generation = 1
}

// Validate validates a new deployment.
func (appStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (appStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	var warnings []string
	return warnings
}

// Canonicalize normalizes the object after validation.
func (appStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is false for deployments.
func (appStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (appStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newApp := obj.(*rbac.App)
	oldApp := old.(*rbac.App)

	// Spec updates bump the generation so that we can distinguish between
	// scaling events and template changes, annotation updates bump the generation
	// because annotations are copied from deployments to their replica sets.
	if !apiequality.Semantic.DeepEqual(newApp.Spec, oldApp.Spec) ||
		!apiequality.Semantic.DeepEqual(newApp.Annotations, oldApp.Annotations) {
		newApp.Generation = oldApp.Generation + 1
	}
}

// ValidateUpdate is the default update validation for an end user.
func (appStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (appStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	var warnings []string
	return warnings
}

func (appStrategy) AllowUnconditionalUpdate() bool {
	return true
}
