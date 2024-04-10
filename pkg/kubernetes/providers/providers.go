package providers

import (
	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	"github.com/kzz45/neverdown/pkg/kubernetes/api/legacyscheme"
	jingxrest "github.com/kzz45/neverdown/pkg/kubernetes/registry/jingx/rest"
	rbacrest "github.com/kzz45/neverdown/pkg/kubernetes/registry/rbac/rest"

	"github.com/kzz45/neverdown/pkg/apiserver/registry/generic"
	"github.com/kzz45/neverdown/pkg/apiserver/server"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/server/storage"
)

type RESTStorageProvider interface {
	GroupName() string
	NewRESTStorage(apiResourceConfigSource storage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (server.APIGroupInfo, error)
}

func LegacyCodec() runtime.Codec {
	return legacyscheme.Codecs.LegacyCodec(
		rbacv1.SchemeGroupVersion,
		jingxv1.SchemeGroupVersion,
	)
}

var RESTStorageProviders = []RESTStorageProvider{
	rbacrest.StorageProvider{},
	jingxrest.StorageProvider{},
}
