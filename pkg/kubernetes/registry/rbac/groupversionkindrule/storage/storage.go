package storage

import (
	rbac "github.com/kzz45/neverdown/pkg/apis/rbac"
	generic "github.com/kzz45/neverdown/pkg/apiserver/registry/generic"
	registry "github.com/kzz45/neverdown/pkg/apiserver/registry/generic/registry"
	rest "github.com/kzz45/neverdown/pkg/apiserver/registry/rest"
	printers "github.com/kzz45/neverdown/pkg/kubernetes/printers"
	internalversion "github.com/kzz45/neverdown/pkg/kubernetes/printers/internalversion"
	storage "github.com/kzz45/neverdown/pkg/kubernetes/printers/storage"
	groupversionkindrule "github.com/kzz45/neverdown/pkg/kubernetes/registry/rbac/groupversionkindrule"

	runtime "k8s.io/apimachinery/pkg/runtime"
)

// GroupVersionKindRuleStorage includes dummy storage for Services and for Scale subresource.
type GroupVersionKindRuleStorage struct {
	GroupVersionKindRule *REST
}

// NewStorage returns new instance of ServiceStorage.
func NewStorage(optsGetter generic.RESTOptionsGetter) (GroupVersionKindRuleStorage, error) {
	serviceRest, err := NewREST(optsGetter)
	if err != nil {
		return GroupVersionKindRuleStorage{}, err
	}

	return GroupVersionKindRuleStorage{
		GroupVersionKindRule: serviceRest,
	}, nil
}

// REST implements a RESTStorage for Deployments.
type REST struct {
	*registry.Store
	categories []string
}

// NewREST returns a RESTStorage object that will work against GroupVersionKindRule.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &registry.Store{
		NewFunc:                   func() runtime.Object { return &rbac.GroupVersionKindRule{} },
		NewListFunc:               func() runtime.Object { return &rbac.GroupVersionKindRuleList{} },
		DefaultQualifiedResource:  rbac.Resource("groupversionkindrules"),
		SingularQualifiedResource: rbac.Resource("groupversionkindrule"),
		CreateStrategy:            groupversionkindrule.Strategy,
		UpdateStrategy:            groupversionkindrule.Strategy,
		DeleteStrategy:            groupversionkindrule.Strategy,

		TableConvertor: storage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(internalversion.AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &REST{store, []string{"all"}}, nil
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"groupversionkindrule"}
}

// Implement CategoriesProvider
var _ rest.CategoriesProvider = &REST{}

// Categories implements the CategoriesProvider interface. Returns a list of categories a resource is part of.
func (r *REST) Categories() []string {
	return r.categories
}

// WithCategories sets categories for REST.
func (r *REST) WithCategories(categories []string) *REST {
	r.categories = categories
	return r
}
