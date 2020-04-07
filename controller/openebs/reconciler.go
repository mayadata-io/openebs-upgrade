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

	"mayadata.io/openebs-upgrade/pkg/utils/metac"
	"mayadata.io/openebs-upgrade/types"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

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
				ObservedOpenEBS:           observedOpenEBS,
				observedOpenEBSComponents: observedOpenEBSComponents,
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
	if resp.DesiredOpenEBSComponets != nil {
		for _, desiredOpenEBSComponent := range resp.DesiredOpenEBSComponets {
			response.Attachments = append(response.Attachments, desiredOpenEBSComponent)
		}
	}

	glog.V(2).Infof(
		"OpenEBS %s %s reconciled successfully: %s",
		request.Watch.GetNamespace(), request.Watch.GetName(),
		metac.GetDetailsFromResponse(response),
	)

	// construct the success handler
	successHandler := &reconcileSuccessHandler{
		openebs:      request.Watch,
		hookResponse: response,
	}
	successHandler.handle()

	return nil
}

// Reconciler enables reconciliation of OpenEBS instance
type Reconciler struct {
	ObservedOpenEBS           *types.OpenEBS
	observedOpenEBSComponents []*unstructured.Unstructured
}

// ReconcilerConfig is a helper structure used to create a
// new instance of Reconciler
type ReconcilerConfig struct {
	ObservedOpenEBS           *unstructured.Unstructured
	observedOpenEBSComponents []*unstructured.Unstructured
}

// ReconcileResponse is a helper struct used to form the response
// of a successful reconciliation
type ReconcileResponse struct {
	DesiredOpenEBS          *unstructured.Unstructured
	DesiredOpenEBSComponets []*unstructured.Unstructured
}

// Planner ensures if any of the instances need
// to be created, or updated.
type Planner struct {
	ObservedOpenEBS           *types.OpenEBS
	observedOpenEBSComponents []*unstructured.Unstructured

	ComponentManifests map[string]*unstructured.Unstructured
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
		ObservedOpenEBS:           &openebsTyped,
		observedOpenEBSComponents: config.observedOpenEBSComponents,
	}, nil
}

// Reconcile runs through the reconciliation logic
//
// NOTE:
//	Due care has been taken to let this logic be idempotent
func (r *Reconciler) Reconcile() (ReconcileResponse, error) {
	planner := Planner{
		ObservedOpenEBS:           r.ObservedOpenEBS,
		observedOpenEBSComponents: r.observedOpenEBSComponents,
	}
	return planner.Plan()
}

// Plan builds the desired instances/components of OpenEBS
func (p *Planner) Plan() (ReconcileResponse, error) {
	err := p.init()
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
		response.DesiredOpenEBSComponets = append(response.DesiredOpenEBSComponets, value)
	}
	return response
}

func (p *Planner) init() error {
	var initFuncs = []func() error{
		p.getManifests,
		// TODO: Remove once 1.9.0 images are available
		p.updateVersionFor190,
		p.setDefaultImagePullPolicyIfNotSet,
		p.setDefaultStoragePathIfNotSet,
		p.setDefaultImagePrefixIfNotSet,
		p.setDefaultStorageConfigIfNotSet,
		p.setAPIServerDefaultsIfNotSet,
		p.setProvisionerDefaultsIfNotSet,
		p.setLocalProvisionerDefaultsIfNotSet,
		p.setSnapshotOperatorDefaultsIfNotSet,
		p.setAdmissionServerDefaultsIfNotSet,
		p.setNDMDefaultsIfNotSet,
		p.setNDMOperatorDefaultsIfNotSet,
		p.setJIVADefaultsIfNotSet,
		p.setCStorDefaultsIfNotSet,
		p.setHelperDefaultsIfNotSet,
		p.setPoliciesDefaultsIfNotSet,
		p.setAnalyticsDefaultsIfNotSet,
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
