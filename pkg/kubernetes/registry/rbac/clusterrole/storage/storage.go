package storage

import (
	rbac "github.com/kzz45/neverdown/pkg/apis/rbac"
	generic "github.com/kzz45/neverdown/pkg/apiserver/registry/generic"
	registry "github.com/kzz45/neverdown/pkg/apiserver/registry/generic/registry"
	rest "github.com/kzz45/neverdown/pkg/apiserver/registry/rest"
	printers "github.com/kzz45/neverdown/pkg/kubernetes/printers"
	internalversion "github.com/kzz45/neverdown/pkg/kubernetes/printers/internalversion"
	storage "github.com/kzz45/neverdown/pkg/kubernetes/printers/storage"
	clusterrole "github.com/kzz45/neverdown/pkg/kubernetes/registry/rbac/clusterrole"

	runtime "k8s.io/apimachinery/pkg/runtime"
)

// ClusterRoleStorage includes dummy storage for Services and for Scale subresource.
type ClusterRoleStorage struct {
	ClusterRole *REST
}

// NewStorage returns new instance of ServiceStorage.
func NewStorage(optsGetter generic.RESTOptionsGetter) (ClusterRoleStorage, error) {
	serviceRest, err := NewREST(optsGetter)
	if err != nil {
		return ClusterRoleStorage{}, err
	}

	return ClusterRoleStorage{
		ClusterRole: serviceRest,
	}, nil
}

// REST implements a RESTStorage for Deployments.
type REST struct {
	*registry.Store
	categories []string
}

// NewREST returns a RESTStorage object that will work against ClusterRole.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &registry.Store{
		NewFunc:                   func() runtime.Object { return &rbac.ClusterRole{} },
		NewListFunc:               func() runtime.Object { return &rbac.ClusterRoleList{} },
		DefaultQualifiedResource:  rbac.Resource("clusterroles"),
		SingularQualifiedResource: rbac.Resource("clusterrole"),
		CreateStrategy:            clusterrole.Strategy,
		UpdateStrategy:            clusterrole.Strategy,
		DeleteStrategy:            clusterrole.Strategy,

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
	return []string{"clusterrole"}
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
