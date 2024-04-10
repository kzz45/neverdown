package storage

import (
	"github.com/kzz45/neverdown/pkg/apis/rbac"

	"github.com/kzz45/neverdown/pkg/apiserver/registry/generic"
	"github.com/kzz45/neverdown/pkg/apiserver/registry/generic/registry"

	"github.com/kzz45/neverdown/pkg/kubernetes/printers"
	"github.com/kzz45/neverdown/pkg/kubernetes/printers/internalversion"
	"github.com/kzz45/neverdown/pkg/kubernetes/printers/storage"
	app "github.com/kzz45/neverdown/pkg/kubernetes/registry/rbac/app"

	"k8s.io/apiserver/pkg/registry/rest"

	runtime "k8s.io/apimachinery/pkg/runtime"
)

// AppStorage includes dummy storage for Services and for Scale subresource.
type AppStorage struct {
	App *REST
}

// NewStorage returns new instance of ServiceStorage.
func NewStorage(optsGetter generic.RESTOptionsGetter) (AppStorage, error) {
	serviceRest, err := NewREST(optsGetter)
	if err != nil {
		return AppStorage{}, err
	}

	return AppStorage{
		App: serviceRest,
	}, nil
}

// REST implements a RESTStorage for Deployments.
type REST struct {
	*registry.Store
	categories []string
}

// NewREST returns a RESTStorage object that will work against Recorder.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &registry.Store{
		NewFunc:                   func() runtime.Object { return &rbac.App{} },
		NewListFunc:               func() runtime.Object { return &rbac.AppList{} },
		DefaultQualifiedResource:  rbac.Resource("apps"),
		SingularQualifiedResource: rbac.Resource("app"),
		CreateStrategy:            app.Strategy,
		UpdateStrategy:            app.Strategy,
		DeleteStrategy:            app.Strategy,

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
	return []string{"app"}
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
