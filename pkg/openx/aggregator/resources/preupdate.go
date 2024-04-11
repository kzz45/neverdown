package resources

import (
	"errors"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
)

var ErrModifyReplicasWhenTheHpaOpened = errors.New("err: the replicas of the Application (in openx) could not be modified when the hpa had been opened")
var ErrMinReplicasAboutTheHpa = errors.New("err: the minReplicas of the HPA should be greater than or equal to 1")
var ErrMaxReplicasAboutTheHpa = errors.New("err: the maxReplicas of the HPA should be greater than or equal to minReplicas")

func (r *Resources) prepareUpdateOpenx(openx *openxv1.Openx) error {
	cache, err := r.OpenXInformerFactory.Openx().V1().Openxes().Lister().Openxes(openx.Namespace).Get(openx.Name)
	if err != nil {
		return err
	}
	apps := make(map[string]openxv1.App, 0)
	for _, v := range cache.Spec.Applications {
		apps[v.AppName] = v
	}
	for _, a1 := range openx.Spec.Applications {
		a2, ok := apps[a1.AppName]
		if !ok {
			continue
		}
		if a2.HorizontalPodAutoscalerSpec != nil && *a1.Replicas != *a2.Replicas {
			return ErrModifyReplicasWhenTheHpaOpened
		}
		if a1.HorizontalPodAutoscalerSpec != nil {
			if *a1.HorizontalPodAutoscalerSpec.MinReplicas <= 0 {
				return ErrMinReplicasAboutTheHpa
			}
			if *a1.HorizontalPodAutoscalerSpec.MinReplicas > *a1.HorizontalPodAutoscalerSpec.MinReplicas {
				return ErrMaxReplicasAboutTheHpa
			}
		}
	}
	return nil
}
