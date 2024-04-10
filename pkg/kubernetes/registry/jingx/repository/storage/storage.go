package storage

import (
	"github.com/kzz45/neverdown/pkg/apis/jingx"
	"github.com/kzz45/neverdown/pkg/apiserver/registry/generic"
	"github.com/kzz45/neverdown/pkg/apiserver/registry/generic/registry"
	"github.com/kzz45/neverdown/pkg/kubernetes/printers"
	"github.com/kzz45/neverdown/pkg/kubernetes/printers/internalversion"
	"github.com/kzz45/neverdown/pkg/kubernetes/printers/storage"
	"github.com/kzz45/neverdown/pkg/kubernetes/registry/jingx/repository"

	"k8s.io/apiserver/pkg/registry/rest"

	"k8s.io/apimachinery/pkg/runtime"
)

// RepositoryStorage includes dummy storage for Services and for Scale subresource.
type RepositoryStorage struct {
	Repository *REST
}

// NewStorage returns new instance of ServiceStorage.
func NewStorage(optsGetter generic.RESTOptionsGetter) (RepositoryStorage, error) {
	serviceRest, err := NewREST(optsGetter)
	if err != nil {
		return RepositoryStorage{}, err
	}

	return RepositoryStorage{
		Repository: serviceRest,
	}, nil
}

// REST implements a RESTStorage for Deployments.
type REST struct {
	*registry.Store
	categories []string
}

// NewREST returns a RESTStorage object that will work against Repository.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &registry.Store{
		NewFunc:                   func() runtime.Object { return &jingx.Repository{} },
		NewListFunc:               func() runtime.Object { return &jingx.RepositoryList{} },
		DefaultQualifiedResource:  jingx.Resource("repositories"),
		SingularQualifiedResource: jingx.Resource("repository"),
		CreateStrategy:            repository.Strategy,
		UpdateStrategy:            repository.Strategy,
		DeleteStrategy:            repository.Strategy,

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
	return []string{"repository"}
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
