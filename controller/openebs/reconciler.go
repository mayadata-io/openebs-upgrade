/*
Copyright 2020 The MayaData Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openebs

import (
	"encoding/json"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/pkg/utils/metac"
	"mayadata.io/openebs-upgrade/types"
	"openebs.io/metac/controller/generic"
)

type reconcileErrHandler struct {
	openebs      *unstructured.Unstructured
	hookResponse *generic.SyncHookResponse
}

type reconcileSuccessHandler struct {
	openebs      *unstructured.Unstructured
	hookResponse *generic.SyncHookResponse
}

func (h *reconcileErrHandler) handle(err error) {
	// Error has been handled elaborately. This logic ensures
	// error message is propagated to the resource & hence seen via
	// 'kubectl get OpenEBS -oyaml'.
	//
	// In addition, logging has been done to check for error messages
	// from this pod's logs.
	glog.Errorf(
		"Failed to reconcile OpenEBS %s %s: %+v",
		h.openebs.GetNamespace(), h.openebs.GetName(), err,
	)
	// response status will be set against the watch's status by metac
	h.hookResponse.Status = map[string]interface{}{}
	h.hookResponse.Status["phase"] = types.OpenEBSStatusPhaseFailed
	h.hookResponse.Status["reason"] = err.Error()
	// this will stop further reconciliation at metac since there was
	// an error
	h.hookResponse.SkipReconcile = true
}

// This method is being used for setting the success response
// against watch's status by metac.
func (h *reconcileSuccessHandler) handle() {
	// response status will be set against the watch's status by metac
	h.hookResponse.Status = map[string]interface{}{}
	h.hookResponse.Status["phase"] = types.OpenEBSStatusPhaseOnline
}

// Sync implements the idempotent logic to reconcile OpenEBS
//
// NOTE:
// 	SyncHookRequest is the payload received as part of reconcile
// request. Similarly, SyncHookResponse is the payload sent as a
// response as part of reconcile request.
//
// NOTE:
//	SyncHookRequest uses OpenEBS as the watched resource.
// SyncHookResponse has the resources that forms the desired state
// w.r.t the watched resource.
//
// NOTE:
//	Returning error will panic this process. We would rather want this
// controller to run continuously. Hence, the errors are logged and at
// the same time, these errors are posted against OpenEBS's
// status.
func Sync(request *generic.SyncHookRequest, response *generic.SyncHookResponse) error {
	if request == nil {
		return errors.Errorf("Failed to reconcile OpenEBS: Nil request found")
	}
	if response == nil {
		return errors.Errorf("Failed to reconcile OpenEBS: Nil response found")
	}
	// Nothing needs to be done if there are no attachments in request
	//
	// NOTE:
	// 	It is expected to have OpenEBS as an attachment
	// resource as well as the resource under watch.
	if request.Attachments == nil || request.Attachments.IsEmpty() {
		response.SkipReconcile = true
		return nil
	}

	glog.V(3).Infof(
		"Will reconcile OpenEBS %s %s:",
		request.Watch.GetNamespace(), request.Watch.GetName(),
	)

	// construct the error handler
	errHandler := &reconcileErrHandler{
		openebs:      request.Watch,
		hookResponse: response,
	}

	var observedOpenEBS *unstructured.Unstructured
	var observedOpenEBSComponents []*unstructured.Unstructured
	// observedOpenEBSCRDs will store the details of all the OpenEBS
	// related CRDs present in the cluster.
	var observedOpenEBSCRDs []*unstructured.Unstructured
	// observedOpenEBSClusterRoleAndRoleBindings will store the details of all the OpenEBS
	// related cluster role and cluster role bindings present in the cluster.
	var observedOpenEBSClusterRoleAndRoleBindings []*unstructured.Unstructured
	var observedCStorCSIDriver *unstructured.Unstructured
	for _, attachment := range request.Attachments.List() {
		// this watch resource must be present in the list of attachments
		if request.Watch.GetUID() == attachment.GetUID() &&
			attachment.GetKind() == string(types.KindOpenEBS) {
			// this is the required OpenEBS
			observedOpenEBS = attachment
			// add this to the response later after completion of its
			// reconcile logic
			continue
		}
		// If the attachments are not of Kind: OpenEBS then it will be
		// considered as an OpenEBS component.
		if attachment.GetKind() != string(types.KindOpenEBS) {
			observedOpenEBSComponents = append(observedOpenEBSComponents, attachment)
		}
		if attachment.GetKind() == types.KindCustomResourceDefinition {
			observedOpenEBSCRDs = append(observedOpenEBSCRDs, attachment)
		}
		if attachment.GetKind() == types.KindClusterRole || attachment.GetKind() == types.KindClusterRoleBinding {
			observedOpenEBSClusterRoleAndRoleBindings = append(observedOpenEBSClusterRoleAndRoleBindings, attachment)
		}
		if attachment.GetKind() == types.KindCSIDriver &&
			attachment.GetName() == types.CStorCSIDriverNameKey {
			observedCStorCSIDriver = attachment
		}
	}

	if observedOpenEBS == nil {
		errHandler.handle(
			errors.Errorf("Can't reconcile: OpenEBS not found in attachments"),
		)
		return nil
	}

	// reconciler is the one that will perform reconciliation of
	// OpenEBS resource
	reconciler, err :=
		NewReconciler(
			ReconcilerConfig{
				ObservedOpenEBS:                           observedOpenEBS,
				ObservedOpenEBSComponents:                 observedOpenEBSComponents,
				ObservedOpenEBSCRDs:                       observedOpenEBSCRDs,
				ObservedOpenEBSClusterRoleAndRoleBindings: observedOpenEBSClusterRoleAndRoleBindings,
				ObservedCStorCSIDriver:                    observedCStorCSIDriver,
			})
	if err != nil {
		errHandler.handle(err)
		return nil
	}
	resp, err := reconciler.Reconcile()
	if err != nil {
		errHandler.handle(err)
		return nil
	}
	// add all the desired OpenEBS components as attachments in the response
	if resp.DesiredOpenEBSComponents != nil {
		for _, desiredOpenEBSComponent := range resp.DesiredOpenEBSComponents {
			response.Attachments = append(response.Attachments, desiredOpenEBSComponent)
		}
	}
	// update the components that needs to be deleted
	if resp.ExplicitDeletes != nil {
		for _, componentToDelete := range resp.ExplicitDeletes {
			response.ExplicitDeletes = append(response.ExplicitDeletes, componentToDelete)
		}
	}
	// isOpenEBSExplicitlyUpdated states if OpenEBS CR is going to be updated explicitly
	// or not, this happens only in case of ISCSI installation where OpenEBS is updated
	// explicitly if ISCSI client installation is successful.
	var isOpenEBSExplicitlyUpdated bool
	// update the components that needs to be explicitly updated
	if resp.ExplicitUpdates != nil {
		for _, componentToUpdate := range resp.ExplicitUpdates {
			if componentToUpdate.GetKind() == string(types.KindOpenEBS) {
				isOpenEBSExplicitlyUpdated = true
			}
			response.ExplicitUpdates = append(response.ExplicitUpdates, componentToUpdate)
		}
	}

	glog.V(2).Infof(
		"OpenEBS %s %s reconciled successfully: %s",
		request.Watch.GetNamespace(), request.Watch.GetName(),
		metac.GetDetailsFromResponse(response),
	)

	// If OpenEBS is not being updated explicitly then use this to update OpenEBS status.
	if !isOpenEBSExplicitlyUpdated {
		// construct the success handler
		successHandler := &reconcileSuccessHandler{
			openebs:      request.Watch,
			hookResponse: response,
		}
		successHandler.handle()
	}

	return nil
}

// Reconciler enables reconciliation of OpenEBS instance
type Reconciler struct {
	ObservedOpenEBS                           *types.OpenEBS
	ObservedOpenEBSComponents                 []*unstructured.Unstructured
	ObservedOpenEBSCRDs                       []*unstructured.Unstructured
	ObservedOpenEBSClusterRoleAndRoleBindings []*unstructured.Unstructured
	ObservedCStorCSIDriver                    *unstructured.Unstructured
}

// ReconcilerConfig is a helper structure used to create a
// new instance of Reconciler
type ReconcilerConfig struct {
	ObservedOpenEBS                           *unstructured.Unstructured
	ObservedOpenEBSComponents                 []*unstructured.Unstructured
	ObservedOpenEBSCRDs                       []*unstructured.Unstructured
	ObservedOpenEBSClusterRoleAndRoleBindings []*unstructured.Unstructured
	ObservedCStorCSIDriver                    *unstructured.Unstructured
}

// ReconcileResponse is a helper struct used to form the response
// of a successful reconciliation
type ReconcileResponse struct {
	DesiredOpenEBS           *unstructured.Unstructured
	DesiredOpenEBSComponents []*unstructured.Unstructured
	ExplicitDeletes          []*unstructured.Unstructured
	ExplicitUpdates          []*unstructured.Unstructured
}

// Planner ensures if any of the instances need
// to be created, or updated.
type Planner struct {
	ObservedOpenEBS                           *types.OpenEBS
	ObservedOpenEBSComponents                 []*unstructured.Unstructured
	ObservedOpenEBSCRDs                       []*unstructured.Unstructured
	ObservedOpenEBSClusterRoleAndRoleBindings []*unstructured.Unstructured
	ObservedCStorCSIDriver                    *unstructured.Unstructured

	DesiredOpenEBSCRDs []*unstructured.Unstructured

	ComponentManifests map[string]*unstructured.Unstructured
	ExplicitDeletes    []*unstructured.Unstructured
	ExplicitUpdates    []*unstructured.Unstructured
}

// NewReconciler returns a new instance of Reconciler
func NewReconciler(config ReconcilerConfig) (*Reconciler, error) {
	// transform OpenEBS from unstructured to typed
	var openebsTyped types.OpenEBS
	openebsRaw, err := config.ObservedOpenEBS.MarshalJSON()
	if err != nil {
		return nil, errors.Wrapf(err, "Can't marshal OpenEBS")
	}
	err = json.Unmarshal(openebsRaw, &openebsTyped)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't unmarshal OpenEBS")
	}
	// use above constructed object to build Reconciler instance
	return &Reconciler{
		ObservedOpenEBS:                           &openebsTyped,
		ObservedOpenEBSComponents:                 config.ObservedOpenEBSComponents,
		ObservedOpenEBSCRDs:                       config.ObservedOpenEBSCRDs,
		ObservedOpenEBSClusterRoleAndRoleBindings: config.ObservedOpenEBSClusterRoleAndRoleBindings,
		ObservedCStorCSIDriver:                    config.ObservedCStorCSIDriver,
	}, nil
}

// Reconcile runs through the reconciliation logic
//
// NOTE:
//	Due care has been taken to let this logic be idempotent
func (r *Reconciler) Reconcile() (ReconcileResponse, error) {
	planner := Planner{
		ObservedOpenEBS:                           r.ObservedOpenEBS,
		ObservedOpenEBSComponents:                 r.ObservedOpenEBSComponents,
		ObservedOpenEBSCRDs:                       r.ObservedOpenEBSCRDs,
		ObservedOpenEBSClusterRoleAndRoleBindings: r.ObservedOpenEBSClusterRoleAndRoleBindings,
		ObservedCStorCSIDriver:                    r.ObservedCStorCSIDriver,
	}
	return planner.Plan()
}

// Plan builds the desired instances/components of OpenEBS
func (p *Planner) Plan() (ReconcileResponse, error) {
	// check if some dependencies/tools/components needs to be installed before OpenEBS
	// installation.
	err := p.getPreInstallationManifests()
	if err != nil {
		return ReconcileResponse{}, err
	}
	err = p.init()
	if err != nil {
		return ReconcileResponse{}, err
	}
	return p.getDesiredOpenEBSComponents(), nil
}

// getDesiredOpenEBSComponents gets all the desired OpenEBS components which
// needs to be created or updated.
func (p *Planner) getDesiredOpenEBSComponents() ReconcileResponse {
	response := ReconcileResponse{}
	for key, value := range p.ComponentManifests {
		if key == "_" {
			continue
		}
		response.DesiredOpenEBSComponents = append(response.DesiredOpenEBSComponents, value)
	}
	// update the components that needs to be deleted
	for _, componentToDelete := range p.ExplicitDeletes {
		response.ExplicitDeletes = append(response.ExplicitDeletes, componentToDelete)
	}
	// update the components that needs to be explicitly updated
	for _, componentToUpdate := range p.ExplicitUpdates {
		response.ExplicitUpdates = append(response.ExplicitUpdates, componentToUpdate)
	}
	// add the observed OpenEBS CRDs to desired OpenEBS CRDs that are not already present
	// in the desiredOpenEBS components list.
	if len(p.ObservedOpenEBSCRDs) > 0 {
		for _, observedCRD := range p.ObservedOpenEBSCRDs {
			key := observedCRD.GetName() + "_" + observedCRD.GetKind()
			if desiredCRD, exist := p.ComponentManifests[key]; exist {
				// If already exists then check if the APIVersion is same or not, if
				// not then add it to the desired OpenEBS components list.
				if desiredCRD.GetAPIVersion() != observedCRD.GetAPIVersion() {
					response.DesiredOpenEBSComponents = append(response.DesiredOpenEBSComponents,
						observedCRD)
				}
			} else {
				response.DesiredOpenEBSComponents = append(response.DesiredOpenEBSComponents,
					observedCRD)
			}
		}
	}

	// add the observed OpenEBS clusterroles & clusterrolebindings to desired OpenEBS clusterroles & clusterrolebindings
	// that are not already present in the desiredOpenEBS components list.
	if len(p.ObservedOpenEBSClusterRoleAndRoleBindings) > 0 {
		for _, observedOpenEBSClusterRoleAndRoleBindings := range p.ObservedOpenEBSClusterRoleAndRoleBindings {
			key := observedOpenEBSClusterRoleAndRoleBindings.GetName() + "_" + observedOpenEBSClusterRoleAndRoleBindings.GetKind()
			if desiredClusterRoleAndRoleBinding, exist := p.ComponentManifests[key]; exist {
				// If already exists then check if the APIVersion is same or not, if
				// not then add it to the desired OpenEBS components list.
				if desiredClusterRoleAndRoleBinding.GetAPIVersion() != observedOpenEBSClusterRoleAndRoleBindings.GetAPIVersion() {
					response.DesiredOpenEBSComponents = append(response.DesiredOpenEBSComponents,
						observedOpenEBSClusterRoleAndRoleBindings)
				}
			} else {
				response.DesiredOpenEBSComponents = append(response.DesiredOpenEBSComponents,
					observedOpenEBSClusterRoleAndRoleBindings)
			}
		}
	}

	return response
}

func (p *Planner) init() error {
	var initFuncs = []func() error{
		p.getManifests,
		p.setDefaultImagePullPolicyIfNotSet,
		p.setDefaultStoragePathIfNotSet,
		p.setDefaultImagePrefixIfNotSet,
		p.setImageTagSuffixIfPresent,
		p.setDefaultStorageConfigIfNotSet,
		p.setAPIServerDefaultsIfNotSet,
		p.setProvisionerDefaultsIfNotSet,
		p.setLocalProvisionerDefaultsIfNotSet,
		p.setSnapshotOperatorDefaultsIfNotSet,
		p.setAdmissionServerDefaultsIfNotSet,
		p.setNDMDefaultsIfNotSet,
		p.setNDMOperatorDefaultsIfNotSet,
		p.setNDMConfigMapDefaultsIfNotSet,
		p.setJIVADefaultsIfNotSet,
		p.setCStorDefaultsIfNotSet,
		p.setMayastorDefaultsIfNotSet,
		p.setHelperDefaultsIfNotSet,
		p.setPoliciesDefaultsIfNotSet,
		p.setAnalyticsDefaultsIfNotSet,
		p.getDesiredValuesFromObservedResources,
		p.removeDisabledManifests,
		p.getDesiredManifests,
	}
	for _, fn := range initFuncs {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}
