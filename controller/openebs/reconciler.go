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

	"github.com/mayadata-io/openebs-operator/k8s"
	"github.com/mayadata-io/openebs-operator/pkg/utils/metac"
	"github.com/mayadata-io/openebs-operator/types"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"openebs.io/metac/controller/generic"
)

type reconcileErrHandler struct {
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

	conds, mergeErr :=
		k8s.MergeStatusConditions(
			h.openebs, types.MakeOpenEBSReconcileErrCond(err),
		)
	if mergeErr != nil {
		glog.Errorf(
			"Failed to reconcile OpenEBS %s %s: Can't set status conditions: %+v",
			h.openebs.GetNamespace(), h.openebs.GetName(), mergeErr,
		)
		// Note: Merge error will reset the conditions which will make
		// things worse since various controllers will be reconciling
		// based on these conditions.
		//
		// Hence it is better to set response status as nil to let metac
		// preserve old status conditions if any.
		h.hookResponse.Status = nil
	} else {
		// response status will be set against the watch's status by metac
		h.hookResponse.Status = map[string]interface{}{}
		h.hookResponse.Status["phase"] = types.OpenEBSStatusPhaseError
		h.hookResponse.Status["conditions"] = conds
	}
	// this will stop further reconciliation at metac since there was
	// an error
	h.hookResponse.SkipReconcile = true
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

	var openebsObj *unstructured.Unstructured
	for _, attachment := range request.Attachments.List() {
		// this watch resource must be present in the list of attachments
		if request.Watch.GetUID() == attachment.GetUID() &&
			attachment.GetKind() == string(types.KindOpenEBS) {
			// this is the required OpenEBS
			openebsObj = attachment
			// add this to the response later after completion of its
			// reconcile logic
			continue
		}
		response.Attachments = append(response.Attachments, attachment)
	}

	if openebsObj == nil {
		errHandler.handle(
			errors.Errorf("Can't reconcile: OpenEBS not found in attachments"),
		)
		return nil
	}

	// reconciler is the one that will perform reconciliation of
	// OpenEBS resource
	reconciler, err :=
		NewReconciler(
			openebsObj,
			request.Attachments.List(),
		)
	if err != nil {
		errHandler.handle(err)
		return nil
	}
	_, err = reconciler.Reconcile()
	if err != nil {
		errHandler.handle(err)
		return nil
	}

	// Note: We are not updating the updated response back to the OpenEBS
	// object.
	// add updated OpenEBS to response
	//response.Attachments = append(response.Attachments, op.OpenEBS)

	glog.V(2).Infof(
		"OpenEBS %s %s reconciled successfully: %s",
		request.Watch.GetNamespace(), request.Watch.GetName(),
		metac.GetDetailsFromResponse(response),
	)

	return nil
}

// Reconciler enables reconciliation of OpenEBS instance
type Reconciler struct {
	OpenEBS   *types.OpenEBS
	Resources []*unstructured.Unstructured
}

// ReconcileResponse is a helper struct used to form the response
// of a successful reconciliation
type ReconcileResponse struct {
	OpenEBS *unstructured.Unstructured
}

// NewReconciler returns a new instance of Reconciler
func NewReconciler(
	openebs *unstructured.Unstructured,
	resources []*unstructured.Unstructured,
) (*Reconciler, error) {
	r := &Reconciler{
		Resources: resources,
	}

	// transform OpenEBS from unstructured to typed
	var openebsTyped types.OpenEBS
	openebsRaw, err := openebs.MarshalJSON()
	if err != nil {
		return nil, errors.Wrapf(err, "Can't marshal OpenEBS")
	}
	err = json.Unmarshal(openebsRaw, &openebsTyped)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't unmarshal OpenEBS")
	}

	// update the reconciler instance with config & related fields
	r.OpenEBS = &openebsTyped

	return r, nil
}

// Reconcile runs through the reconciliation logic
//
// NOTE:
//	Due care has been taken to let this logic be idempotent
func (r *Reconciler) Reconcile() (ReconcileResponse, error) {
	// Updating the previousOpenEBS field so that it can be used
	// in order to return the reconciler response.
	syncFns := []func() error{
		r.validateOpenEBSAndSetDefaultsIfNotSet,
		r.createOrUpdateComponents,
	}
	for _, syncFn := range syncFns {
		err := syncFn()
		if err != nil {
			return ReconcileResponse{}, err
		}
	}
	// reset previous errors if any
	r.resetOpenEBSReconcileErrorIfAny()
	return r.makeReconcileResponse()
}

// makeReconcileResponse builds reconcile response
func (r *Reconciler) makeReconcileResponse() (ReconcileResponse, error) {
	// convert updated OpenEBS from typed to unstruct
	openebsRaw, err := json.Marshal(r.OpenEBS)
	if err != nil {
		return ReconcileResponse{},
			errors.Wrapf(err, "Can't marshal OpenEBS")
	}
	var openebs unstructured.Unstructured
	err = json.Unmarshal(openebsRaw, &openebs)
	if err != nil {
		return ReconcileResponse{},
			errors.Wrapf(err, "Can't unmarshal OpenEBS")
	}

	return ReconcileResponse{
		OpenEBS: &openebs,
	}, nil
}

// resetOpenEBSReconcileErrorIfAny removes ReconcileError
// if any associated with OpenEBS. This removes
// ReconcileError if it ever happened during previous
// reconciliations.
func (r *Reconciler) resetOpenEBSReconcileErrorIfAny() {
	types.MergeNoReconcileErrorOnOpenEBS(r.OpenEBS)
}

func (r *Reconciler) validateOpenEBSAndSetDefaultsIfNotSet() error {
	// set all the defaults for all the components
	// if not already set.
	setDefaultFns := []func() error{
		r.setDefaultImagePullPolicyIfNotSet,
		r.setDefaultStoragePathIfNotSet,
		r.setDefaultImagePrefixIfNotSet,
		r.setDefaultStorageConfigIfNotSet,
		r.setAPIServerDefaultsIfNotSet,
		r.setProvisionerDefaultsIfNotSet,
		r.setLocalProvisionerDefaultsIfNotSet,
		r.setSnapshotOperatorDefaultsIfNotSet,
		r.setAdmissionServerDefaultsIfNotSet,
		r.setNDMDefaultsIfNotSet,
		r.setNDMOperatorDefaultsIfNotSet,
		r.setJIVADefaultsIfNotSet,
		r.setCStorDefaultsIfNotSet,
		r.setHelperDefaultsIfNotSet,
		r.setPoliciesDefaultsIfNotSet,
		r.setAnalyticsDefaultsIfNotSet,
	}
	for _, setDefaultFn := range setDefaultFns {
		err := setDefaultFn()
		if err != nil {
			return err
		}
	}
	return nil
}

// createOrUpdateComponents createsupdates all the required OpenEBS components
// as per the manifests for a given OpenEBS version.
func (r *Reconciler) createOrUpdateComponents() error {
	// get the manifest which can be used to install/update
	// OpenEBS components.
	manifests, err := r.getManifests()
	if err != nil {
		return errors.Errorf(
			"Error getting manifests for creating/updating components: %+v", err)
	}
	// remove manifests which are disabled or should not be installed
	// TODO: Delete the components which are disabled after installation.
	manifests, err = r.removeDisabledManifests(manifests)
	if err != nil {
		return err
	}
	// update the manifests as per the given configuration in the
	// OpenEBS CR.
	manifests, err = r.updateManifests(manifests)
	if err != nil {
		return err
	}
	// create/update all the OpenEBS components as per the
	// updated manifests.
	err = deployComponents(manifests)
	if err != nil {
		return err
	}

	return nil
}
