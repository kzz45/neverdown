package controlplane

import (
	"fmt"

	"github.com/kzz45/neverdown/pkg/client-go/clientset/versioned/scheme"
	"github.com/kzz45/neverdown/pkg/kubernetes/api/legacyscheme"
	"github.com/kzz45/neverdown/pkg/kubernetes/providers"

	"github.com/kzz45/neverdown/pkg/apiserver/registry/generic"
	"github.com/kzz45/neverdown/pkg/apiserver/server"
	"github.com/kzz45/neverdown/pkg/apiserver/server/options"

	serverstorage "k8s.io/apiserver/pkg/server/storage"

	"k8s.io/apiserver/pkg/storage/storagebackend"
	"k8s.io/klog/v2"
)

type Instance struct {
	GenericAPIServer *server.GenericAPIServer

	restOptionsGetter generic.RESTOptionsGetter
}

func New(storeConfig *storagebackend.Config) *Instance {
	storeConfig.Codec = providers.LegacyCodec()
	opts := options.NewEtcdOptions(storeConfig)

	cRESTOptionsGetter := &options.SimpleRestOptionsFactory{Options: *opts}

	c := server.NewConfig(scheme.Codecs)
	c.Complete()

	ic := server.AConfig{}
	ic.Config = c
	s, err := ic.New("discovery-controlplane", server.NewEmptyDelegate())
	if err != nil {
		klog.Fatal(err)
	}

	i := &Instance{
		GenericAPIServer:  s,
		restOptionsGetter: cRESTOptionsGetter,
	}
	if err := i.InstallAPIs(nil, nil, providers.RESTStorageProviders...); err != nil {
		klog.Fatal(err)
	}
	return i
}

// InstallAPIs will install the APIs for the restStorageProviders if they are enabled.
func (m *Instance) InstallAPIs(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, restStorageProviders ...providers.RESTStorageProvider) error {
	klog.Info(legacyscheme.Scheme.AllKnownTypes())

	apiGroupsInfo := make([]*server.APIGroupInfo, 0)
	for _, restStorageBuilder := range restStorageProviders {
		apiGroupInfo, err := restStorageBuilder.NewRESTStorage(apiResourceConfigSource, m.restOptionsGetter)
		if err != nil {
			klog.Errorf("problem initializing API group %q : %v", restStorageBuilder.GroupName(), err)
			return err
		}
		apiGroupsInfo = append(apiGroupsInfo, &apiGroupInfo)
	}
	if err := m.GenericAPIServer.InstallAPIGroups(apiGroupsInfo...); err != nil {
		return fmt.Errorf("error in registering group versions: %v", err)
	}
	return nil
}
