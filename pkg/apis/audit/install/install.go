package install

import (
	"github.com/kzz45/neverdown/pkg/apis/audit"
	v1 "github.com/kzz45/neverdown/pkg/apis/audit/v1"
	"github.com/kzz45/neverdown/pkg/kubernetes/api/legacyscheme"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func init() {
	Install(legacyscheme.Scheme)
}

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(audit.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1.SchemeGroupVersion))
}
