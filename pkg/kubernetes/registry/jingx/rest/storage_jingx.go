package rest

import (
	"github.com/kzz45/neverdown/pkg/apis/jingx"
	v1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	"github.com/kzz45/neverdown/pkg/kubernetes/api/legacyscheme"

	eventstorage "github.com/kzz45/neverdown/pkg/kubernetes/registry/jingx/event/storage"
	projectstorage "github.com/kzz45/neverdown/pkg/kubernetes/registry/jingx/project/storage"
	repositorystorage "github.com/kzz45/neverdown/pkg/kubernetes/registry/jingx/repository/storage"
	tagstorage "github.com/kzz45/neverdown/pkg/kubernetes/registry/jingx/tag/storage"

	"github.com/kzz45/neverdown/pkg/apiserver/registry/generic"
	"github.com/kzz45/neverdown/pkg/apiserver/registry/rest"
	genericapiserver "github.com/kzz45/neverdown/pkg/apiserver/server"

	serverstorage "k8s.io/apiserver/pkg/server/storage"
)

// StorageProvider is a REST storage provider for discovery.k8s.io.
type StorageProvider struct{}

// NewRESTStorage returns a new storage provider.
func (p StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(jingx.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)

	if storageMap, err := p.v1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
		return genericapiserver.APIGroupInfo{}, err
	} else if len(storageMap) > 0 {
		apiGroupInfo.VersionedResourcesStorageMap[v1.SchemeGroupVersion.Version] = storageMap
	}

	return apiGroupInfo, nil
}

func (p StorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	// events
	eventstorage, err := eventstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["events"] = eventstorage.Event

	// projects
	projectstorage, err := projectstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["projects"] = projectstorage.Project

	// repositories
	repositorystorage, err := repositorystorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["repositories"] = repositorystorage.Repository

	// tags
	tagstorage, err := tagstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["tags"] = tagstorage.Tag

	return storage, nil
}

// GroupName is the group name for the storage provider.
func (p StorageProvider) GroupName() string {
	return jingx.GroupName
}
