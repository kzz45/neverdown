package controller

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

type BaseControllerRefManager struct {
	Controller metav1.Object
	Selector   labels.Selector

	canAdoptErr  error
	canAdoptOnce sync.Once
	CanAdoptFunc func() error
}

func (m *BaseControllerRefManager) CanAdopt() error {
	m.canAdoptOnce.Do(func() {
		if m.CanAdoptFunc != nil {
			m.canAdoptErr = m.CanAdoptFunc()
		}
	})
	return m.canAdoptErr
}

// ClaimObject tries to take ownership of an object for this controller.
//
// It will reconcile the following:
//   - Adopt orphans if the match function returns true.
//   - Release owned objects if the match function returns false.
//
// A non-nil error is returned if some form of reconciliation was attempted and
// failed. Usually, controllers should try again later in case reconciliation
// is still needed.
//
// If the error is nil, either the reconciliation succeeded, or no
// reconciliation was necessary. The returned boolean indicates whether you now
// own the object.
//
// No reconciliation will be attempted if the controller is being deleted.
func (m *BaseControllerRefManager) ClaimObject(obj metav1.Object, match func(metav1.Object) bool, adopt, release func(metav1.Object) error) (bool, error) {
	controllerRef := metav1.GetControllerOfNoCopy(obj)
	if controllerRef != nil {
		if controllerRef.UID != m.Controller.GetUID() {
			// Owned by someone else. Ignore.
			return false, nil
		}
		if match(obj) {
			// We already own it and the selector matches.
			// Return true (successfully claimed) before checking deletion timestamp.
			// We're still allowed to claim things we already own while being deleted
			// because doing so requires taking no actions.
			return true, nil
		}
		// Owned by us but selector doesn't match.
		// Try to release, unless we're being deleted.
		if m.Controller.GetDeletionTimestamp() != nil {
			return false, nil
		}
		if err := release(obj); err != nil {
			// If the pod no longer exists, ignore the error.
			if errors.IsNotFound(err) {
				return false, nil
			}
			// Either someone else released it, or there was a transient error.
			// The controller should requeue and try again if it's still stale.
			return false, err
		}
		// Successfully released.
		return false, nil
	}

	// It's an orphan.
	if m.Controller.GetDeletionTimestamp() != nil || !match(obj) {
		// Ignore if we're being deleted or selector doesn't match.
		return false, nil
	}
	if obj.GetDeletionTimestamp() != nil {
		// Ignore if the object is being deleted
		return false, nil
	}
	// Selector matches. Try to adopt.
	if err := adopt(obj); err != nil {
		// If the pod no longer exists, ignore the error.
		if errors.IsNotFound(err) {
			return false, nil
		}
		// Either someone else claimed it first, or there was a transient error.
		// The controller should requeue and try again if it's still orphaned.
		return false, err
	}
	// Successfully adopted.
	return true, nil
}

// ServiceControllerRefManager is used to manage controllerRef of Service.
// Three methods are defined on this object 1: Classify 2: AdoptReplicaSet and
// 3: ReleaseReplicaSet which are used to classify the Service into appropriate
// categories and accordingly adopt or release them. See comments on these functions
// for more details.
type ServiceControllerRefManager struct {
	BaseControllerRefManager
	controllerKind schema.GroupVersionKind
	svcControl     ServiceControlInterface
}

// NewServiceControllerRefManager returns a ServiceControllerRefManager that exposes
// methods to manage the controllerRef of Services.
//
// The CanAdopt() function can be used to perform a potentially expensive check
// (such as a live GET from the API server) prior to the first adoption.
// It will only be called (at most once) if an adoption is actually attempted.
// If CanAdopt() returns a non-nil error, all adoptions will fail.
//
// NOTE: Once CanAdopt() is called, it will not be called again by the same
//
//	ReplicaSetControllerRefManager instance. Create a new instance if it
//	makes sense to check CanAdopt() again (e.g. in a different sync pass).
func NewServiceControllerRefManager(
	svcControl ServiceControlInterface,
	controller metav1.Object,
	selector labels.Selector,
	controllerKind schema.GroupVersionKind,
	canAdopt func() error,
) *ServiceControllerRefManager {
	return &ServiceControllerRefManager{
		BaseControllerRefManager: BaseControllerRefManager{
			Controller:   controller,
			Selector:     selector,
			CanAdoptFunc: canAdopt,
		},
		controllerKind: controllerKind,
		svcControl:     svcControl,
	}
}

func (m *ServiceControllerRefManager) ClaimServices(services []*corev1.Service) ([]*corev1.Service, error) {
	var claimed []*corev1.Service
	var errlist []error

	match := func(obj metav1.Object) bool {
		return m.Selector.Matches(labels.Set(obj.GetLabels()))
	}
	adopt := func(obj metav1.Object) error {
		return m.AdoptService(obj.(*corev1.Service))
	}
	release := func(obj metav1.Object) error {
		return m.ReleaseService(obj.(*corev1.Service))
	}

	for _, rs := range services {
		ok, err := m.ClaimObject(rs, match, adopt, release)
		if err != nil {
			errlist = append(errlist, err)
			continue
		}
		if ok {
			claimed = append(claimed, rs)
		}
	}
	return claimed, utilerrors.NewAggregate(errlist)
}

// AdoptService sends a patch to take control of the Service. It returns
// the error if the patching fails.
func (m *ServiceControllerRefManager) AdoptService(rs *corev1.Service) error {
	if err := m.CanAdopt(); err != nil {
		return fmt.Errorf("can't adopt Service %v/%v (%v): %v", rs.Namespace, rs.Name, rs.UID, err)
	}
	// Note that ValidateOwnerReferences() will reject this patch if another
	// OwnerReference exists with controller=true.
	patchBytes, err := ownerRefControllerPatch(m.Controller, m.controllerKind, rs.UID)
	if err != nil {
		return err
	}
	return m.svcControl.Patch(rs.Namespace, rs.Name, patchBytes)
}

// ReleaseService sends a patch to free the Service from the control of the Deployment controller.
// It returns the error if the patching fails. 404 and 422 errors are ignored.
func (m *ServiceControllerRefManager) ReleaseService(svc *corev1.Service) error {
	zaplogger.Sugar().Infof("patching Service %s_%s to remove its controllerRef to %s/%s:%s",
		svc.Namespace, svc.Name, m.controllerKind.GroupVersion(), m.controllerKind.Kind, m.Controller.GetName())
	patchBytes, err := deleteOwnerRefStrategicMergePatch(svc.UID, m.Controller.GetUID())
	if err != nil {
		return err
	}
	err = m.svcControl.Patch(svc.Namespace, svc.Name, patchBytes)
	if err != nil {
		if errors.IsNotFound(err) {
			// If the ReplicaSet no longer exists, ignore it.
			return nil
		}
		if errors.IsInvalid(err) {
			// Invalid error will be returned in two cases: 1. the ReplicaSet
			// has no owner reference, 2. the uid of the ReplicaSet doesn't
			// match, which means the ReplicaSet is deleted and then recreated.
			// In both cases, the error can be ignored.
			return nil
		}
	}
	return err
}

// StatefulSetControllerRefManager is used to manage controllerRef of Service.
// Three methods are defined on this object 1: Classify 2: AdoptReplicaSet and
// 3: ReleaseReplicaSet which are used to classify the Service into appropriate
// categories and accordingly adopt or release them. See comments on these functions
// for more details.
type StatefulSetControllerRefManager struct {
	BaseControllerRefManager
	controllerKind schema.GroupVersionKind
	stsControl     StatefulSetControlInterface
}

// NewStatefulSetControllerRefManager returns a StatefulSetControllerRefManager that exposes
// methods to manage the controllerRef of Services.
//
// The CanAdopt() function can be used to perform a potentially expensive check
// (such as a live GET from the API server) prior to the first adoption.
// It will only be called (at most once) if an adoption is actually attempted.
// If CanAdopt() returns a non-nil error, all adoptions will fail.
//
// NOTE: Once CanAdopt() is called, it will not be called again by the same
//
//	ReplicaSetControllerRefManager instance. Create a new instance if it
//	makes sense to check CanAdopt() again (e.g. in a different sync pass).
func NewStatefulSetControllerRefManager(
	stsControl StatefulSetControlInterface,
	controller metav1.Object,
	selector labels.Selector,
	controllerKind schema.GroupVersionKind,
	canAdopt func() error,
) *StatefulSetControllerRefManager {
	return &StatefulSetControllerRefManager{
		BaseControllerRefManager: BaseControllerRefManager{
			Controller:   controller,
			Selector:     selector,
			CanAdoptFunc: canAdopt,
		},
		controllerKind: controllerKind,
		stsControl:     stsControl,
	}
}

func (m *StatefulSetControllerRefManager) ClaimStatefulSets(stses []*appsv1.StatefulSet) ([]*appsv1.StatefulSet, error) {
	var claimed []*appsv1.StatefulSet
	var errlist []error

	match := func(obj metav1.Object) bool {
		return m.Selector.Matches(labels.Set(obj.GetLabels()))
	}
	adopt := func(obj metav1.Object) error {
		return m.AdoptStatefulSet(obj.(*appsv1.StatefulSet))
	}
	release := func(obj metav1.Object) error {
		return m.ReleaseStatefulSet(obj.(*appsv1.StatefulSet))
	}

	for _, sts := range stses {
		ok, err := m.ClaimObject(sts, match, adopt, release)
		if err != nil {
			errlist = append(errlist, err)
			continue
		}
		if ok {
			claimed = append(claimed, sts)
		}
	}
	return claimed, utilerrors.NewAggregate(errlist)
}

// AdoptStatefulSet sends a patch to take control of the StatefulSet. It returns
// the error if the patching fails.
func (m *StatefulSetControllerRefManager) AdoptStatefulSet(rs *appsv1.StatefulSet) error {
	if err := m.CanAdopt(); err != nil {
		return fmt.Errorf("can't adopt StatefulSet %v/%v (%v): %v", rs.Namespace, rs.Name, rs.UID, err)
	}
	// Note that ValidateOwnerReferences() will reject this patch if another
	// OwnerReference exists with controller=true.
	patchBytes, err := ownerRefControllerPatch(m.Controller, m.controllerKind, rs.UID)
	if err != nil {
		return err
	}
	return m.stsControl.Patch(rs.Namespace, rs.Name, patchBytes)
}

// ReleaseStatefulSet sends a patch to free the Service from the control of the Deployment controller.
// It returns the error if the patching fails. 404 and 422 errors are ignored.
func (m *StatefulSetControllerRefManager) ReleaseStatefulSet(sts *appsv1.StatefulSet) error {
	zaplogger.Sugar().Infof("patching StatefulSet %s_%s to remove its controllerRef to %s/%s:%s",
		sts.Namespace, sts.Name, m.controllerKind.GroupVersion(), m.controllerKind.Kind, m.Controller.GetName())
	patchBytes, err := deleteOwnerRefStrategicMergePatch(sts.UID, m.Controller.GetUID())
	if err != nil {
		return err
	}
	err = m.stsControl.Patch(sts.Namespace, sts.Name, patchBytes)
	if err != nil {
		if errors.IsNotFound(err) {
			// If the StatefulSet no longer exists, ignore it.
			return nil
		}
		if errors.IsInvalid(err) {
			// Invalid error will be returned in two cases: 1. the ReplicaSet
			// has no owner reference, 2. the uid of the ReplicaSet doesn't
			// match, which means the ReplicaSet is deleted and then recreated.
			// In both cases, the error can be ignored.
			return nil
		}
	}
	return err
}

// DeploymentControllerRefManager is used to manage controllerRef of Service.
// Three methods are defined on this object 1: Classify 2: AdoptReplicaSet and
// 3: ReleaseReplicaSet which are used to classify the Service into appropriate
// categories and accordingly adopt or release them. See comments on these functions
// for more details.
type DeploymentControllerRefManager struct {
	BaseControllerRefManager
	controllerKind schema.GroupVersionKind
	dpControl      DeploymentControlInterface
}

// NewDeploymentControllerRefManager returns a DeploymentControllerRefManager that exposes
// methods to manage the controllerRef of Services.
//
// The CanAdopt() function can be used to perform a potentially expensive check
// (such as a live GET from the API server) prior to the first adoption.
// It will only be called (at most once) if an adoption is actually attempted.
// If CanAdopt() returns a non-nil error, all adoptions will fail.
//
// NOTE: Once CanAdopt() is called, it will not be called again by the same
//
//	ReplicaSetControllerRefManager instance. Create a new instance if it
//	makes sense to check CanAdopt() again (e.g. in a different sync pass).
func NewDeploymentControllerRefManager(
	stsControl DeploymentControlInterface,
	controller metav1.Object,
	selector labels.Selector,
	controllerKind schema.GroupVersionKind,
	canAdopt func() error,
) *DeploymentControllerRefManager {
	return &DeploymentControllerRefManager{
		BaseControllerRefManager: BaseControllerRefManager{
			Controller:   controller,
			Selector:     selector,
			CanAdoptFunc: canAdopt,
		},
		controllerKind: controllerKind,
		dpControl:      stsControl,
	}
}

func (m *DeploymentControllerRefManager) ClaimDeployments(stses []*appsv1.Deployment) ([]*appsv1.Deployment, error) {
	var claimed []*appsv1.Deployment
	var errlist []error

	match := func(obj metav1.Object) bool {
		return m.Selector.Matches(labels.Set(obj.GetLabels()))
	}
	adopt := func(obj metav1.Object) error {
		return m.AdoptDeployment(obj.(*appsv1.Deployment))
	}
	release := func(obj metav1.Object) error {
		return m.ReleaseDeployment(obj.(*appsv1.Deployment))
	}

	for _, sts := range stses {
		ok, err := m.ClaimObject(sts, match, adopt, release)
		if err != nil {
			errlist = append(errlist, err)
			continue
		}
		if ok {
			claimed = append(claimed, sts)
		}
	}
	return claimed, utilerrors.NewAggregate(errlist)
}

// AdoptDeployment sends a patch to take control of the Deployment. It returns
// the error if the patching fails.
func (m *DeploymentControllerRefManager) AdoptDeployment(rs *appsv1.Deployment) error {
	if err := m.CanAdopt(); err != nil {
		return fmt.Errorf("can't adopt Deployment %v/%v (%v): %v", rs.Namespace, rs.Name, rs.UID, err)
	}
	// Note that ValidateOwnerReferences() will reject this patch if another
	// OwnerReference exists with controller=true.
	patchBytes, err := ownerRefControllerPatch(m.Controller, m.controllerKind, rs.UID)
	if err != nil {
		return err
	}
	return m.dpControl.Patch(rs.Namespace, rs.Name, patchBytes)
}

// ReleaseDeployment sends a patch to free the Service from the control of the Deployment controller.
// It returns the error if the patching fails. 404 and 422 errors are ignored.
func (m *DeploymentControllerRefManager) ReleaseDeployment(sts *appsv1.Deployment) error {
	zaplogger.Sugar().Infof("patching Deployment %s_%s to remove its controllerRef to %s/%s:%s",
		sts.Namespace, sts.Name, m.controllerKind.GroupVersion(), m.controllerKind.Kind, m.Controller.GetName())
	patchBytes, err := deleteOwnerRefStrategicMergePatch(sts.UID, m.Controller.GetUID())
	if err != nil {
		return err
	}
	err = m.dpControl.Patch(sts.Namespace, sts.Name, patchBytes)
	if err != nil {
		if errors.IsNotFound(err) {
			// If the Deployment no longer exists, ignore it.
			return nil
		}
		if errors.IsInvalid(err) {
			// Invalid error will be returned in two cases: 1. the ReplicaSet
			// has no owner reference, 2. the uid of the ReplicaSet doesn't
			// match, which means the ReplicaSet is deleted and then recreated.
			// In both cases, the error can be ignored.
			return nil
		}
	}
	return err
}

type HorizontalPodAutoscalerRefManager struct {
	BaseControllerRefManager
	controllerKind schema.GroupVersionKind
	hpaControl     HorizontalPodAutoscalerControlInterface
}

func NewHorizontalPodAutoscalerControllerRefManager(
	stsControl HorizontalPodAutoscalerControlInterface,
	controller metav1.Object,
	selector labels.Selector,
	controllerKind schema.GroupVersionKind,
	canAdopt func() error,
) *HorizontalPodAutoscalerRefManager {
	return &HorizontalPodAutoscalerRefManager{
		BaseControllerRefManager: BaseControllerRefManager{
			Controller:   controller,
			Selector:     selector,
			CanAdoptFunc: canAdopt,
		},
		controllerKind: controllerKind,
		hpaControl:     stsControl,
	}
}

func (m *HorizontalPodAutoscalerRefManager) ClaimHorizontalPodAutoscalers(stses []*autoscalingv2.HorizontalPodAutoscaler) ([]*autoscalingv2.HorizontalPodAutoscaler, error) {
	var claimed []*autoscalingv2.HorizontalPodAutoscaler
	var errlist []error

	match := func(obj metav1.Object) bool {
		return m.Selector.Matches(labels.Set(obj.GetLabels()))
	}
	adopt := func(obj metav1.Object) error {
		return m.AdoptHorizontalPodAutoscaler(obj.(*autoscalingv2.HorizontalPodAutoscaler))
	}
	release := func(obj metav1.Object) error {
		return m.ReleaseHorizontalPodAutoscaler(obj.(*autoscalingv2.HorizontalPodAutoscaler))
	}

	for _, sts := range stses {
		ok, err := m.ClaimObject(sts, match, adopt, release)
		if err != nil {
			errlist = append(errlist, err)
			continue
		}
		if ok {
			claimed = append(claimed, sts)
		}
	}
	return claimed, utilerrors.NewAggregate(errlist)
}

func (m *HorizontalPodAutoscalerRefManager) AdoptHorizontalPodAutoscaler(rs *autoscalingv2.HorizontalPodAutoscaler) error {
	if err := m.CanAdopt(); err != nil {
		return fmt.Errorf("can't adopt HorizontalPodAutoscaler %v/%v (%v): %v", rs.Namespace, rs.Name, rs.UID, err)
	}
	// Note that ValidateOwnerReferences() will reject this patch if another
	// OwnerReference exists with controller=true.
	patchBytes, err := ownerRefControllerPatch(m.Controller, m.controllerKind, rs.UID)
	if err != nil {
		return err
	}
	return m.hpaControl.Patch(rs.Namespace, rs.Name, patchBytes)
}

func (m *HorizontalPodAutoscalerRefManager) ReleaseHorizontalPodAutoscaler(sts *autoscalingv2.HorizontalPodAutoscaler) error {
	zaplogger.Sugar().Infof("patching HorizontalPodAutoscaler %s_%s to remove its controllerRef to %s/%s:%s",
		sts.Namespace, sts.Name, m.controllerKind.GroupVersion(), m.controllerKind.Kind, m.Controller.GetName())
	patchBytes, err := deleteOwnerRefStrategicMergePatch(sts.UID, m.Controller.GetUID())
	if err != nil {
		return err
	}
	err = m.hpaControl.Patch(sts.Namespace, sts.Name, patchBytes)
	if err != nil {
		if errors.IsNotFound(err) {
			// If the HorizontalPodAutoscaler no longer exists, ignore it.
			return nil
		}
		if errors.IsInvalid(err) {
			// Invalid error will be returned in two cases: 1. the ReplicaSet
			// has no owner reference, 2. the uid of the ReplicaSet doesn't
			// match, which means the ReplicaSet is deleted and then recreated.
			// In both cases, the error can be ignored.
			return nil
		}
	}
	return err
}

type objectForDeleteOwnerRefStrategicMergePatch struct {
	Metadata objectMetaForMergePatch `json:"metadata"`
}

type objectMetaForMergePatch struct {
	UID             types.UID           `json:"uid"`
	OwnerReferences []map[string]string `json:"ownerReferences"`
}

func deleteOwnerRefStrategicMergePatch(dependentUID types.UID, ownerUIDs ...types.UID) ([]byte, error) {
	var pieces []map[string]string
	for _, ownerUID := range ownerUIDs {
		pieces = append(pieces, map[string]string{"$patch": "delete", "uid": string(ownerUID)})
	}
	patch := objectForDeleteOwnerRefStrategicMergePatch{
		Metadata: objectMetaForMergePatch{
			UID:             dependentUID,
			OwnerReferences: pieces,
		},
	}
	patchBytes, err := json.Marshal(&patch)
	if err != nil {
		return nil, err
	}
	return patchBytes, nil
}

type objectForAddOwnerRefPatch struct {
	Metadata objectMetaForPatch `json:"metadata"`
}

type objectMetaForPatch struct {
	OwnerReferences []metav1.OwnerReference `json:"ownerReferences"`
	UID             types.UID               `json:"uid"`
}

func ownerRefControllerPatch(controller metav1.Object, controllerKind schema.GroupVersionKind, uid types.UID) ([]byte, error) {
	blockOwnerDeletion := true
	isController := true
	addControllerPatch := objectForAddOwnerRefPatch{
		Metadata: objectMetaForPatch{
			UID: uid,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         controllerKind.GroupVersion().String(),
					Kind:               controllerKind.Kind,
					Name:               controller.GetName(),
					UID:                controller.GetUID(),
					Controller:         &isController,
					BlockOwnerDeletion: &blockOwnerDeletion,
				},
			},
		},
	}
	patchBytes, err := json.Marshal(&addControllerPatch)
	if err != nil {
		return nil, err
	}
	return patchBytes, nil
}
