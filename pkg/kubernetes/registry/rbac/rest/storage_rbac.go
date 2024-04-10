package rest

import (
	"github.com/kzz45/neverdown/pkg/apis/rbac"
	v1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	"github.com/kzz45/neverdown/pkg/kubernetes/api/legacyscheme"
	appstorage "github.com/kzz45/neverdown/pkg/kubernetes/registry/rbac/app/storage"
	appserviceaccountstorage "github.com/kzz45/neverdown/pkg/kubernetes/registry/rbac/appserviceaccount/storage"
	clusterrolestorage "github.com/kzz45/neverdown/pkg/kubernetes/registry/rbac/clusterrole/storage"
	groupversionkindrulestorage "github.com/kzz45/neverdown/pkg/kubernetes/registry/rbac/groupversionkindrule/storage"
	rbacserviceaccountstorage "github.com/kzz45/neverdown/pkg/kubernetes/registry/rbac/rbacserviceaccount/storage"

	"github.com/kzz45/neverdown/pkg/apiserver/registry/generic"
	"github.com/kzz45/neverdown/pkg/apiserver/registry/rest"
	genericapiserver "github.com/kzz45/neverdown/pkg/apiserver/server"

	serverstorage "k8s.io/apiserver/pkg/server/storage"
)

// StorageProvider is a REST storage provider for discovery.k8s.io.
type StorageProvider struct{}

// NewRESTStorage returns a new storage provider.
func (p StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(rbac.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)

	if storageMap, err := p.v1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
		return genericapiserver.APIGroupInfo{}, err
	} else if len(storageMap) > 0 {
		apiGroupInfo.VersionedResourcesStorageMap[v1.SchemeGroupVersion.Version] = storageMap
	}

	return apiGroupInfo, nil
}

func (p StorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	// apps
	appStorage, err := appstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["apps"] = appStorage.App

	// appserviceaccounts
	appserviceaccountStorage, err := appserviceaccountstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["appserviceaccounts"] = appserviceaccountStorage.AppServiceAccount

	// clusterroles
	clusterroleStorage, err := clusterrolestorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["clusterroles"] = clusterroleStorage.ClusterRole

	// groupversionkindrules
	groupversionkindruleStorage, err := groupversionkindrulestorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["groupversionkindrules"] = groupversionkindruleStorage.GroupVersionKindRule

	// rbacserviceaccounts
	rbacserviceaccountStorage, err := rbacserviceaccountstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["rbacserviceaccounts"] = rbacserviceaccountStorage.RbacServiceAccount

	return storage, nil
}

// GroupName is the group name for the storage provider.
func (p StorageProvider) GroupName() string {
	return rbac.GroupName
}
