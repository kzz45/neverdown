package storage

import (
	"github.com/kzz45/neverdown/pkg/apis/jingx"
	"github.com/kzz45/neverdown/pkg/apiserver/registry/generic"
	"github.com/kzz45/neverdown/pkg/apiserver/registry/generic/registry"
	"github.com/kzz45/neverdown/pkg/kubernetes/printers"
	"github.com/kzz45/neverdown/pkg/kubernetes/printers/internalversion"
	"github.com/kzz45/neverdown/pkg/kubernetes/printers/storage"
	"github.com/kzz45/neverdown/pkg/kubernetes/registry/jingx/project"

	"k8s.io/apiserver/pkg/registry/rest"

	"k8s.io/apimachinery/pkg/runtime"
)

// ProjectStorage includes dummy storage for Services and for Scale subresource.
type ProjectStorage struct {
	Project *REST
}

// NewStorage returns new instance of ServiceStorage.
func NewStorage(optsGetter generic.RESTOptionsGetter) (ProjectStorage, error) {
	serviceRest, err := NewREST(optsGetter)
	if err != nil {
		return ProjectStorage{}, err
	}

	return ProjectStorage{
		Project: serviceRest,
	}, nil
}

// REST implements a RESTStorage for Deployments.
type REST struct {
	*registry.Store
	categories []string
}

// NewREST returns a RESTStorage object that will work against Project.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &registry.Store{
		NewFunc:                   func() runtime.Object { return &jingx.Project{} },
		NewListFunc:               func() runtime.Object { return &jingx.ProjectList{} },
		DefaultQualifiedResource:  jingx.Resource("projects"),
		SingularQualifiedResource: jingx.Resource("project"),
		CreateStrategy:            project.Strategy,
		UpdateStrategy:            project.Strategy,
		DeleteStrategy:            project.Strategy,

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
	return []string{"project"}
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
