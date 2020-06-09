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

package adoptopenebs

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
	adoptopenebs *unstructured.Unstructured
	hookResponse *generic.SyncHookResponse
}

type reconcileSuccessHandler struct {
	adoptopenebs *unstructured.Unstructured
	hookResponse *generic.SyncHookResponse
}

func (h *reconcileErrHandler) handle(err error) {
	// Error has been handled elaborately. This logic ensures
	// error message is propagated to the resource & hence seen via
	// 'kubectl get adoptopenebs -n {some-namespace} -oyaml'.
	//
	// In addition, logging has been done to check for error messages
	// from this pod's logs.
	glog.Errorf(
		"Failed to reconcile adoptOpenEBS %s %s: %+v",
		h.adoptopenebs.GetNamespace(), h.adoptopenebs.GetName(), err,
	)
	// response status will be set against the watch's status by metac
	h.hookResponse.Status = map[string]interface{}{}
	h.hookResponse.Status["phase"] = types.AdoptOpenEBSStatusPhaseFailed
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
	h.hookResponse.Status["phase"] = types.AdoptOpenEBSStatusPhaseOnline
}

// Sync implements the idempotent logic to reconcile AdoptOpenEBS
//
// NOTE:
// 	SyncHookRequest is the payload received as part of reconcile
// request. Similarly, SyncHookResponse is the payload sent as a
// response as part of reconcile request.
//
// NOTE:
//	SyncHookRequest uses AdoptOpenEBS as the watched resource.
// SyncHookResponse has the resources that forms the desired state
// w.r.t the watched resource.
//
// NOTE:
//	Returning error will panic this process. We would rather want this
// controller to run continuously. Hence, the errors are logged and at
// the same time, these errors are posted against AdoptOpenEBS's
// status.
func Sync(request *generic.SyncHookRequest, response *generic.SyncHookResponse) error {
	if request == nil {
		return errors.Errorf("Failed to reconcile AdoptOpenEBS: Nil request found")
	}
	if response == nil {
		return errors.Errorf("Failed to reconcile AdoptOpenEBS: Nil response found")
	}
	// Nothing needs to be done if there are no attachments in request
	//
	// NOTE:
	// 	It is expected to have AdoptOpenEBS as an attachment
	// resource as well as the resource under watch.
	if request.Attachments == nil || request.Attachments.IsEmpty() {
		response.SkipReconcile = true
		return nil
	}

	glog.V(3).Infof(
		"Will reconcile AdoptOpenEBS %s %s:",
		request.Watch.GetNamespace(), request.Watch.GetName(),
	)

	// construct the error handler
	errHandler := &reconcileErrHandler{
		adoptopenebs: request.Watch,
		hookResponse: response,
	}

	var observedAdoptOpenEBS *unstructured.Unstructured
	var observedOpenEBS *unstructured.Unstructured
	var observedAdoptOpenEBSComponents []*unstructured.Unstructured
	for _, attachment := range request.Attachments.List() {
		// this watch resource must be present in the list of attachments
		if request.Watch.GetUID() == attachment.GetUID() &&
			attachment.GetKind() == string(types.KindAdoptOpenEBS) {
			// this is the required AdoptOpenEBS
			observedAdoptOpenEBS = attachment
			// add this to the response later after completion of its
			// reconcile logic
			continue
		}
		if attachment.GetKind() == string(types.KindOpenEBS) {
			// this is the required OpenEBS
			observedOpenEBS = attachment
			continue
		}
		// If the attachments are not of Kind: AdoptOpenEBS then it will be
		// considered as an OpenEBS component.
		if attachment.GetKind() != string(types.KindAdoptOpenEBS) {
			observedAdoptOpenEBSComponents = append(observedAdoptOpenEBSComponents, attachment)
		}
	}

	if observedAdoptOpenEBS == nil {
		errHandler.handle(
			errors.Errorf("Can't reconcile: AdoptOpenEBS not found in attachments"),
		)
		return nil
	}

	// reconciler is the one that will perform reconciliation of
	// AdoptOpenEBS resource
	reconciler, err :=
		NewReconciler(
			ReconcilerConfig{
				ObservedAdoptOpenEBS:           observedAdoptOpenEBS,
				ObservedOpenEBS:                observedOpenEBS,
				ObservedAdoptOpenEBSComponents: observedAdoptOpenEBSComponents,
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
	// add all the desired adoptOpenEBS, OpenEBS and adoptOpenEBS components as attachments
	// in the response.
	if resp.DesiredAdoptOpenEBS != nil {
		response.Attachments = append(response.Attachments, resp.DesiredAdoptOpenEBS)
	}
	if resp.DesiredOpenEBS != nil {
		response.Attachments = append(response.Attachments, resp.DesiredOpenEBS)
	}
	if resp.DesiredAdoptOpenEBSComponents != nil {
		for _, desiredAdoptOpenEBSComponent := range resp.DesiredAdoptOpenEBSComponents {
			response.Attachments = append(response.Attachments, desiredAdoptOpenEBSComponent)
		}
	}

	glog.V(2).Infof(
		"AdoptOpenEBS %s %s reconciled successfully: %s",
		request.Watch.GetNamespace(), request.Watch.GetName(),
		metac.GetDetailsFromResponse(response),
	)

	// construct the success handler
	successHandler := &reconcileSuccessHandler{
		adoptopenebs: request.Watch,
		hookResponse: response,
	}
	successHandler.handle()

	return nil
}

// Reconciler enables reconciliation of AdoptOpenEBS instance
type Reconciler struct {
	ObservedAdoptOpenEBS           *types.AdoptOpenEBS
	ObservedOpenEBS                *types.OpenEBS
	ObservedAdoptOpenEBSComponents []*unstructured.Unstructured
}

// ReconcilerConfig is a helper structure used to create a
// new instance of Reconciler
type ReconcilerConfig struct {
	ObservedAdoptOpenEBS           *unstructured.Unstructured
	ObservedOpenEBS                *unstructured.Unstructured
	ObservedAdoptOpenEBSComponents []*unstructured.Unstructured
}

// ReconcileResponse is a helper struct used to form the response
// of a successful reconciliation
type ReconcileResponse struct {
	DesiredAdoptOpenEBS           *unstructured.Unstructured
	DesiredOpenEBS                *unstructured.Unstructured
	DesiredAdoptOpenEBSComponents []*unstructured.Unstructured
}

// Planner ensures if any of the instances need
// to be created, or updated.
type Planner struct {
	ObservedAdoptOpenEBS           *types.AdoptOpenEBS
	ObservedOpenEBS                *types.OpenEBS
	ObservedAdoptOpenEBSComponents []*unstructured.Unstructured

	DesiredAdoptOpenEBS           *unstructured.Unstructured
	DesiredOpenEBS                *unstructured.Unstructured
	DesiredAdoptOpenEBSComponents []*unstructured.Unstructured

	OpenEBSVersion             string
	DefaultStoragePath         string
	ImagePrefix                string
	CreateDefaultStorageConfig bool
	ImagePullPolicy            string
	JivaCtrlImageTag           string
	JivaReplicaImageTag        string
	JivaReplicaCount           int64
	CStorTargetImageTag        string
	CStorPoolImageTag          string
	CStorPoolMgmtImageTag      string
	CStorVolumeMgmtImageTag    string
	CStorVolumeManagerImageTag string
	CSPIMgmtImageTag           string
	VolumeMonitorImageTag      string
	PoolExporterImageTag       string
	HelperImageTag             string
	EnableAnalytics            bool

	Resources              *unstructured.Unstructured
	APIServerConfig        *unstructured.Unstructured
	ProvisionerConfig      *unstructured.Unstructured
	LocalProvisionerConfig *unstructured.Unstructured
	AdmissionServerConfig  *unstructured.Unstructured
	SnapshotOperatorConfig *unstructured.Unstructured
	NDMDaemonConfig        *unstructured.Unstructured
	NDMOperatorConfig      *unstructured.Unstructured
	NDMConfigMapConfig     *unstructured.Unstructured
	JivaConfig             *unstructured.Unstructured
	CstorConfig            *unstructured.Unstructured
	HelperConfig           *unstructured.Unstructured
	PoliciesConfig         *unstructured.Unstructured
	AnalyticsConfig        *unstructured.Unstructured
}

// NewReconciler returns a new instance of Reconciler
func NewReconciler(config ReconcilerConfig) (*Reconciler, error) {
	// transform adoptOpenEBS from unstructured to typed
	var adoptOpenEBSTyped types.AdoptOpenEBS
	adoptOpenEBSRaw, err := config.ObservedAdoptOpenEBS.MarshalJSON()
	if err != nil {
		return nil, errors.Wrapf(err, "Can't marshal adoptOpenEBS")
	}
	err = json.Unmarshal(adoptOpenEBSRaw, &adoptOpenEBSTyped)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't unmarshal adoptOpenEBS")
	}
	// transform OpenEBS if present from unstructured to typed
	var openebsTyped types.OpenEBS
	if config.ObservedOpenEBS != nil {
		openebsRaw, err := config.ObservedOpenEBS.MarshalJSON()
		if err != nil {
			return nil, errors.Wrapf(err, "Can't marshal OpenEBS")
		}
		err = json.Unmarshal(openebsRaw, &openebsTyped)
		if err != nil {
			return nil, errors.Wrapf(err, "Can't unmarshal OpenEBS")
		}
	}
	// use above constructed object to build Reconciler instance
	return &Reconciler{
		ObservedAdoptOpenEBS:           &adoptOpenEBSTyped,
		ObservedOpenEBS:                &openebsTyped,
		ObservedAdoptOpenEBSComponents: config.ObservedAdoptOpenEBSComponents,
	}, nil
}

// Reconcile runs through the reconciliation logic
//
// NOTE:
//	Due care has been taken to let this logic be idempotent
func (r *Reconciler) Reconcile() (ReconcileResponse, error) {
	planner := Planner{
		ObservedAdoptOpenEBS:           r.ObservedAdoptOpenEBS,
		ObservedOpenEBS:                r.ObservedOpenEBS,
		ObservedAdoptOpenEBSComponents: r.ObservedAdoptOpenEBSComponents,
	}
	return planner.Plan()
}

// Plan builds the desired instances/components of OpenEBS
func (p *Planner) Plan() (ReconcileResponse, error) {
	err := p.init()
	if err != nil {
		return ReconcileResponse{}, err
	}
	return p.plan(), nil
}

// plan forms the desired OpenEBS, AdoptOpenEBS and the AdoptOpenEBS components
// with the help of discovered OpenEBS components.
func (p *Planner) plan() ReconcileResponse {
	response := ReconcileResponse{}
	response.DesiredAdoptOpenEBS = p.DesiredAdoptOpenEBS
	response.DesiredOpenEBS = p.DesiredOpenEBS
	for _, value := range p.DesiredAdoptOpenEBSComponents {
		response.DesiredAdoptOpenEBSComponents = append(
			response.DesiredAdoptOpenEBSComponents, value)
	}
	return response
}

func (p *Planner) init() error {
	var initFuncs = []func() error{
		p.IdentifyOpenEBSVersion,
		p.getDesiredAdoptOpenEBS,
		// Ordering of getDesiredAdoptOpenEBSComponents
		// and getDesiredOpenEBS must be maintained.
		p.getDesiredAdoptOpenEBSComponents,
		p.getDesiredOpenEBS,
	}
	for _, fn := range initFuncs {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

// getDesiredAdoptOpenEBS returns the desired adoptOpenEBS structure.
func (p *Planner) getDesiredAdoptOpenEBS() error {
	desiredAdoptOpenEBS := &unstructured.Unstructured{
		Object: make(map[string]interface{}, 0),
	}
	rawObservedAdoptOpenEBS, err := json.Marshal(p.ObservedAdoptOpenEBS)
	if err != nil {
		return errors.Errorf("Error marshalling observed adoptOpenEBS: %+v, Error: %+v",
			p.ObservedAdoptOpenEBS, err)
	}
	if err = json.Unmarshal(rawObservedAdoptOpenEBS, &desiredAdoptOpenEBS.Object); err != nil {
		return errors.Errorf(
			"Error unmarshalling observed adoptOpenEBS into unstructured, Error: %+v",
			err)
	}
	p.DesiredAdoptOpenEBS = desiredAdoptOpenEBS
	return nil
}

// getDesiredAdoptOpenEBSComponents returns the desired adoptOpenEBS components.
func (p *Planner) getDesiredAdoptOpenEBSComponents() error {
	var err error
	desiredAdoptOpenEBSComponents := p.DesiredAdoptOpenEBSComponents
	for _, observedAdoptOpenEBSComponent := range p.ObservedAdoptOpenEBSComponents {
		openEBSIdentifier := OpenEBSIdentifier{
			Object: observedAdoptOpenEBSComponent,
		}
		// Identify the type of this OpenEBS component and form the config based on that.
		componentType, err := openEBSIdentifier.IdentifyOpenEBSComponentType()
		if err != nil {
			return errors.Errorf(
				"Error identifying OpenEBS component type for component: %+v, Error: %+v",
				observedAdoptOpenEBSComponent, err)
		}
		// If there is no error and still the type is not recognised then continue without
		// forming any config. This could be some unknonwn/yet-to-be-identified OpenEBS resource.
		if componentType == "" {
			continue
		}
		// Form and set the OpenEBS CR config as per the components observed configuration.
		err = p.formComponentOpenEBSConfig(observedAdoptOpenEBSComponent, componentType)
		if err != nil {
			return err
		}
		// update the component labels.
		err = p.addAdoptOpenEBSLabels(observedAdoptOpenEBSComponent)
		if err != nil {
			return err
		}
		desiredAdoptOpenEBSComponents = append(desiredAdoptOpenEBSComponents, observedAdoptOpenEBSComponent)
	}
	// Now, form some of the configs which are common to all the components or can be
	// used across components.
	err = p.formCommonOpenEBSConfig()
	if err != nil {
		return err
	}
	return nil
}

// addAdoptOpenEBSLabels adds the desired labels to the adopted OpenEBS components.
// It adds two labels, one which is specific to adoptOpenEBS i.e., it helps in identifying
// if the components have been adopted or installed, label is
// openebs-upgrade.dao.mayadata.io/adopt=true.
// The second label which is added is a common label which is added both in case
// of installation as well as adoption i.e., openebs-upgrade.dao.mayadata.io/managed=true.
func (p *Planner) addAdoptOpenEBSLabels(component *unstructured.Unstructured) error {
	desiredLabels := component.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	desiredLabels[types.OpenEBSUpgradeDAOManagedLabelKey] =
		types.OpenEBSUpgradeDAOManagedLabelValue
	desiredLabels[types.OpenEBSUpgradeDAOAdoptLabelKey] =
		types.OpenEBSUpgradeDAOAdoptLabelValue
	component.SetLabels(desiredLabels)

	return nil
}

// getDesiredOpenEBS returns the desired OpenEBS structure as per the discovered OpenEBS
// components and their configuration.
func (p *Planner) getDesiredOpenEBS() error {
	openebs := &unstructured.Unstructured{}
	openebs.SetUnstructuredContent(map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":      "my-openebs-123",
			"namespace": p.ObservedAdoptOpenEBS.Namespace,
		},
		"spec": map[string]interface{}{
			"version":                    p.OpenEBSVersion,
			"defaultStoragePath":         p.DefaultStoragePath,
			"createDefaultStorageConfig": p.CreateDefaultStorageConfig,
			"imagePrefix":                p.ImagePrefix,
			"imagePullPolicy":            p.ImagePullPolicy,
			"resources":                  p.Resources,
			"apiServer":                  p.APIServerConfig,
			"provisioner":                p.ProvisionerConfig,
			"localProvisioner":           p.LocalProvisionerConfig,
			"snapshotOperator":           p.SnapshotOperatorConfig,
			"ndmDaemon":                  p.NDMDaemonConfig,
			"ndmOperator":                p.NDMOperatorConfig,
			"ndmConfigMap":               p.NDMConfigMapConfig,
			"jivaConfig":                 p.JivaConfig,
			"cstorConfig":                p.CstorConfig,
			"admissionServer":            p.AdmissionServerConfig,
			"helper":                     p.HelperConfig,
			"policies":                   p.PoliciesConfig,
			"analytics":                  p.AnalyticsConfig,
		},
	})
	openebs.SetKind(string(types.KindOpenEBS))
	openebs.SetAPIVersion(string(types.APIVersionDAOMayaDataV1Alpha1))

	// set the desired OpenEBS structure
	p.DesiredOpenEBS = openebs
	return nil
}
